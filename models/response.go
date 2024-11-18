package models

type Response struct {
	FileName    string  `json:"filename"`
	ArchiveSize float32 `json:"archive_size"`
	TotalSize   float32 `json:"total_size"`
	TotalFiles  float32 `json:"total_files"`
	Files       []File  `json:"files"`
}
