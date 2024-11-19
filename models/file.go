package models

type File struct {
	FilePath string  `json:"file_path"`
	Size     float64 `json:"size"`
	MimeType string  `json:"mimetype"`
	Content  []byte  `json:"-"`
}
