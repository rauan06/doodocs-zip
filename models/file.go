package models

type File struct {
	FilePath string  `json:"file_path"`
	Size     float32 `json:"size"`
	MimeType string  `json:"mimetype"`
}
