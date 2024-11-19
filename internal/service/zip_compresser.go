package service

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log/slog"
	"zip/models"
)

type ZipCompresser interface {
	CreateZip(files []models.File) (*bytes.Buffer, error)
}

func (z *Zip) CreateZip(files []models.File) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	defer closeZipWriter(zipWriter)

	for _, file := range files {
		if err := writeFileToZip(zipWriter, file); err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func writeFileToZip(zipWriter *zip.Writer, file models.File) error {
	fileHeader := zip.FileHeader{
		Name:   file.FilePath,
		Method: zip.Deflate, // Compression method
	}

	zipEntry, err := zipWriter.CreateHeader(&fileHeader)
	if err != nil {
		return fmt.Errorf("couldn't create zip entry for '%s': %w", file.FilePath, err)
	}

	if _, err := zipEntry.Write(file.Content); err != nil {
		return fmt.Errorf("couldn't write content to zip entry for '%s': %w", file.FilePath, err)
	}

	return nil
}

func closeZipWriter(zipWriter *zip.Writer) {
	if err := zipWriter.Close(); err != nil {
		slog.Debug(fmt.Sprintf("error closing zip writer: %v\n", err))
	}
}
