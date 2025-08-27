package constants

type momoResultCode struct {
	Success                 int
	Pending                 int
	PaymentProcessorPending int
	PaymentConfirm          int
}

type momo struct {
	CreatePaymentEndpoint string
	ResultCode            momoResultCode
}

var Momo = momo{
	CreatePaymentEndpoint: "/v2/gateway/api/create",
	ResultCode: momoResultCode{
		Success:                 0,
		Pending:                 7000,
		PaymentProcessorPending: 7002,
		PaymentConfirm:          9000,
	},
}
