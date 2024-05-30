package config

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/joho/godotenv"
)

func SetupCloudinary() (*cloudinary.Cloudinary, error) {
	cldSecret := os.Getenv("CLOUDINARY_API_SECRET")
	cldName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	cldKey := os.Getenv("CLOUDINARY_API_KEY")

	cld, err := cloudinary.NewFromParams(cldName, cldKey, cldSecret)
	if err != nil {
		return nil, err
	}

	return cld, nil
}

func EnvCloudUploadFolder() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
}