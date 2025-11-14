package helpers

import (
	"backend-city/structs"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadConfig struct {
	File           *multipart.FileHeader // File yang akan diupload
	AllowedTypes   []string              // Ekstensi file yang diperbolehkan
	MaxSize        int64                 // Ukuran maksimum file (dalam byte)
	DestinationDir string                // Folder tujuan penyimpanan
}

type UploadResult struct {
	FileName string
	FilePath string
	Error    error
	Response *structs.ErrorResponse
}

// Untuk nama file yang sudah di-slugify
func SlugifyFilename(filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	// Slugify base name saja
	slugBase := Slugify(base)

	// Gabungkan dengan ekstensi asli
	return slugBase + ext
}

func UploadFile(c *gin.Context, config UploadConfig) UploadResult {
	// Cek apakah file ada
	if config.File == nil {
		return UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "File is required",
				Errors:  map[string]string{"file": "No file was uploaded"},
			},
		}
	}

	// Cek ukuran file
	if config.File.Size > config.MaxSize {
		return UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "File to large",
				Errors:  map[string]string{"file": fmt.Sprintf("Maximum file size is: %dMB", config.MaxSize/(1<<20))},
			},
		}
	}

	// Cek tipe file
	ext := strings.ToLower(filepath.Ext(config.File.Filename))
	allowed := false
	for _, t := range config.AllowedTypes {
		if ext == t {
			allowed = true
			break
		}
	}
	if !allowed {
		return UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "Invalid file type",
				Errors:  map[string]string{"file": fmt.Sprintf("Allowed file types: %v", config.AllowedTypes)},
			},
		}
	}

	// Generate UUID sebagai nama file
	uuidName := uuid.New().String()

	// Pertahankan ekstensi asli
	filename := uuidName + ext
	filePath := filepath.Join(config.DestinationDir, filename)

	// Buat folder tujuan jika belum ada
	if err := os.MkdirAll(config.DestinationDir, 0755); err != nil {
		return UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "Failed to create upload directory",
				Errors:  map[string]string{"system": err.Error()},
			},
		}
	}

	// Simpan file ke folder tujuan
	if err := c.SaveUploadedFile(config.File, filePath); err != nil {
		return UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "Failed to save file",
				Errors:  map[string]string{"file": err.Error()},
			},
		}
	}

	// Kembalikan hasil upload
	return UploadResult{
		FileName: filename,
		FilePath: filePath,
	}
}
