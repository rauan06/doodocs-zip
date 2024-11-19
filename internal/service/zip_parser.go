package service

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"zip/models"
)

type ZipParser interface {
	Unzip(r io.Reader) error
	CreateResponse() (*models.Response, error)
}

type Zip struct {
	ZipReader *zip.Reader
	Size      int64
}

// Unzip reads the contents of the provided reader, extracting zip data
func (z *Zip) Unzip(r io.Reader) error {
	zipData, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("couldn't read from request body: %w", err)
	}

	z.Size = int64(len(zipData))

	zipReader, err := zip.NewReader(bytes.NewReader(zipData), z.Size)
	if err != nil {
		return fmt.Errorf("couldn't create zip reader: %w", err)
	}

	z.ZipReader = zipReader
	return nil
}

// CreateResponse generates a response based on the contents of the zip file
func (z *Zip) CreateResponse() (*models.Response, error) {
	if z.ZipReader == nil {
		return nil, fmt.Errorf("can't create response without calling Unzip() first")
	}

	var response models.Response
	var totalSize, compressedSize float64
	var totalFiles float64

	for _, file := range z.ZipReader.File {
		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("error opening file '%s' in zip: %w", file.Name, err)
		}
		defer rc.Close()

		content, err := io.ReadAll(rc)
		if err != nil {
			return nil, fmt.Errorf("couldn't read file '%s' from zip: %w", file.Name, err)
		}

		mimeType := http.DetectContentType(content[:512])

		response.Files = append(response.Files, models.File{
			FilePath: file.Name,
			Size:     float64(file.UncompressedSize64),
			MimeType: mimeType,
		})

		totalSize += float64(file.UncompressedSize64)
		compressedSize += float64(file.CompressedSize64)
		totalFiles++
	}

	response.TotalSize = totalSize
	response.TotalFiles = totalFiles
	response.ArchiveSize = compressedSize

	return &response, nil
}
