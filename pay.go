package weixinpay

import (
	"encoding/xml"
	"io"
)

var (
	tradeStates = map[string]string{
		"SUCCESS":    "支付成功",
		"REFUND":     "转入退款",
		"NOTPAY":     "未支付",
		"CLOSED":     "已关闭",
		"REVOKED":    "已撤销（刷卡支付）",
		"USERPAYING": "用户支付中",
		"PAYERROR":   "支付失败(其他原因，如银行返回失败)",
	}
)

func GetTradeState(s string) string {
	return tradeStates[s]
}

type PayResult struct {
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
	TotalFee       int64    `xml:"total_fee"`
	FeeType        string   `xml:"fee_type"`
	CashFee        int64    `xml:"cash_fee"`
	CashFeeType    string   `xml:"cash_fee_type"`
	CouponFee      int64    `xml:"coupon_fee"`
	CouponCount    int      `xml:"coupon_count"`
	// TODO:
	//coupon_id_$n
	//coupon_fee_$n
	TransactionId string `xml:"transaction_id"`
	OutTradeNO    string `xml:"out_trade_no"`
	ProductId     string `xml:"product_id"`
	Attach        string `xml:"attach"`   // 商家数据包
	TimeEnd       string `xml:"time_end"` // 支付完成时间 yyyyMMddHHmmss
}

func (r *PayResult) IsSuccess() bool {
	return r.ReturnCode == "SUCCESS"
}

func ParsePayResult(data []byte) (*PayResult, error) {
	var result PayResult
	err := xml.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type NotifyResponse struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"` // SUCCESS/FAIL
	ReturnMsg  string   `xml:"return_msg"`
}

type ScanResponse struct {
	Params
	AppKey string
}

func (s *ScanResponse) WriteTo(w io.Writer) (int, error) {
	sign := Sign(s.Params, s.AppKey)
	s.Params = append(s.Params, Param{"sign", sign})
	return w.Write([]byte(s.Params.ToXmlString()))
}
