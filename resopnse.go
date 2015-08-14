package weixinpay

import "encoding/xml"

// PlaceOrderResult represent place order reponse message from weixin pay.
// For field explanation refer to: http://pay.weixin.qq.com/wiki/doc/api/app.php?chapter=9_1
type PlaceOrderResult struct {
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
func ParsePlaceOrderResult(resp []byte) (*PlaceOrderResult, error) {
	var placeOrderResult PlaceOrderResult
	err := xml.Unmarshal(resp, &placeOrderResult)
	if err != nil {
		return nil, err
	}

	return &placeOrderResult, nil
}

// QueryOrder Result represent query response message from weixin pay
// Refer to http://pay.weixin.qq.com/wiki/doc/api/app.php?chapter=9_2&index=4
type QueryOrderResult struct {
	XMLName        xml.Name `xml:"xml"`
	ReturnCode     string   `xml:"return_code"`
	ReturnMsg      string   `xml:"return_msg"`
	AppId          string   `xml:"appid"`
	MchId          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	ResultCode     string   `xml:"result_code"`
	ErrCode        string   `xml:"err_code"`
	ErrCodeDesc    string   `xml:"err_code_des"`
	DeviceInfo     string   `xml:"device_info"`
	OpenId         string   `xml:"open_id"`
	IsSubscribe    string   `xml:"is_subscribe"`
	TradeType      string   `xml:"trade_type"`
	TradeState     string   `xml:"trade_state"`
	TradeStateDesc string   `xml:"trade_state_desc"`
	BankType       string   `xml:"bank_type"`
	TotalFee       string   `xml:"total_fee"`
	FeeType        string   `xml:"fee_type"`
	CashFee        string   `xml:"cash_fee"`
	CashFeeType    string   `xml:"cash_fee_type"`
	CouponFee      string   `xml:"coupon_fee"`
	CouponCount    string   `xml:"coupon_count"`
	TransactionId  string   `xml:"transaction_id"`
	OrderId        string   `xml:"out_trade_no"`
	Attach         string   `xml:"attach"`
	TimeEnd        string   `xml:"time_end"`
}

func ParseQueryOrderResult(resp []byte) (*QueryOrderResult, error) {
	var queryOrderResult QueryOrderResult
	err := xml.Unmarshal(resp, &queryOrderResult)
	if err != nil {
		return nil, err
	}

	return &queryOrderResult, nil
}
