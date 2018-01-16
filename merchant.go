package weixinpay

import (
	"errors"
	"fmt"

	"github.com/lunny/log"
)

var (
	SUCCESS = "SUCCESS"
	FAIL    = "FAIL"
)

type Merchant struct {
	AppId     string // 微信公众账号或开放平台APP的唯一标识
	AppKey    string // API密钥
	AppSecret string // API接口密码
	MchId     string // 微信支付商户号
}

func NewMerchant(appid, appkey, mchid, appsecret string) *Merchant {
	return &Merchant{
		AppId:     appid,
		AppKey:    appkey,
		MchId:     mchid,
		AppSecret: appsecret,
	}
}

func (m *Merchant) IsValid() bool {
	return m.AppId != "" && m.MchId != "" && m.AppKey != ""
}

// sign and return xml
func (m *Merchant) Sign(params Params) string {
	sign := Sign(params, m.AppKey)
	params = append(params, Param{"sign", sign})
	return params.ToXmlString()
}

var (
	NATIVE = "NATIVE"
	JSAPI  = "JSAPI"
	APP    = "APP"
	WAP    = "WAP"
)

// 统一下单 https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_1
func (m *Merchant) PlaceOrder(orderId, goodsname, desc, clientIp, notifyUrl string, amount int64, tradeType string) (*PlaceOrderResponse, error) {
	var params = Params{
		{"appid", m.AppId},
		{"body", goodsname},
		{"detail", desc},
		{"mch_id", m.MchId},
		{"nonce_str", NewNonceString()},
		{"notify_url", notifyUrl},
		{"out_trade_no", orderId},
		{"product_id", orderId},
		{"spbill_create_ip", clientIp},
		{"total_fee", fmt.Sprintf("%d", amount)},
		{"trade_type", tradeType},
	}

	postData := []byte(m.Sign(params))
	log.Debug(string(postData))
	data, err := doHttpPost(PlaceOrderUrl, postData)
	if err != nil {
		return nil, err
	}

	log.Debug(string(data))

	resp, err := ParsePlaceOrderResponse(data)
	if err != nil {
		return nil, err
	}

	log.Debug(resp)

	if resp.IsSuccess() {
		ok, err := Verify(resp, m.AppKey, resp.Sign)
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, errors.New("signature error")
		}
	}
	return resp, nil
}

// 统一下单 https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_1 JSAPI
func (m *Merchant) PlaceOrderJSAPI(orderId, goodsname, desc, clientIp, notifyUrl string, amount int64, openID string) (*PlaceOrderResponse, error) {
	var params = Params{
		{"appid", m.AppId},
		{"body", goodsname},
		{"detail", desc},
		{"mch_id", m.MchId},
		{"nonce_str", NewNonceString()},
		{"notify_url", notifyUrl},
		{"out_trade_no", orderId},
		{"product_id", orderId},
		{"spbill_create_ip", clientIp},
		{"total_fee", fmt.Sprintf("%d", amount)},
		{"openid", openID},
		{"trade_type", "JSAPI"},
	}

	postData := []byte(m.Sign(params))
	log.Debug(string(postData))
	data, err := doHttpPost(PlaceOrderUrl, postData)
	if err != nil {
		return nil, err
	}

	log.Debug(string(data))

	resp, err := ParsePlaceOrderResponse(data)
	if err != nil {
		return nil, err
	}

	log.Debug(resp)

	if resp.IsSuccess() {
		ok, err := Verify(resp, m.AppKey, resp.Sign)
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, errors.New("signature error")
		}
	}
	return resp, nil
}

func (m *Merchant) CloseOrder(orderId string) (*CloseOrderResponse, error) {
	var params = Params{
		{"appid", m.AppId},
		{"mch_id", m.MchId},
		{"out_trade_no", orderId},
		{"nonce_str", NewNonceString()},
	}

	data, err := doHttpPost(CloseOrderUrl, []byte(m.Sign(params)))
	if err != nil {
		return nil, err
	}

	return ParseCloseOrderResponse(data)
}

// 根据微信支付订单号查询订单 https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_2
func (m *Merchant) QueryOrderByTransId(transId string) (*PayResult, error) {
	var params = Params{
		{"appid", m.AppId},
		{"mch_id", m.MchId},
		{"nonce_str", NewNonceString()},
		{"transaction_id", transId},
	}

	data, err := doHttpPost(QueryOrderUrl, []byte(m.Sign(params)))
	if err != nil {
		return nil, err
	}

	return ParsePayResult(data)
}

// 根据商户订单号查询订单 https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_2
func (m *Merchant) QueryOrderByOrderId(orderId string) ([]byte, *PayResult, error) {
	var params = Params{
		{"appid", m.AppId},
		{"mch_id", m.MchId},
		{"nonce_str", NewNonceString()},
		{"out_trade_no", orderId},
	}

	data, err := doHttpPost(QueryOrderUrl, []byte(m.Sign(params)))
	if err != nil {
		return nil, nil, err
	}

	res, err := ParsePayResult(data)
	if err != nil {
		return nil, nil, err
	}
	return data, res, nil
}

// 生成二维码链接
// weixin：//wxpay/bizpayurl?sign=XXXXX&appid=XXXXX&mch_id=XXXXX&product_id=XXXXXX&time_stamp=XXXXXX&nonce_str=XXXXX
func (m *Merchant) GenQRLink(productId string) string {
	var params = Params{
		{"appid", m.AppId},
		{"mch_id", m.MchId},
		{"nonce_str", NewNonceString()},
		{"product_id", productId},
		{"time_stamp", NewTimestampString()},
	}

	sign := Sign(params, m.AppKey)
	params = append(params, Param{"sign", sign})
	fmt.Println(params)
	return fmt.Sprintf("weixin://wxpay/bizpayurl?%s", params.ToQueryString())
}

func (m *Merchant) NewScanResponse(returnCode, returnMsg, prepayId, resultCode, errCodeDes string) *ScanResponse {
	return &ScanResponse{
		Params: Params{
			{"return_code", returnCode},
			{"return_msg", returnMsg},
			{"appid", m.AppId},
			{"mch_id", m.MchId},
			{"nonce_str", NewNonceString()},
			{"prepay_id", prepayId},
			{"result_code", resultCode},
			{"err_code_des", errCodeDes},
		},
		AppKey: m.AppKey,
	}
}
