package weixinpay

import (
	"fmt"
	"net/url"
	"reflect"
)

type Param struct {
	Key, Value string
}

type Params []Param

func (p Params) Len() int {
	return len(p)
}

func stringLess(a, b string) bool {
	var i int
	for {
		if i+1 > len(a) {
			return true
		}
		if i+1 > len(b) {
			return false
		}

		if a[i] != b[i] {
			return a[i] < b[i]
		}

		i++
	}
}

func (p Params) Less(i, j int) bool {
	return stringLess(p[i].Key, p[j].Key)
}

func (p Params) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Params) ToQueryString() (res string) {
	for _, param := range p {
		if len(param.Key) > 0 && len(param.Value) > 0 {
			res += param.Key + "=" + param.Value + "&"
		}
	}
	if len(res) > 0 {
		res = res[0 : len(res)-1]
	}
	return
}

func (p Params) FromQuery(vals url.Values) {
	for k, v := range vals {
		if k == "sign" {
			continue
		}
		if v[0] == "" {
			continue
		}
		p = append(p, Param{k, v[0]})
	}
}

func (p Params) FromUrl(reqUrl string) error {
	u, err := url.Parse(reqUrl)
	if err != nil {
		return err
	}

	p.FromQuery(u.Query())
	return nil
}

// ToXmlString convert the map[string]string to xml string
func (p Params) ToXmlString() string {
	xml := "<xml>"
	for _, s := range p {
		xml = xml + fmt.Sprintf("<%s>%s</%s>", s.Key, s.Value, s.Key)
	}
	xml = xml + "</xml>"

	return xml
}

// ToParams convert the xml struct to map[string]string
func ToParams(in interface{}) (Params, error) {
	out := make(Params, 0)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get("xml"); tagv != "" && tagv != "xml" && tagv != "sign" {
			// set key of map to value in struct field
			out = append(out, Param{tagv, v.Field(i).String()})
		}
	}
	return out, nil
}
