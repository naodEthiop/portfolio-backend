package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SupabaseStorage struct {
	SupabaseURL        string
	ServiceRoleKey     string
	Bucket             string
	PublicObjectBase   string
	MaxBytes           int64
	HTTPClient         *http.Client
}

func NewSupabaseStorage(supabaseURL, serviceRoleKey, bucket, publicObjectBase string, maxBytes int64) (*SupabaseStorage, error) {
	if strings.TrimSpace(supabaseURL) == "" {
		return nil, fmt.Errorf("SUPABASE_URL is required when STORAGE_PROVIDER=supabase")
	}
	if strings.TrimSpace(serviceRoleKey) == "" {
		return nil, fmt.Errorf("SUPABASE_SERVICE_ROLE_KEY is required when STORAGE_PROVIDER=supabase")
	}
	if strings.TrimSpace(bucket) == "" {
		return nil, fmt.Errorf("SUPABASE_STORAGE_BUCKET is required when STORAGE_PROVIDER=supabase")
	}

	base := strings.TrimRight(strings.TrimSpace(supabaseURL), "/")
	publicBase := strings.TrimSpace(publicObjectBase)
	if publicBase == "" {
		publicBase = base + "/storage/v1/object/public/" + url.PathEscape(bucket)
	} else {
		publicBase = strings.TrimRight(publicBase, "/")
	}

	return &SupabaseStorage{
		SupabaseURL:      base,
		ServiceRoleKey:   strings.TrimSpace(serviceRoleKey),
		Bucket:           strings.TrimSpace(bucket),
		PublicObjectBase: publicBase,
		MaxBytes:         maxBytes,
		HTTPClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}, nil
}

func (s *SupabaseStorage) SaveImage(folder string, fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size > s.MaxBytes {
		return "", fmt.Errorf("file exceeds max size of %d bytes", s.MaxBytes)
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	sniff := make([]byte, 512)
	bytesRead, err := io.ReadFull(src, sniff)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
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
	objectPath := strings.Trim(folder, "/") + "/" + filename
	uploadURL := s.SupabaseURL + "/storage/v1/object/" + url.PathEscape(s.Bucket) + "/" + escapePathSegments(objectPath)

	req, err := http.NewRequest(http.MethodPut, uploadURL, src)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.ServiceRoleKey)
	req.Header.Set("apikey", s.ServiceRoleKey)
	req.Header.Set("Content-Type", mimeType)
	req.Header.Set("cache-control", "3600")
	req.Header.Set("x-upsert", "false")
	req.ContentLength = fileHeader.Size

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("supabase storage upload failed: %s (%d) %s", resp.Status, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return s.PublicObjectBase + "/" + objectPath, nil
}

func escapePathSegments(path string) string {
	parts := strings.Split(path, "/")
	escaped := make([]string, 0, len(parts))
	for _, p := range parts {
		if p == "" {
			continue
		}
		escaped = append(escaped, url.PathEscape(p))
	}
	return strings.Join(escaped, "/")
}
