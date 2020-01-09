package soter

import (
	"context"

	"github.com/TRON-US/soter-sdk-golang/utils"
)

type AutopayRawData struct {
	Autopay   bool  `json:"autopay"`
	Timestamp int64 `json:"timestamp"`
}

type AutopayPayload struct {
	UserAddress string         `json:"user_address"`
	RawData     AutopayRawData `json:"raw_data"`
	Signature   string         `json:"signature"`
}

func getAutopayPayload(enable bool, userAddress, privateKey string) (string, error) {
	rawData := AutopayRawData{
		Autopay:   enable,
		Timestamp: utils.GetUnixTimeNow(),
	}
	rawString, err := utils.GetStructRawString(rawData)
	if err != nil {
		return "", nil
	}
	signature, err := utils.GetSignature(rawString, privateKey)
	if err != nil {
		return "", nil
	}
	payload := AutopayPayload{
		UserAddress: userAddress,
		RawData:     rawData,
		Signature:   signature,
	}
	return utils.GetStructRawString(payload)
}

func (s *Shell) Autopay(enable bool) (SoterResponse, error) {
	payload, err := getAutopayPayload(enable, s.userAddress, s.privateKey)
	if err != nil {
		return SoterResponse{}, err
	}
	var out SoterResponse
	rb := s.Request("autopay")
	rb = rb.BodyString(payload)
	rb = rb.Header("Content-Type", "application/json")
	rb.SetMethod("POST")
	err = rb.Exec(context.Background(), &out)
	return out, err
}