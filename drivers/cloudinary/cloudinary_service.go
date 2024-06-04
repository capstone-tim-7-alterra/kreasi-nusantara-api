package cloudinary

import (
	"context"
	"kreasi-nusantara-api/config"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService interface {
	UploadImage(ctx context.Context, input multipart.File) (string, error)
	UploadVideo(ctx context.Context, input multipart.File) (string, error)
}

type cloudinaryService struct {
	cloudinary *cloudinary.Cloudinary
}

func NewCloudinaryService(cloudinary *cloudinary.Cloudinary) CloudinaryService {
	return &cloudinaryService{
		cloudinary: cloudinary,
	}
}

func (cs *cloudinaryService) UploadImage(ctx context.Context, input multipart.File) (string, error) {
	uploadParams := uploader.UploadParams{
		Folder: config.EnvCloudUploadFolder(),
	}

	result, err := cs.cloudinary.Upload.Upload(ctx, input, uploadParams)
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}

func (cs *cloudinaryService) UploadVideo(ctx context.Context, input multipart.File) (string, error) {
	uploadParams := uploader.UploadParams{
		Folder: config.EnvCloudUploadFolder(),
	}

	result, err := cs.cloudinary.Upload.Upload(ctx, input, uploadParams)
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}
