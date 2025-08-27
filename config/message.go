package config

import "github.com/gin-gonic/gin"

type DataResponse struct {
	Message       string      `json:"message"`
	Data          interface{} `json:"data"`
	ValidateError interface{} `json:"validateError"`
	JWTToken      string      `json:"jwtToken"`
	RefreshToken  string      `json:"refreshToken"`
}

func NewDataResponse(ctx *gin.Context) *DataResponse {
	// jwt, _ := ctx.Get("jwt")

	// if jwt == nil {
	// 	jwt = ""
	// }

	jwt := ctx.GetString("jwt")

	refreshToken := ctx.GetString("refreshToken")

	return &DataResponse{
		Message:       "",
		Data:          nil,
		ValidateError: nil,
		// JWTToken:      jwt.(string),
		JWTToken:     jwt,
		RefreshToken: refreshToken,
	}
}

var messageList = map[string]string{
	// System error messages
	"SYSTEM_ERROR":      "MSG_S0001",
	"CONCURRENCY_ERROR": "MSG_S0002",

	// User error messages
	"NO_IMPORT_FILE":         "MSG_E0001",
	"INVALID_IMPORT":         "MSG_E0002",
	"INVALID_CREDENTIALS":    "MSG_E0003",
	"EMAIL_NOT_VERIFIED":     "MSG_E0004",
	"TOKEN_VERIFY_FAILED":    "MSG_E0005",
	"TOKEN_REFRESH_FAILED":   "MSG_E0006",
	"PERMISSION_DENIED":      "MSG_E0007",
	"USER_NOT_FOUND":         "MSG_E0008",
	"REQUEST_SPAM":           "MSG_E0009",
	"IMAGE_NOT_FOUND":        "MSG_E0010",
	"DATA_NOT_FOUND":         "MSG_E0011",
	"EMAIL_ALREADY_EXISTS":   "MSG_E0012",
	"SSN_ALREADY_EXISTS":     "MSG_E0013",
	"PHONE_ALREADY_EXISTS":   "MSG_E0014",
	"OLD_SSN_ALREADY_EXISTS": "MSG_E0015",
	"UPDATE_FAILED":          "MSG_E0016",
	"CREATE_FAILED":          "MSG_E0017",
	"PASSWORD_INCORRECT":     "MSG_E0018",
	"FILE_NOT_FOUND":         "MSG_E0019",
	"BILL_ALREADY_PAID":      "MSG_E0020",
	"IPN_PAYLOAD_INVALID":    "MSG_E0021",

	// Parameter error messages
	"INVALID_PARAMETER":    "MSG_V0001",
	"PARAMETER_VALIDATION": "MSG_V0002",

	// Success messages
	"GET_SUCCESS":         "MSG_I0001",
	"CREATE_SUCCESS":      "MSG_I0002",
	"UPDATE_SUCCESS":      "MSG_I0003",
	"DELETE_SUCCESS":      "MSG_I0004",
	"LOGIN_SUCCESS":       "MSG_I0005",
	"LOGOUT_SUCCESS":      "MSG_I0006",
	"VERIFY_SUCCESS":      "MSG_I0007",
	"EMAIL_SENT":          "MSG_I0008",
	"PASSWORD_RESET":      "MSG_I0009",
	"PAYMENT_INITIALIZED": "MSG_I0010",
	"PAYMENT_COMPLETED":   "MSG_I0011",
}

func GetMessageCode(key string) string {
	var message = ""
	if msg, ok := messageList[key]; ok {
		message = msg
	}

	return message
}
