package structs

type MoMoCreatePaymentPayload struct {
	PartnerCode  string `json:"partnerCode"`
	AccessKey    string `json:"accessKey"`
	RequestID    string `json:"requestId"`
	Amount       string `json:"amount"`
	OrderID      string `json:"orderId"`
	OrderInfo    string `json:"orderInfo"`
	PartnerName  string `json:"partnerName"`
	StoreId      string `json:"storeId"`
	OrderGroupId string `json:"orderGroupId"`
	Lang         string `json:"lang"`
	AutoCapture  bool   `json:"autoCapture"`
	RedirectUrl  string `json:"redirectUrl"`
	IpnUrl       string `json:"ipnUrl"`
	ExtraData    string `json:"extraData"`
	RequestType  string `json:"requestType"`
	Signature    string `json:"signature"`
}

type MoMoCreatePaymentResponse struct {
	PartnerCode  string `json:"partnerCode"`
	RequestId    string `json:"requestId"`
	OrderId      string `json:"orderId"`
	Amount       string `json:"amount"`
	ResponseTime string `json:"responseTime"`
	Message      string `json:"message"`
	ResultCode   int    `json:"resultCode"`
	PayUrl       string `json:"payUrl"`
}

type MoMoIPNPayload struct {
	OrderType    string  `json:"orderType"`
	Amount       float64 `json:"amount"`
	PartnerCode  string  `json:"partnerCode"`
	OrderID      string  `json:"orderId"`
	ExtraData    string  `json:"extraData"`
	Signature    string  `json:"signature"`
	TransId      int64   `json:"transId"`
	ResponseTime int64   `json:"responseTime"`
	ResultCode   int     `json:"resultCode"`
	Message      string  `json:"message"`
	PayType      string  `json:"payType"`
	RequestID    string  `json:"requestId"`
	OrderInfo    string  `json:"orderInfo"`
}

type MoMoQueryPaymentPayload struct {
	PartnerCode string `json:"partnerCode"`
	RequestID   string `json:"requestId"`
	OrderID     string `json:"orderId"`
	Lang        string `json:"lang"`
	Signature   string `json:"signature"`
}

type MoMoQueryPaymentResponse struct {
	PartnerCode   string      `json:"partnerCode"`
	RequestID     string      `json:"requestId"`
	OrderID       string      `json:"orderId"`
	ExtraData     string      `json:"extraData"`
	Amount        string      `json:"amount"`
	TransID       string      `son:"transId"`
	PayType       string      `json:"payType"`
	ResultCode    int         `json:"resultCode"`
	RefundTrans   interface{} `json:"refundTrans"`
	Message       string      `json:"message"`
	ResponseTime  string      `json:"responseTime"`
	PaymentOption string      `json:"paymentOption"`
	PromotionInfo interface{} `json:"promotionInfo"`
}
