package structs

import "mime/multipart"

type UploadStruct struct {
	UploadType int `form:"uploadType" validate:"min=1,max=3"`
	File       *multipart.FileHeader `form:"file" validate:"required"`
}
