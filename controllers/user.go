package controllers

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/services"
	"api/structs"
	"api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{userService: services.NewUserService()}
}

func (c *UserController) GetStaffList(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	users := &[]models.UserModel{}

	if err := c.userService.GetStaffList(ctx, users); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Data = users
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(200, response)
}

func (c *UserController) GetStaffDetail(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	user := &models.UserModel{}

	if err := c.userService.GetStaffDetail(ctx, user, id); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Data = user
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(200, response)
}

func (c *UserController) GetStaffSchedule(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	schedules := []models.ManagerScheduleModel{}

	if err := c.userService.GetStaffSchedule(ctx, &schedules, id); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Data = schedules
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(200, response)
}

func (c *UserController) GetStaffRelatedContract(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	contracts := []models.ContractModel{}

	if err := c.userService.GetStaffRelatedContract(ctx, &contracts, id); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Data = contracts
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(200, response)
}

func (c *UserController) GetStaffRelatedTicket(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	limitStr := ctx.DefaultQuery("limit", "500")
	offsetStr := ctx.DefaultQuery("offset", "0")
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	limit, err := strconv.ParseInt(limitStr, 10, 64)

	if err != nil {
		limit = 500
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)

	if err != nil {
		offset = 0
	}

	if startDate == "" {
		startDate = utils.GetFirstDayOfMonth()
	}

	if endDate == "" {
		endDate = utils.GetCurrentDate()
	}

	tickets := []structs.SupportTicket{}

	if err := c.userService.GetStaffRelatedTicket(ctx, &tickets, id, limit, offset, startDate, endDate); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Data = tickets
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(200, response)
}

func (c *UserController) DeleteStaffs(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	type deleteIDs struct {
		IDs []int64 `json:"IDs" validate:"required"`
	}

	input := &deleteIDs{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := constants.Validate.Struct(input); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	if err := c.userService.DeleteUsers(ctx, input.IDs); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Message = config.GetMessageCode("DELETE_SUCCESS")
	ctx.JSON(200, response)
}

func (c *UserController) AddStaff(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	newStaff := &structs.NewStaff{}

	if err := ctx.ShouldBind(newStaff); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	var err error
	newStaff.ProfileImage, err = ctx.FormFile("profileImage")
	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	newStaff.FrontSSNImage, err = ctx.FormFile("frontSSNImage")
	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	newStaff.BackSSNImage, err = ctx.FormFile("backSSNImage")
	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := constants.Validate.Struct(newStaff); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	if messageCode, err := c.userService.CheckDuplicateData(ctx, newStaff.Email, newStaff.SSN, newStaff.Phone, newStaff.OldSSN); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	} else if messageCode != "" {
		response.Message = messageCode
		ctx.JSON(400, response)
		return
	}

	if err := c.userService.CreateStaff(ctx, newStaff); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Message = config.GetMessageCode("CREATE_SUCCESS")
	ctx.JSON(200, response)
}

func (c *UserController) UpdateStaff(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	editStaff := &structs.EditStaff{}

	staffID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := ctx.ShouldBind(editStaff); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	editStaff.ID = staffID

	if err := constants.Validate.Struct(editStaff); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	checkUser := &models.UserModel{}

	if err := c.userService.GetStaffDetail(ctx, checkUser, staffID); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	if checkUser.ID == 0 {
		response.Message = config.GetMessageCode("USER_NOT_FOUND")
		ctx.JSON(400, response)
		return
	}

	if err := c.userService.UpdateStaff(ctx, editStaff); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	ctx.JSON(200, response)
}

func (c *UserController) GetCustomerList(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	users := []models.UserModel{}

	limitStr := ctx.DefaultQuery("limit", "500")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)

	if err != nil {
		limit = 500
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)

	if err != nil {
		offset = 0
	}

	if err := c.userService.GetCustomerList(ctx, &users, limit, offset); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Data = users
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(200, response)
}

func (c *UserController) DeleteCustomers(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	type deleteIDs struct {
		IDs []int64 `json:"IDs" validate:"required"`
	}

	input := &deleteIDs{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := constants.Validate.Struct(input); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	if err := c.userService.DeleteUsers(ctx, input.IDs); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Message = config.GetMessageCode("DELETE_SUCCESS")
	ctx.JSON(200, response)
}

func (c *UserController) GetCustomerDetail(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	user := &models.UserModel{}

	if err := c.userService.GetCustomerDetail(ctx, user, id); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Data = user
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(200, response)
}

func (c *UserController) GetCustomerContract(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	contracts := []structs.Contract{}
	if err := c.userService.GetCustomerContract(ctx, &contracts, id); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Data = contracts
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(200, response)
}

func (c *UserController) GetCustomerTicket(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	tickets := []structs.SupportTicket{}
	if err := c.userService.GetCustomerTicket(ctx, &tickets, id); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Data = tickets
	response.Message = config.GetMessageCode("GET_SUCCESS")
	ctx.JSON(200, response)
}

func (c *UserController) AddCustomer(ctx *gin.Context) {
	response := config.NewDataResponse(ctx)
	newCustomer := &structs.NewCustomer{}

	if err := ctx.ShouldBind(newCustomer); err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	var err error
	newCustomer.ProfileImage, err = ctx.FormFile("profileImage")
	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	newCustomer.FrontSSNImage, err = ctx.FormFile("frontSSNImage")
	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	newCustomer.BackSSNImage, err = ctx.FormFile("backSSNImage")
	if err != nil {
		response.Message = config.GetMessageCode("INVALID_PARAMETER")
		ctx.JSON(400, response)
		return
	}

	if err := constants.Validate.Struct(newCustomer); err != nil {
		response.Message = config.GetMessageCode("PARAMETER_VALIDATION")
		response.ValidateError = err.Error()
		ctx.JSON(400, response)
		return
	}

	if messageCode, err := c.userService.CheckDuplicateData(ctx, newCustomer.Email, newCustomer.SSN, newCustomer.Phone, newCustomer.OldSSN); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	} else if messageCode != "" {
		response.Message = messageCode
		ctx.JSON(400, response)
		return
	}

	if err := c.userService.CreateCustomer(ctx, newCustomer); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Message = config.GetMessageCode("CREATE_SUCCESS")
	ctx.JSON(200, response)
}
