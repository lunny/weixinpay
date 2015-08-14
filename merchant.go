package weixinpay

import (
	"fmt"
)

type Merchant struct {
    AppId         string // 微信公众账号或开放平台APP的唯一标识
	AppKey        string // API密钥
	AppSecret string // API接口密码
	MchId         string // 微信支付商户号
}

func NewMerchant(appid, appkey, mchid, appsecret string) *Merchant {
	return &Merchant{
		AppId: appid,
		AppKey: appkey,
		MchId: mchid,
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

// 统一下单 https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_1
func (m *Merchant) PlaceOrder(orderId, desc, clientIp, notifyUrl string, amount int64) (*PlaceOrderResponse, error) {
	var params = Params{
		{"appid", m.AppId},
		{"body", desc},
		{"mch_id", m.MchId},
		{"nonce_str", NewNonceString()},
		{"notify_url", notifyUrl},
		{"out_trade_no", orderId},
		{"product_id", orderId},
		{"spbill_create_ip", clientIp},
		{"total_fee", fmt.Sprintf("%d", amount)},
		{"trade_type", "NATIVE"},
	}

	data, err := doHttpPost(PlaceOrderUrl, []byte(m.Sign(params)))
	if err != nil {
		return nil, err
	}

	return ParsePlaceOrderResponse(data)
}

// 查询订单 https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_2
func (m *Merchant) QueryOrder(transId string) (*QueryOrderResponse, error) {
	var params = Params{
		{"appid", m.AppId},
		{"mch_id", m.MchId},
		{"transaction_id", transId},
		{"nonce_str", NewNonceString()},
	}

	data, err := doHttpPost(QueryOrderUrl, []byte(m.Sign(params)))
	if err != nil {
		return nil, err
	}

	return ParseQueryOrderResponse(data)
}

// 生成二维码链接
func (m *Merchant) GenQRLink(productId string) (string, error) {
	link := fmt.Sprintf("weixin://wxpay/bizpayurl?appid=%s&mch_id=%s&product_id=%s&time_stamp=%s&nonce_str=%s",
		m.AppId, m.MchId, productId, NewTimestampString(), NewNonceString())

	var params = make(Params, 0)
	err := params.FromUrl(link)
	if err != nil {
		return "", err
	}

	return link + "&sign=" + m.Sign(params), nil
}