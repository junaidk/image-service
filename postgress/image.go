package postgress

import (
	"context"

	"github.com/google/uuid"
	imageapi "github.com/junaidk/image-service"
)

type ImageService struct {
	db *DB
}

func NewImageService(db *DB) *ImageService {
	return &ImageService{db: db}
}

func (s *ImageService) CreateImage(ctx context.Context, img *imageapi.Image) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
	INSERT INTO images_metadata (
		hash,
		type,
		size_width,
		size_height,
		camera_model,
		lens_model,
		location_lat,
		location_long
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	ON CONFLICT (hash)
	DO NOTHING 
	`,
		img.Hash,
		img.MetaData.Type,
		img.MetaData.Size.Width,
		img.MetaData.Size.Height,
		NewNullString(img.MetaData.CameraModel),
		NewNullString(img.MetaData.LensModel),
		img.MetaData.Location.Latitude,
		img.MetaData.Location.Longitude,
	)
	if err != nil {
		return FormatError(err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO images (
			id,
			name,
			hash
		)
		VALUES ($1,$2,$3)
		RETURNING id
		`,
		img.ID,
		img.Name,
		img.Hash,
	)
	if err != nil {
		return FormatError(err)
	}

	return tx.Commit()
}

func (s *ImageService) GetImage(ctx context.Context, id uuid.UUID) (*imageapi.Image, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `
	SELECT
		id,
		name,
		hash,
		created_at
	FROM images
	WHERE id=$1
	LIMIT 1
	`,
		id,
	)

	img := imageapi.Image{}
	err = row.Scan(
		&img.ID,
		&img.Name,
		&img.Hash,
		&img.CreateAt,
	)

	if err != nil {
		return nil, FormatError(err)
	}

	return &img, nil
}
