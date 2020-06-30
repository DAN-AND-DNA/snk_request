package snk_request

import (
	"bytes"
	"errors"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	ERR_BAD_URL    = errors.New("bad url")
	ERR_BAD_METHOD = errors.New("bad method")
	json           = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Snk_Config struct {
	Read_timeout int
}

type Snk_request struct {
	Connect_timeout int
	Read_timeout    int
	Write_timeout   int

	http_client http.Client
}

func New() *Snk_request {
	request := &Snk_request{
		Connect_timeout: 3,
		Read_timeout:    10,
		Write_timeout:   10,
	}

	request.http_client = http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(network, addr, time.Duration(request.Connect_timeout)*time.Second)
				if err != nil {
					return nil, err
				}

				c.SetReadDeadline(time.Now().Add(time.Duration(request.Read_timeout) * time.Second))
				c.SetWriteDeadline(time.Now().Add(time.Duration(request.Write_timeout) * time.Second))

				return c, nil
			},
		},
	}

	return request
}

func New_timeout(connect_timeout, read_timeout, write_timeout int) *Snk_request {
	if connect_timeout <= 0 {
		connect_timeout = 3
	}

	if read_timeout <= 0 {
		read_timeout = 10
	}

	if write_timeout <= 0 {
		write_timeout = 10
	}

	request := &Snk_request{
		Connect_timeout: connect_timeout,
		Read_timeout:    read_timeout,
		Write_timeout:   write_timeout,
	}

	request.http_client = http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(network, addr, time.Duration(request.Connect_timeout)*time.Second)
				if err != nil {
					return nil, err
				}

				c.SetReadDeadline(time.Now().Add(time.Duration(request.Read_timeout) * time.Second))
				c.SetWriteDeadline(time.Now().Add(time.Duration(request.Write_timeout) * time.Second))

				return c, nil
			},
		},
	}

	return request
}

type Before_set struct {
	url         string
	method      string
	http_client http.Client
}

func (this *Snk_request) Get(url string) *Before_set {
	return &Before_set{
		url:         url,
		method:      "GET",
		http_client: this.http_client,
	}
}

func (this *Snk_request) Post(url string) *Before_set {
	return &Before_set{
		url:         url,
		method:      "POST",
		http_client: this.http_client,
	}
}

type Before_send struct {
	url         string
	method      string
	http_client http.Client
	headers     map[string]string
}

func (this *Before_set) Set(headers ...string) *Before_send {
	bs := &Before_send{
		url:         this.url,
		method:      this.method,
		http_client: this.http_client,
		headers:     map[string]string{},
	}

	if len(headers) == 0 {
		return bs
	}

	if len(headers)%2 != 0 {
		headers = append(headers, "")
	}

	for i := 0; i < len(headers); i += 2 {
		bs.headers[headers[i]] = headers[i+1]
	}

	return bs
}

func (this *Before_send) Send(body interface{}) *Before_end {
	be := &Before_end{
		url:         this.url,
		method:      this.method,
		http_client: this.http_client,
		req_headers: map[string]string{},
	}

	for key, vals := range this.headers {
		be.req_headers[key] = vals
	}

	if body == nil {
		be.body = nil
		return be
	}

	switch body.(type) {
	case string:
		str_body, _ := body.(string)
		be.body = []byte(str_body)

		var check_json map[string]interface{}
		if err := json.Unmarshal([]byte(str_body), &check_json); err == nil {
			be.req_headers["Content-Type"] = "application/json"
		} else {
			if str_body != "" {
				be.req_headers["Content-Type"] = "text/plain"
			} else {
				// PASS
			}
		}

	case []byte:
		byte_body, _ := body.([]byte)
		be.body = byte_body

		var check_json map[string]interface{}
		if err := json.Unmarshal(byte_body, &check_json); err == nil {
			be.req_headers["Content-Type"] = "application/json"
		} else {
			if byte_body != nil {
				be.req_headers["Content-Type"] = "text/plain"
			} else {
				// PASS
			}
		}

	default:
		json_data, err := json.Marshal(body)
		if err == nil {
			return be
		} else {
			be.req_headers["Content-Type"] = "application/json"
			be.body = json_data
		}
	}
	return be
}

type Before_end struct {
	url         string
	method      string
	http_client http.Client
	req_headers map[string]string
	Header      http.Header
	body        []byte
}

type Resp struct {
	Header http.Header
}

func (this *Before_end) End() (Resp, []byte, error) {

	resp := Resp{}

	if this.method != "GET" && this.method != "POST" {
		return resp, nil, ERR_BAD_METHOD
	}

	if this.url == "" {
		return resp, nil, ERR_BAD_URL
	}

	if _, err := url.Parse(this.url); err != nil {
		return resp, nil, ERR_BAD_URL
	}

	req, err := http.NewRequest(this.method, this.url, bytes.NewBuffer(this.body))
	if err != nil {
		return resp, nil, err
	}

	for k, v := range this.req_headers {
		req.Header.Set(k, v)
	}

	if res, err := this.http_client.Do(req); err == nil {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return resp, nil, err
		}
		resp.Header = res.Header.Clone()
		return resp, body, nil
	} else {
		return resp, nil, err
	}
	return resp, nil, err
}

func (this *Before_end) End_benchmark() (Resp, []byte, error) {

	resp := Resp{}

	if this.method != "GET" && this.method != "POST" {
		return resp, nil, ERR_BAD_METHOD
	}

	if this.url == "" {
		return resp, nil, ERR_BAD_URL
	}

	if _, err := url.Parse(this.url); err != nil {
		return resp, nil, ERR_BAD_URL
	}

	req, err := http.NewRequest(this.method, this.url, bytes.NewBuffer(this.body))
	if err != nil {
		return resp, nil, err
	}

	for k, v := range this.req_headers {
		req.Header.Set(k, v)
	}

	if res, err := this.http_client.Do(req); err == nil {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return resp, nil, err
		}
		resp.Header = res.Header.Clone()
		return resp, body, nil
	} else {
		return resp, nil, err
	}

}
