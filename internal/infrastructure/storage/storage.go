package storage

import "mime/multipart"

type ImageStorage interface {
	SaveImage(folder string, fileHeader *multipart.FileHeader) (string, error)
}

