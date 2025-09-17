package utils

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/structs"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
)

var (
	momoBaseURL      = ""
	momoPartnerCode  = ""
	momoAccessKey    = ""
	momoSecretKey    = ""
	momoEnv          = ""
	apmClientBaseURL = ""
)

func InitMoMoConfig() {
	momoBaseURL = config.GetEnv("MOMO_BASE_URL")
	momoPartnerCode = config.GetEnv("MOMO_PARTNER_CODE")
	momoAccessKey = config.GetEnv("MOMO_ACCESS_KEY")
	momoSecretKey = config.GetEnv("MOMO_SECRET_KEY")
	momoEnv = config.GetEnv("MOMO_ENV")
	apmClientBaseURL = config.GetEnv("APM_CLIENT_BASE_URL")
}

func CreateMoMoPayment(bill *models.BillModel, requestID, orderID uint64, momoResponse *structs.MoMoCreatePaymentResponse) error {
	if momoBaseURL == "" {
		return errors.New("MOMO_BASE_URL is not set")
	}

	if momoPartnerCode == "" {
		return errors.New("MOMO_PARTNER_CODE is not set")
	}

	if momoAccessKey == "" {
		return errors.New("MOMO_ACCESS_KEY is not set")
	}

	if momoSecretKey == "" {
		return errors.New("MOMO_SECRET_KEY is not set")
	}

	if momoEnv == "" {
		momoEnv = "UAT"
	}

	if apmClientBaseURL == "" {
		return errors.New("APM_CLIENT_BASE_URL is not set")
	}

	var (
		partnerName  = ""
		storeId      = ""
		lang         = "en"
		extraData    = ""
		orderInfo    = fmt.Sprintf("Payment for billing ID: %d", bill.ID)
		amount       = 0.0
		requestType  = "payWithMethod"
		redirectUrl  = fmt.Sprintf("%s/bill/%d", apmClientBaseURL, bill.ID)
		ipnUrl       = fmt.Sprintf("%s/api/bill/%d/momo-confirm", apmClientBaseURL, bill.ID)
		endpoint     = momoBaseURL + constants.Momo.CreatePaymentEndpoint
		orderGroupId = ""
		autoCapture  = true
	)

	if momoEnv == "production" {
		amount = bill.Amount
	} else {
		amount = 10000 // 10k VND for testing in non-production environment
	}

	roundedAmount := int64(math.Round(amount))

	// Build raw signature
	var rawSignature bytes.Buffer
	rawSignature.WriteString("accessKey=")
	rawSignature.WriteString(momoAccessKey)
	rawSignature.WriteString("&amount=")
	rawSignature.WriteString(fmt.Sprintf("%d", roundedAmount))
	rawSignature.WriteString("&extraData=")
	rawSignature.WriteString(extraData)
	rawSignature.WriteString("&ipnUrl=")
	rawSignature.WriteString(ipnUrl)
	rawSignature.WriteString("&orderId=")
	rawSignature.WriteString(fmt.Sprintf("%d", orderID))
	rawSignature.WriteString("&orderInfo=")
	rawSignature.WriteString(orderInfo)
	rawSignature.WriteString("&partnerCode=")
	rawSignature.WriteString(momoPartnerCode)
	rawSignature.WriteString("&redirectUrl=")
	rawSignature.WriteString(redirectUrl)
	rawSignature.WriteString("&requestId=")
	rawSignature.WriteString(fmt.Sprintf("%d", requestID))
	rawSignature.WriteString("&requestType=")
	rawSignature.WriteString(requestType)

	// Create a new HMAC by defining the hash type and the key (as byte array)
	hmac := hmac.New(sha256.New, []byte(momoSecretKey))

	// Write data to it
	hmac.Write(rawSignature.Bytes())
	signature := hex.EncodeToString(hmac.Sum(nil))

	payload := structs.MoMoCreatePaymentPayload{
		PartnerCode:  momoPartnerCode,
		AccessKey:    momoAccessKey,
		RequestID:    fmt.Sprintf("%d", requestID),
		Amount:       fmt.Sprintf("%d", roundedAmount),
		OrderID:      fmt.Sprintf("%d", orderID),
		OrderInfo:    orderInfo,
		PartnerName:  partnerName,
		StoreId:      storeId,
		OrderGroupId: orderGroupId,
		Lang:         lang,
		AutoCapture:  autoCapture,
		RedirectUrl:  redirectUrl,
		IpnUrl:       ipnUrl,
		ExtraData:    extraData,
		RequestType:  requestType,
		Signature:    signature,
	}

	var jsonPayload []byte
	var err error
	jsonPayload, err = json.Marshal(payload)
	if err != nil {
		return err
	}

	//send HTTP to momo endpoint
	response, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	json.NewDecoder(response.Body).Decode(momoResponse)

	if momoResponse.ResultCode != constants.Momo.ResultCode.Success {
		return errors.New(momoResponse.Message)
	}

	return nil
}

func GetMoMoPaymentStatus(bill *models.BillModel) (int, error) {
	if momoBaseURL == "" {
		return -1, errors.New("MOMO_BASE_URL is not set")
	}

	if momoPartnerCode == "" {
		return -1, errors.New("MOMO_PARTNER_CODE is not set")
	}

	if momoAccessKey == "" {
		return -1, errors.New("MOMO_ACCESS_KEY is not set")
	}

	if momoSecretKey == "" {
		return -1, errors.New("MOMO_SECRET_KEY is not set")
	}

	if momoEnv == "" {
		momoEnv = "UAT"
	}

	var (
		lang     = "en"
		endpoint = momoBaseURL + constants.Momo.QueryPaymentEndPoint
	)

	// Build raw signature
	var rawSignature bytes.Buffer
	rawSignature.WriteString("accessKey=")
	rawSignature.WriteString(momoAccessKey)
	rawSignature.WriteString("&orderId=")
	rawSignature.WriteString(bill.OrderID.String)
	rawSignature.WriteString("&partnerCode=")
	rawSignature.WriteString(momoPartnerCode)
	rawSignature.WriteString("&requestId=")
	rawSignature.WriteString(bill.RequestID.String)

	// Create a new HMAC by defining the hash type and the key (as byte array)
	hmac := hmac.New(sha256.New, []byte(momoSecretKey))

	// Write data to it
	hmac.Write(rawSignature.Bytes())
	signature := hex.EncodeToString(hmac.Sum(nil))

	payload := structs.MoMoQueryPaymentPayload{
		PartnerCode: momoPartnerCode,
		RequestID:   bill.RequestID.String,
		OrderID:     bill.OrderID.String,
		Lang:        lang,
		Signature:   signature,
	}

	var jsonPayload []byte
	var err error
	jsonPayload, err = json.Marshal(payload)
	if err != nil {
		return -1, err
	}

	//send HTTP to momo endpoint
	response, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return -1, err
	}

	momoResponse := &structs.MoMoQueryPaymentResponse{}

	json.NewDecoder(response.Body).Decode(momoResponse)

	return momoResponse.ResultCode, nil
}

func CheckIPNPayload(bill *models.BillModel, payload *structs.MoMoIPNPayload) bool {
	if bill.RequestID.String != payload.RequestID {
		return false
	}

	if bill.OrderID.String != payload.OrderID {
		return false
	}

	if momoEnv == "production" {
		if payload.Amount != bill.Amount {
			return false
		}
	} else {
		if payload.Amount != 10000 {
			return false
		}
	}

	if payload.PartnerCode != momoPartnerCode {
		return false
	}

	return true
}
