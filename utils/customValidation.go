package utils

import (
	"mime/multipart"
)

//IsPdf fujction to check the extension of the files
func IsPdf(header *multipart.FileHeader) bool {
	uploadedheader := header.Header.Get("Content-Type")
	if uploadedheader != "application/pdf" {
		return false
	}
	return true
}
