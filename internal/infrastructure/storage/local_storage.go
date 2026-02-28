package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type LocalStorage struct {
	BaseDir  string
	MaxBytes int64
}

func NewLocalStorage(baseDir string, maxBytes int64) *LocalStorage {
	return &LocalStorage{BaseDir: baseDir, MaxBytes: maxBytes}
}

func (s *LocalStorage) SaveImage(folder string, fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size > s.MaxBytes {
		return "", fmt.Errorf("file exceeds max size of %d bytes", s.MaxBytes)
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	sniff := make([]byte, 512)
	bytesRead, err := src.Read(sniff)
	if err != nil && err != io.EOF {
		return "", err
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(sniff[:bytesRead])
	if !isAllowedMIME(mimeType) {
		return "", fmt.Errorf("unsupported content type: %s", mimeType)
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" || !isAllowedExtension(ext) {
		ext = mimeToExt(mimeType)
	}

	filename := uuid.NewString() + ext
	dir := filepath.Join(s.BaseDir, folder)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return "", err
	}

	fullPath := filepath.Join(dir, filename)
	dst, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o640)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return "/uploads/" + folder + "/" + filename, nil
}

func isAllowedMIME(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png", "image/webp", "image/avif":
		return true
	default:
		return false
	}
}

func mimeToExt(contentType string) string {
	switch contentType {
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/avif":
		return ".avif"
	default:
		return ".jpg"
	}
}

func isAllowedExtension(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".avif":
		return true
	default:
		return false
	}
}
