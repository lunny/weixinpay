package weixinpay

import (
	"encoding/xml"
)

type NotifyRequest struct {
	XMLName     xml.Name `xml:"xml"`
	ReturnCode  string   `xml:"return_code"` // SUCCESS/FAIL
	ReturnMsg   string   `xml:"return_msg"`
	AppId string `xml:"appid"`
	MchId string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr string `xml:"nonce_str"`
	Sign string `xml:"sign"`
	ResultCode  string   `xml:"result_code"` // SUCCESS/FAIL
	ErrCode     string   `xml:"err_code"`
	ErrCodeDesc string   `xml:"err_code_des"`
	OpenId string `xml:"openid"`
	IsSubscribe string `xml:"is_subscribe"` // Y/N
	TradeType string `xml:"trade_type"`
	BankType string `xml:"bank_type"`
	TotalFee int64 `xml:"total_fee"` // 订单总金额，单位为分
	FeeType string `xml:"fee_type"`
	CashFee int64 `xml:"cash_fee"`
	CashFeeType string `xml:"cash_fee_type"`
	CouponFee  int64 `xml:"coupon_fee"`
	CouponCount int `xml:"coupon_count"`
	//coupon_id_$n
	//coupon_fee_$n
	TransactionId string `xml:"transaction_id"`
	OutTradeNO string `xml:"out_trade_no"`
	Attach string `xml:"attach"` // 商家数据包
	TimeEnd string `xml:"time_end"` // 支付完成时间 yyyyMMddHHmmss
}

func ParseNotifyRequest(data []byte) (*NotifyRequest, error) {
	var notifyReq NotifyRequest
	err := xml.Unmarshal(data, &notifyReq)
	if err != nil {
		return nil, err
	}

	return &notifyReq, nil
}

type NotifyResponse struct {
	XMLName     xml.Name `xml:"xml"`
	ReturnCode  string   `xml:"return_code"` // SUCCESS/FAIL
	ReturnMsg   string   `xml:"return_msg"`
}