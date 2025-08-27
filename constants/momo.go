package constants

type momoResultCode struct {
	Success                 int
	Pending                 int
	PaymentProcessorPending int
	PaymentAuthorized       int
}

type momo struct {
	CreatePaymentEndpoint string
	QueryPaymentEndPoint  string
	ResultCode            momoResultCode
}

var Momo = momo{
	CreatePaymentEndpoint: "/v2/gateway/api/create",
	QueryPaymentEndPoint:  "/v2/gateway/api/query",
	ResultCode: momoResultCode{
		Success:                 0,
		Pending:                 7000,
		PaymentProcessorPending: 7002,
		PaymentAuthorized:       9000,
	},
}
