package models

type Response struct {
	FileName    string
	ArchiveSize float32
	TotalSize   float32
	TotalFiles  float32
	Files       []File
}
