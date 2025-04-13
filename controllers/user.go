package controllers

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/services"
	"api/structs"
	"api/utils"
	"strconv"
	"strings"

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
	quartersStr := ctx.Query("quarters")
	var quarters []struct {
		Year     int
		Quarters []int
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)

	if err != nil {
		limit = 500
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)

	if err != nil {
		offset = 0
	}

	if quartersStr != "" {
		quartersArray := strings.Split(quartersStr, ",")
		for _, quarter := range quartersArray {
			quarterParts := strings.Split(quarter, "-")
			if len(quarterParts) == 2 {
				year, err := strconv.Atoi(quarterParts[0])
				if err != nil {
					response.Message = config.GetMessageCode("INVALID_PARAMETER")
					ctx.JSON(400, response)
					return
				}
				yearExists := false
				for _, q := range quarters {
					if q.Year == year {
						yearExists = true
						break
					}
				}

				if !yearExists {
					quarters = append(quarters, struct {
						Year     int
						Quarters []int
					}{Year: year})
				}
				quarterNum, err := strconv.Atoi(quarterParts[1])
				if err != nil {
					response.Message = config.GetMessageCode("INVALID_PARAMETER")
					ctx.JSON(400, response)
					return
				}
				for i := range quarters {
					if quarters[i].Year == year {
						quarters[i].Quarters = append(quarters[i].Quarters, quarterNum)
						break
					}
				}
			}
		}
	} else {
		currentYear := utils.GetCurrentYear()
		currentQuarter := utils.GetCurrentQuarter()

		quarters = append(quarters, struct {
			Year     int
			Quarters []int
		}{Year: currentYear, Quarters: []int{currentQuarter}})
	}

	tickets := []structs.SupportTicket{}

	if err := c.userService.GetStaffRelatedTicket(ctx, &tickets, id, limit, offset, quarters); err != nil {
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

	if err := c.userService.UpdateStaff(ctx, editStaff); err != nil {
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		ctx.JSON(500, response)
		return
	}

	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	ctx.JSON(200, response)
}
