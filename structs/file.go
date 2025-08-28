package structs

import "net/textproto"

type CustomFileStruct struct {
	Filename string
	Header   textproto.MIMEHeader
	Size     int64
	Content  []byte
}

type ImageValidation struct {
	Type string `validate:"image_type"`
	Size int64  `validate:"image_size"`
}

type FileValidation struct {
	Type string `validate:"file_type"`
	Size int64  `validate:"file_size"`
}

type UploadValidation struct {
	Type string `validate:"upload_type"`
	Size int64  `validate:"upload_size"`
}

type ValidateEditBuildingFile struct {
	NewBuildingImages []ImageValidation      `validate:"dive"`
	Rooms             []ValidateEditRoomFile `validate:"dive"`
	NewRooms          []ValidateAddRoomFile  `validate:"dive"`
}

type ValidateEditRoomFile struct {
	NewImages []ImageValidation `validate:"dive"`
}

type ValidateEditRoomFile2 struct {
	NewRoomImages []ImageValidation `validate:"dive"`
}

type ValidateAddBuildingFile struct {
	Images []ImageValidation     `validate:"dive"`
	Rooms  []ValidateAddRoomFile `validate:"dive"`
}

type ValidateAddRoomFile struct {
	Images []ImageValidation `validate:"dive"`
}

type ValidateUserFile struct {
	ProfileImage  ImageValidation
	FrontSSNImage ImageValidation
	BackSSNImage  ImageValidation
}

type ValidateEditContractFile struct {
	NewFiles []FileValidation `validate:"dive"`
}

type ValidateUserFile2 struct {
	ProfileImage  *ImageValidation
	FrontSSNImage *ImageValidation
	BackSSNImage  *ImageValidation
}

type ValidateSupportTicketFile struct {
	Images []ImageValidation `validate:"dive"`
}

type ValidateNoticeFile struct {
	Files []FileValidation `validate:"dive"`
}

type ValidateUploadFile struct {
	File UploadValidation
}
