package weixinpay

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PaymentRequest struct {
	AppId     string
	PartnerId string
	PrepayId  string
	Package   string
	NonceStr  string
	Timestamp string
	Sign      string
}

// Session is abstact of Transaction handler. With Session, we can get prepay id
type Session struct {
	PrepayId string
	CodeUrl string
	Config   *Config
}

// Initialized the Session with specific config
func NewSession(cfg *Config) (*Session, error) {
	if cfg.AppId == "" ||
		cfg.MchId == "" ||
		cfg.AppKey == "" ||
		cfg.NotifyUrl == "" ||
		cfg.QueryOrderUrl == "" ||
		cfg.PlaceOrderUrl == "" {
		return &Session{Config: cfg}, errors.New("config field canot empty string")
	}

	return &Session{Config: cfg}, nil
}

// Submit the order to weixin pay and return the prepay id if success,
// Prepay id is used for app to start a payment
// If fail, error is not nil, check error for more information
// amount is 
func (this *Session) PrepareOrder(orderId string, amount int64, desc string, clientIp string) error {
	odrInXml := this.signedOrderRequestXmlString(orderId, desc, clientIp, amount)
	resp, err := doHttpPost(this.Config.PlaceOrderUrl, []byte(odrInXml))
	if err != nil {
		return err
	}

	placeOrderResult, err := ParsePlaceOrderResult(resp)
	if err != nil {
		return err
	}

	if placeOrderResult.ReturnCode != "SUCCESS" {
		return fmt.Errorf("return code:%s, return desc:%s", placeOrderResult.ReturnCode, placeOrderResult.ReturnMsg)
	}

	if placeOrderResult.ResultCode != "SUCCESS" {
		return fmt.Errorf("resutl code:%s, result desc:%s", placeOrderResult.ErrCode, placeOrderResult.ErrCodeDesc)
	}

	this.PrepayId = placeOrderResult.PrepayId
	this.CodeUrl = placeOrderResult.CodeUrl

	return nil
}

func (this *Session) GenQRLink(productId string) (string, error) {
	link := fmt.Sprintf("weixin://wxpay/bizpayurl?appid=%s&mch_id=%s&product_id=%s&time_stamp=%s&nonce_str=%s",
		this.Config.AppId, this.Config.MchId, productId, NewTimestampString(), NewNonceString())
	sign, err := SignUrl(link, this.Config.AppKey)
	if err != nil {
		return "", err
	}
	return link + "&sign=" + sign, nil
}

// Query the order from weixin pay server by transaction id of weixin pay
func (this *Session) QueryOrder(transId string) (*QueryOrderResult, error) {
	var params = Params{
		{"appid", this.Config.AppId},
		{"mch_id", this.Config.MchId},
		{"transaction_id", transId},
		{"nonce_str", NewNonceString()},
	}

	queryXml := this.Sign(params)

	resp, err := doHttpPost(this.Config.QueryOrderUrl, []byte(queryXml))
	if err != nil {
		return nil, err
	}

	queryOrderResult, err := ParseQueryOrderResult(resp)
	if err != nil {
		return nil, err
	}

	return queryOrderResult, nil
}

// NewPaymentRequest build the payment request structure for app to start a payment.
// Return stuct of PaymentRequest, please refer to http://pay.weixin.qq.com/wiki/doc/api/app.php?chapter=9_12&index=2
func (this *Session) NewPaymentRequest(prepayId string) PaymentRequest {
	var (
		nonce = NewNonceString()
		timeStamp = NewTimestampString()
	)

	var params = Params{
		{"appid", this.Config.AppId},
		{"partnerid", this.Config.MchId},
		{"prepayid", prepayId},
		{"package", "Sign=WXPay"},
		{"noncestr", nonce},
		{"timestamp", timeStamp},
	}

	sign := Sign(params, this.Config.AppKey)

	return PaymentRequest{
		AppId:     this.Config.AppId,
		PartnerId: this.Config.MchId,
		PrepayId:  prepayId,
		Package:   "Sign=WXPay",
		NonceStr:  nonce,
		Timestamp: timeStamp,
		Sign:      sign,
	}
}

// sign and return xml
func (this *Session) Sign(params Params) string {
	sign := Sign(params, this.Config.AppKey)
	params = append(params, Param{"sign", sign})
	return params.ToXmlString()
}

func (this *Session) signedOrderRequestXmlString(orderId, desc, clientIp string, amount int64) string {
	if clientIp == "" {
		clientIp = this.Config.ServerIP
	}
	var params = Params{
		{"appid", this.Config.AppId},
		{"body", desc},
		{"mch_id", this.Config.MchId},
		{"nonce_str", NewNonceString()},
		{"notify_url", this.Config.NotifyUrl},
		{"out_trade_no", orderId},
		{"product_id", orderId},
		{"spbill_create_ip", clientIp},
		{"total_fee", fmt.Sprintf("%d", amount)},
		{"trade_type", "NATIVE"},
	}

	return this.Sign(params)
}

// doRequest post the order in xml format with a sign
func doHttpPost(targetUrl string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return []byte(""), err
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded;charset=UTF-8")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), err
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	return respData, nil
}
