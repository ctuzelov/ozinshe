package models

import "mime/multipart"

type SavePhoto struct {
	File_form    *multipart.Form
	UploadPath   string
	MaxImageSize int64
}
