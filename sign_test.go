package weixinpay

import (
	"testing"
)

func TestSign(t *testing.T) {
	var params = Params{
		{"appid", "wxd930ea5d5a258f4f"},
		{"mch_id",	"10000100"},
		{"device_info",	"1000"},
		{"body", "test"},
		{"nonce_str", "ibuaiVcKdpRxkhJA"},
	}

	correctSign := "192006250b4c09247ec02edce69f6a2d"
	correctXml := `<xml><appid>wxd930ea5d5a258f4f</appid><body>test</body><device_info>1000</device_info><mch_id>10000100</mch_id><nonce_str>ibuaiVcKdpRxkhJA</nonce_str><sign>9A0A8659F005D6984697E2CA0A9CF3B7</sign></xml>`
	sign := Sign(params, correctSign)
	if sign != "9A0A8659F005D6984697E2CA0A9CF3B7" {
		t.Fatal("sign error:", sign, "is not equal to", correctSign)
		return
	}

	params = append(params, Param{"sign", sign})
	if params.ToXmlString() != correctXml {
		t.Fatal("sign error:", params.ToXmlString(), "is not equal to", correctXml)
		return
	}
}