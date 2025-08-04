package structs

import "mime/multipart"

type NewStaff struct {
	FirstName        string `form:"firstName" validate:"required"`
	LastName         string `form:"lastName" validate:"required"`
	MiddleName       string `form:"middleName"`
	SSN              string `form:"ssn" validate:"required,alphanum,len=12"`
	OldSSN           string `form:"oldSSN" validate:"omitempty,alphanum,len=9"`
	Dob              string `form:"dob" validate:"required,datetime=2006-01-02,dob_18"`
	Pob              string `form:"pob" validate:"required"`
	Phone            string `form:"phone" validate:"required,alphanum,len=10"`
	PermanentAddress string `form:"permanentAddress" validate:"required"`
	TemporaryAddress string `form:"temporaryAddress" validate:"required"`
	Email            string `form:"email" validate:"required,email"`
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
	FirstName        string                `form:"firstName" validate:"required"`
	LastName         string                `form:"lastName" validate:"required"`
	MiddleName       string                `form:"middleName"`
	SSN              string                `form:"ssn" validate:"required,alphanum,len=12"`
	OldSSN           string                `form:"oldSSN" validate:"omitempty,alphanum,len=9"`
	Dob              string                `form:"dob" validate:"required,datetime=2006-01-02,dob_18"`
	Pob              string                `form:"pob" validate:"required"`
	Phone            string                `form:"phone" validate:"required,alphanum,len=10"`
	PermanentAddress string                `form:"permanentAddress" validate:"required"`
	TemporaryAddress string                `form:"temporaryAddress" validate:"required"`
	Email            string                `form:"email" validate:"required,email"`
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
	FirstName        string `form:"firstName" validate:"required"`
	LastName         string `form:"lastName" validate:"required"`
	MiddleName       string `form:"middleName"`
	SSN              string `form:"ssn" validate:"required,alphanum,len=12"`
	OldSSN           string `form:"oldSSN" validate:"omitempty,alphanum,len=9"`
	Dob              string `form:"dob" validate:"required,datetime=2006-01-02,dob_18"`
	Pob              string `form:"pob" validate:"required"`
	Phone            string `form:"phone" validate:"required,alphanum,len=10"`
	PermanentAddress string `form:"permanentAddress" validate:"required"`
	TemporaryAddress string `form:"temporaryAddress" validate:"required"`
	Gender           int    `form:"gender" validate:"required,min=1,max=3"`
	NewProfileImage  *multipart.FileHeader
	NewFrontSSNImage *multipart.FileHeader
	NewBackSSNImage  *multipart.FileHeader
}
