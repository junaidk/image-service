package http

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	imageapi "github.com/junaidk/image-service"
	"github.com/junaidk/image-service/internal/image"
	"github.com/junaidk/image-service/internal/token"
)

func (s *Server) imageRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/upload/{token}", s.createImageHandler)
	r.Get("/{image-id}", s.showImageHandler)
	return r
}

func (s *Server) createImageHandler(w http.ResponseWriter, r *http.Request) {
	expToken := chi.URLParam(r, "token")
	if !token.Validate(expToken) {
		s.badRequestResponse(w, r, fmt.Errorf("invalid toekn"))
		return
	}

	// 10 MB file size
	reqSize := int64(10 * 1024 * 1024)

	r.Body = http.MaxBytesReader(w, r.Body, reqSize)
	err := r.ParseMultipartForm(reqSize)

	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	imagesResp := map[string]string{}

	files := r.MultipartForm.File["images"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			s.serverErrorResponse(w, r, err)
			return
		}
		defer file.Close()

		imgData := &imageapi.Image{}

		tmp, err := image.GetMetadata(file, path.Ext(fileHeader.Filename))
		if err == nil {
			imgData = tmp
		} else {
			s.logError(r, err)
		}
		_, _ = file.Seek(0, io.SeekStart)

		hash := sha256.New()
		io.Copy(hash, file)
		hashValue := hex.EncodeToString(hash.Sum(nil))

		imgData.ID = uuid.New()
		imgData.Hash = hashValue
		imgData.Name = fileHeader.Filename

		err = s.ImageService.CreateImage(r.Context(), imgData)
		isImgExist := false
		if err != nil {
			switch imageapi.ErrorCode(err) {
			case imageapi.ERRCONFLICT:
				slog.Info("image already exists")
				isImgExist = true
			default:
				s.serverErrorResponse(w, r, err)
				return
			}
		}

		if !isImgExist {
			// Create a new file in the local file system
			_, _ = file.Seek(0, io.SeekStart)
			dst, err := os.Create(filepath.Join(s.ImageDir, hashValue+path.Ext(fileHeader.Filename)))
			if err != nil {
				s.serverErrorResponse(w, r, err)
				return
			}
			defer dst.Close()

			if _, err := io.Copy(dst, file); err != nil {
				s.serverErrorResponse(w, r, err)
				return
			}
		}

		imagesResp[fileHeader.Filename] = imgData.ID.String()

	}

	resp := envelope{
		"images": imagesResp,
	}
	err = s.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
		return
	}

}

func (s *Server) showImageHandler(w http.ResponseWriter, r *http.Request) {
	imgID := chi.URLParam(r, "image-id")

	id, err := uuid.Parse(imgID)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}
	img, err := s.ImageService.GetImage(r.Context(), id)
	if err != nil {
		switch imageapi.ErrorCode(err) {
		case imageapi.ERRNOTFOUND:
			s.notFoundResponse(w, r)
		default:
			s.serverErrorResponse(w, r, err)
			return
		}

	}

	imagePath := filepath.Join(s.ImageDir, img.Hash+filepath.Ext(img.Name))

	// Open the file
	file, err := os.Open(imagePath)
	if err != nil {
		if os.IsNotExist(err) {
			s.notFoundResponse(w, r)
		} else {
			s.serverErrorResponse(w, r, err)
		}
		return
	}
	defer file.Close()

	// Set the appropriate headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", img.Name))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Serve the file
	http.ServeFile(w, r, imagePath)
}
