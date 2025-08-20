package constants

import "fmt"

type supportTicketStatusStruct struct {
	PENDING  int
	APPROVED int
	REJECTED int
}

type contractStatusStruct struct {
	ACTIVE                int
	EXPIRED               int
	CANCELLED             int
	WAITING_FOR_SIGNATURE int
	NOT_IN_EFFECT         int
}

type roomStatusStruct struct {
	RENTED      int
	SOLD        int
	AVAILABLE   int
	MAINTENANCE int
	UNAVAILABLE int
}

type userGenderStruct struct {
	MALE   int
	FEMALE int
	OTHER  int
}

type uploadTypeStruct struct {
	ADD_CUSTOMERS int
	ADD_CONTRACTS int
	ADD_BILLS     int
}

type importProcessResultStruct struct {
	SUCCESS int
	FAILED  int
}

type billStatusStruct struct {
	UN_PAID    int
	PAID       int
	OVERDUE    int
	PROCESSING int
	CANCELLED  int
}

type residentRelationshipStruct struct {
	SPOUSE int
	CHILD  int
	PARENT int
	OTHER  int
}

type fileUploadStruct struct {
	AllowedImageTypes  []string
	AllowedFileTypes   []string
	AllowedUploadTypes []string
	MaxImageSize       int64
	MaxFileSize        int64
	MaxUploadSize      int64
	MaxImageSizeStr    string
	MaxFileSizeStr     string
	MaxUploadSizeStr   string
}

type contractTypeStruct struct {
	RENT int
	BUY  int
}

type Notification struct {
	MarkedStatus   int
	UnmarkedStatus int
	ReadStatus     int
	UnreadStatus   int
}

type WebsocketSignalType struct {
	NewInbox       int
	NewImportant   int
	NewSent        int
	UploadCustomer int
	UploadContract int
	UploadBill     int
}

var Common = struct {
	SupportTicketStatus  supportTicketStatusStruct
	ContractType         contractTypeStruct
	ContractStatus       contractStatusStruct
	RoomStatus           roomStatusStruct
	UserGender           userGenderStruct
	UploadType           uploadTypeStruct
	BillStatus           billStatusStruct
	ResidentRelationship residentRelationshipStruct
	FileUpload           fileUploadStruct
	EmailTokenLimit      int
	NewPasswordLength    int
	Notification         Notification
	WebsocketSignalType  WebsocketSignalType
	ImportProcessResult  importProcessResultStruct
}{
	SupportTicketStatus: supportTicketStatusStruct{
		PENDING:  1,
		APPROVED: 2,
		REJECTED: 3,
	},
	ContractType: contractTypeStruct{
		RENT: 1,
		BUY:  2,
	},
	ContractStatus: contractStatusStruct{
		ACTIVE:                1,
		EXPIRED:               2,
		CANCELLED:             3,
		WAITING_FOR_SIGNATURE: 4,
		NOT_IN_EFFECT:         5,
	},
	RoomStatus: roomStatusStruct{
		RENTED:      1,
		SOLD:        2,
		AVAILABLE:   3,
		MAINTENANCE: 4,
		UNAVAILABLE: 5,
	},
	UserGender: userGenderStruct{
		MALE:   1,
		FEMALE: 2,
		OTHER:  3,
	},
	UploadType: uploadTypeStruct{
		ADD_CUSTOMERS: 1,
		ADD_CONTRACTS: 2,
		ADD_BILLS:     3,
	},
	BillStatus: billStatusStruct{
		UN_PAID:    1,
		PAID:       2,
		OVERDUE:    3,
		PROCESSING: 4,
		CANCELLED:  5,
	},
	ResidentRelationship: residentRelationshipStruct{
		SPOUSE: 1,
		CHILD:  2,
		PARENT: 3,
		OTHER:  4,
	},
	EmailTokenLimit:   5,
	NewPasswordLength: 8,
	FileUpload: fileUploadStruct{
		AllowedImageTypes: []string{"image/jpeg", "image/jpg", "image/png"},
		AllowedFileTypes: []string{
			"application/pdf",
			"image/jpeg",
			"image/jpg",
			"image/png",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		},
		AllowedUploadTypes: []string{"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.ms-excel"},
		MaxImageSize:       2 * 1024 * 1024,  // 2 MB
		MaxFileSize:        10 * 1024 * 1024, // 10 MB
		MaxUploadSize:      20 * 1024 * 1024, // 20 MB
		MaxImageSizeStr:    "2MB",
		MaxFileSizeStr:     "10MB",
		MaxUploadSizeStr:   "20MB",
	},
	Notification: Notification{
		MarkedStatus:   1,
		UnmarkedStatus: 0,
		ReadStatus:     1,
		UnreadStatus:   0,
	},
	WebsocketSignalType: WebsocketSignalType{
		NewInbox:       1,
		NewImportant:   2,
		NewSent:        3,
		UploadCustomer: 4,
		UploadContract: 5,
		UploadBill:     6,
	},
	ImportProcessResult: importProcessResultStruct{
		SUCCESS: 1,
		FAILED:  2,
	},
}

func GetRoomImageURL(imagePrefix, buildingID, roomNo, fileName string) string {
	return fmt.Sprintf("%s/buildings/%s/rooms/%s/%s", imagePrefix, buildingID, roomNo, fileName)
}

func GetBuildingImageURL(imagePrefix, buildingID, fileName string) string {
	return fmt.Sprintf("%s/buildings/%s/%s", imagePrefix, buildingID, fileName)
}

func GetUserImageURL(imagePrefix, userID, fileName string) string {
	return fmt.Sprintf("%s/users/%s/%s", imagePrefix, userID, fileName)
}

func GetContractFileURL(filePrefix, contractID, fileName string) string {
	return fmt.Sprintf("%s/contracts/%s/%s", filePrefix, contractID, fileName)
}

func GetTicketImageURL(filePrefix, ticketID, fileName string) string {
	return fmt.Sprintf("%s/tickets/%s/%s", filePrefix, ticketID, fileName)
}

func GetNotificationFileURL(filePrefix, notificationID, fileName string) string {
	return fmt.Sprintf("%s/notifications/%s/%s", filePrefix, notificationID, fileName)
}

func GetUploadFileURL(filePrefix, uploadID, fileName string) string {
	return fmt.Sprintf("%s/uploads/%s/%s", filePrefix, uploadID, fileName)
}
