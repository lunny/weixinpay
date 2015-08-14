package weixinpay

import (
	"sort"
	"testing"
)

func TestParam(t *testing.T) {
	params := Params{
		{"nonce_str", "35dcf9064d9b84f971d6120f6c652ff7"},
		{"out_trade_no", "0123456"},
		{"body", "一元洗车"},
		{"bod", "hehe"},
	}

	correctParams := Params{
		{"bod", "hehe"},
		{"body", "一元洗车"},
		{"nonce_str", "35dcf9064d9b84f971d6120f6c652ff7"},
		{"out_trade_no", "0123456"},
	}

	sort.Sort(params)

	for i, v := range params {
		if correctParams[i].Key != v.Key {
			t.Fatal(i, correctParams[i].Key, "is not equal to", v)
			return
		}
	}
}
