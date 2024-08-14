package postgress

import (
	"context"

	imageapi "github.com/junaidk/image-service"
)

type StatisticsService struct {
	db *DB
}

func NewStatisticsService(db *DB) *StatisticsService {
	return &StatisticsService{db: db}
}

func (s *StatisticsService) GetStatistics(ctx context.Context) (*imageapi.Statistics, error) {
	stat := imageapi.Statistics{}
	_ = stat
	topImgFormat, err := s.getFormatStatistics(ctx)
	if err != nil {
		return nil, err
	}

	topCameraModel, err := s.getModelStatistics(ctx)
	if err != nil {
		return nil, err
	}

	topFrequencyStatistics, err := s.getFrequencyStatistics(ctx)
	if err != nil {
		return nil, err
	}

	stat.TopImageFormat = topImgFormat
	stat.TopCameraModel = topCameraModel
	stat.ImageUploadFrequency = topFrequencyStatistics

	return &stat, nil
}

// the most popular image format
func (s *StatisticsService) getFormatStatistics(ctx context.Context) (imageapi.ImageFormatItem, error) {
	row := s.db.db.QueryRowContext(ctx, `
		SELECT 
			type,
			COUNT(*) AS format_count
		FROM 
			images_metadata
		GROUP BY 
			type
		ORDER BY 
			format_count DESC
		LIMIT 1;
	`,
	)

	imgFormat := imageapi.ImageFormatItem{}
	err := row.Scan(
		&imgFormat.Type,
		&imgFormat.Count,
	)

	if err != nil {
		return imageapi.ImageFormatItem{}, err
	}
	return imgFormat, nil
}

// the top 10 most popular camera models
func (s *StatisticsService) getModelStatistics(ctx context.Context) ([]imageapi.CameraModelItem, error) {
	rows, err := s.db.db.QueryContext(ctx, `
		SELECT 
			camera_model, 
			COUNT(*) AS count
		FROM 
			images_metadata
		WHERE 
			camera_model IS NOT NULL
		GROUP BY 
			camera_model
		ORDER BY 
			count DESC
		LIMIT 10;
	`,
	)

	if err != nil {
		return nil, err
	}

	imgFormat := make([]imageapi.CameraModelItem, 0)

	for rows.Next() {
		var item imageapi.CameraModelItem
		if err := rows.Scan(
			&item.Name,
			&item.Count,
		); err != nil {
			return nil, err
		}
		imgFormat = append(imgFormat, item)
	}

	return imgFormat, nil
}

// image upload frequency per day for the past 30 days
func (s *StatisticsService) getFrequencyStatistics(ctx context.Context) ([]imageapi.FrequencyItem, error) {
	rows, err := s.db.db.QueryContext(ctx, `
		SELECT 
			DATE(created_at) AS upload_date, 
			COUNT(*) AS images_uploaded
		FROM 
			images
		WHERE 
			created_at >= NOW() - INTERVAL '30 days'
		GROUP BY 
			DATE(created_at)
		ORDER BY 
			upload_date DESC;
		`,
	)

	if err != nil {
		return nil, err
	}

	imgFormat := make([]imageapi.FrequencyItem, 0)

	for rows.Next() {
		var item imageapi.FrequencyItem
		if err := rows.Scan(
			&item.Day,
			&item.Count,
		); err != nil {
			return nil, err
		}
		imgFormat = append(imgFormat, item)
	}

	return imgFormat, nil
}
