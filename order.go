package weixinpay

import "encoding/xml"

// PlaceOrderResponse represent place order reponse message from weixin pay.
// For field explanation refer to: http://pay.weixin.qq.com/wiki/doc/api/app.php?chapter=9_1
type PlaceOrderResponse struct {
	XMLName     xml.Name `xml:"xml"`
	ReturnCode  string   `xml:"return_code"`
	ReturnMsg   string   `xml:"return_msg"`
	AppId       string   `xml:"appid"`
	MchId       string   `xml:"mch_id"`
	DeviceInfo  string   `xml:"device_info"`
	NonceStr    string   `xml:"nonce_str"`
	Sign        string   `xml:"sign"`
	ResultCode  string   `xml:"result_code"`
	ErrCode     string   `xml:"err_code"`
	ErrCodeDesc string   `xml:"err_code_des"`
	TradeType   string   `xml:"trade_type"`
	PrepayId    string   `xml:"prepay_id"`
	CodeUrl     string   `xml:"code_url"`
}

// Parse the reponse message from weixin pay to struct of PlaceOrderResult
func ParsePlaceOrderResponse(data []byte) (*PlaceOrderResponse, error) {
	var resp PlaceOrderResponse
	err := xml.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (p *PlaceOrderResponse) IsSuccess() bool {
	return p.ReturnCode == "SUCCESS"
}

type CloseOrderResponse struct {
	XMLName     xml.Name `xml:"xml"`
	ReturnCode  string   `xml:"return_code"`
	ReturnMsg   string   `xml:"return_msg"`
	AppId       string   `xml:"appid"`
	MchId       string   `xml:"mch_id"`
	NonceStr    string   `xml:"nonce_str"`
	Sign        string   `xml:"sign"`
	ResultCode  string   `xml:"result_code"`
	ErrCode     string   `xml:"err_code"`
	ErrCodeDesc string   `xml:"err_code_des"`
}

func ParseCloseOrderResponse(data []byte) (*CloseOrderResponse, error) {
	var resp CloseOrderResponse
	err := xml.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}