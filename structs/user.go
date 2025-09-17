package structs

import "mime/multipart"

type NewStaff struct {
	FirstName        string `form:"firstName" validate:"required,max=255"`
	LastName         string `form:"lastName" validate:"required,max=255"`
	MiddleName       string `form:"middleName" validate:"omitempty,max=255"`
	SSN              string `form:"ssn" validate:"required,number,len=12"`
	OldSSN           string `form:"oldSSN" validate:"omitempty,number,len=9"`
	Dob              string `form:"dob" validate:"required,datetime=2006-01-02,dob_18"`
	Pob              string `form:"pob" validate:"required,max=255"`
	Phone            string `form:"phone" validate:"required,number,len=10"`
	PermanentAddress string `form:"permanentAddress" validate:"required,max=255"`
	TemporaryAddress string `form:"temporaryAddress" validate:"required,max=255"`
	Email            string `form:"email" validate:"required,max=255,email"`
	Gender           int    `form:"gender" validate:"required,min=1,max=3"`
	Schedules        []struct {
		BuildingID int64  `form:"buildingID"  validate:"required"`
		StartDate  string `form:"startDate"  validate:"required,datetime=2006-01-02"`
		EndDate    string `form:"endDate" validate:"omitempty,datetime=2006-01-02,check_date_equal_or_after=StartDate"`
	} `form:"schedules[]" validate:"dive"`
	ProfileImage  *multipart.FileHeader `validate:"required"`
	FrontSSNImage *multipart.FileHeader `validate:"required"`
	BackSSNImage  *multipart.FileHeader `validate:"required"`
}

type NewCustomer struct {
	FirstName        string                `form:"firstName" validate:"required,max=255"`
	LastName         string                `form:"lastName" validate:"required,max=255"`
	MiddleName       string                `form:"middleName" validate:"omitempty,max=255"`
	SSN              string                `form:"ssn" validate:"required,number,len=12"`
	OldSSN           string                `form:"oldSSN" validate:"omitempty,number,len=9"`
	Dob              string                `form:"dob" validate:"required,datetime=2006-01-02,dob_18"`
	Pob              string                `form:"pob" validate:"required,max=255"`
	Phone            string                `form:"phone" validate:"required,number,len=10"`
	PermanentAddress string                `form:"permanentAddress" validate:"required,max=255"`
	TemporaryAddress string                `form:"temporaryAddress" validate:"required,max=255"`
	Email            string                `form:"email" validate:"required,max=255,email"`
	Gender           int                   `form:"gender" validate:"required,min=1,max=3"`
	ProfileImage     *multipart.FileHeader `validate:"required"`
	FrontSSNImage    *multipart.FileHeader `validate:"required"`
	BackSSNImage     *multipart.FileHeader `validate:"required"`
}

type EditStaff struct {
	ID        int64 `validate:"required"`
	Schedules []struct {
		ID         int64  `form:"id" validate:"required"`
		BuildingID int64  `form:"buildingID"  validate:"required"`
		StartDate  string `form:"startDate"  validate:"required,datetime=2006-01-02"`
		EndDate    string `form:"endDate" validate:"omitempty,datetime=2006-01-02,check_date_equal_or_after=StartDate"`
	} `form:"schedules[]" validate:"dive"`
	NewSchedules []struct {
		BuildingID int64  `form:"buildingID"  validate:"required"`
		StartDate  string `form:"startDate"  validate:"required,datetime=2006-01-02"`
		EndDate    string `form:"endDate" validate:"omitempty,datetime=2006-01-02,check_date_equal_or_after=StartDate"`
	} `form:"newSchedules[]" validate:"dive"`
	DeletedSchedules []int64 `form:"deletedSchedules[]"`
}

type UpdateProfile struct {
	FirstName        string `form:"firstName" validate:"required,max=255"`
	LastName         string `form:"lastName" validate:"required,max=255"`
	MiddleName       string `form:"middleName" validate:"omitempty,max=255"`
	SSN              string `form:"ssn" validate:"required,number,len=12"`
	OldSSN           string `form:"oldSSN" validate:"omitempty,number,len=9"`
	Dob              string `form:"dob" validate:"required,datetime=2006-01-02,dob_18"`
	Pob              string `form:"pob" validate:"required,max=255"`
	Phone            string `form:"phone" validate:"required,number,len=10"`
	PermanentAddress string `form:"permanentAddress" validate:"required,max=255"`
	TemporaryAddress string `form:"temporaryAddress" validate:"required,max=255"`
	Gender           int    `form:"gender" validate:"required,min=1,max=3"`
	NewProfileImage  *multipart.FileHeader
	NewFrontSSNImage *multipart.FileHeader
	NewBackSSNImage  *multipart.FileHeader
}

type ChangePassword struct {
	OldPassword        string `json:"oldPassword" validate:"required"`
	NewPassword        string `json:"newPassword" validate:"required,max=30,password"`
	ConfirmNewPassword string `json:"confirmNewPassword" validate:"required,eqfield=NewPassword"`
}

type ChangeEmail struct {
	Password string `json:"password" validate:"required"`
	NewEmail string `json:"newEmail" validate:"required,max=255,email"`
}

type NewUploadCustomer struct {
	FirstName        string `validate:"required,max=255"`
	LastName         string `validate:"required,max=255"`
	MiddleName       string `validate:"omitempty,max=255"`
	SSN              string `validate:"required,number,len=12"`
	OldSSN           string `validate:"omitempty,number,len=9"`
	Dob              string `validate:"required,datetime=2006-01-02,dob_18"`
	Pob              string `validate:"required,max=255"`
	Phone            string `validate:"required,number,len=10"`
	PermanentAddress string `validate:"required,max=255"`
	TemporaryAddress string `validate:"required,max=255"`
	Email            string `validate:"required,max=255,email"`
	Gender           int    `validate:"required,min=1,max=3"`
	ProfileImage     string `validate:"required"`
	FrontSSNImage    string `validate:"required"`
	BackSSNImage     string `validate:"required"`
}
