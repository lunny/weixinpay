package weixinpay

import (
	"fmt"
	"crypto/md5"
	"net/url"
	"sort"
	"strconv"
	"time"
)

func SignUrl(reqUrl, key string) (string, error) {
	u, err := url.Parse(reqUrl)
	if err != nil {
		return "", err
	}

	return SignQuery(u.Query(), key), nil
}

func SignQuery(vals url.Values, key string) string {
	var params = make(Params, 0)
	params.FromQuery(vals)
	return Sign(params, key)
}

// Sign the parameter in form of map[string]string with app key.
// Empty string and "sign" key is excluded before sign.
// Please refer to http://pay.weixin.qq.com/wiki/doc/api/app.php?chapter=4_3
func Sign(params Params, key string) string {
	sort.Sort(params)
	preSignWithKey := params.ToQueryString() + "&key=" + key
	return fmt.Sprintf("%X", md5.Sum([]byte(preSignWithKey)))
}

// NewNonceString return random string in 32 characters
func NewNonceString() string {
	nonce := strconv.FormatInt(time.Now().UnixNano(), 36)
	return fmt.Sprintf("%x", md5.Sum([]byte(nonce)))
}

const ChinaTimeZoneOffset = 8 * 60 * 60 //Beijing(UTC+8:00)

// NewTimestampString return
func NewTimestampString() string {
	return fmt.Sprintf("%d", time.Now().Unix()+ChinaTimeZoneOffset)
}
