package httpclient

import (
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/imroc/req"
	"io/ioutil"
	"net/http"
	"siteLetterJob/mdata"
	"strings"
	"time"
)

var (
	defaultHttpClient *http.Client
)

// init HTTPClient  默认开启长链接(http 1.1之后) 开启http keepalive功能，也即是否重用连接，
func init() {
	defaultHttpClient = newClient(true)
}

// Basic Authentication
type BasicAuth struct {
	Username string
	Password string
}

//httpClient 非必须，为nil 将使用默认提供的httpclient
func PostV2(path string, postData []byte, header map[string]string, httpClient *http.Client) ([]byte, error) {
	var s []byte
	if httpClient == nil {
		return []byte(""), errors.New("请传入httpClient")
	}
	payload := bytes.NewReader(postData)
	req, err := http.NewRequest("POST", path, payload)
	if err != nil {
		return s, err
	}
	for key, value := range header {
		req.Header.Add(key, value)
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return s, err
	}
	if res != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return s, err
	}
	return body, nil
}

func ProxyPostJson(path, body string, header map[string]string) ([]byte, error) {
	var s []byte
	myClient := GetShortProxyNotifyClient(time.Second * 30)

	request, _ := http.NewRequest("POST", path, strings.NewReader(body))
	request.Header.Add("content-type", "application/json")
	if header != nil {
		for key, value := range header {
			request.Header.Add(key, value)
		}
	}
	res, err := myClient.Do(request)
	if err != nil {
		return nil, err
	}
	if res.Status == "404 " {
		return nil, err
	}
	defer res.Body.Close()
	body1, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return s, err
	}
	return body1, nil
}

func Post(path string, header req.Header, body interface{}) ([]byte, error) {

	resp, err := req.Post(path, header, body)
	if err != nil {
		return nil, err
	}

	bs, err := resp.ToBytes()
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func POST(path string, data []byte, header map[string]string, basicAuth ...BasicAuth) ([]byte, error) {
	isJson := mdata.Cjson.Valid(data)
	if !isJson {
		return []byte(""), errors.New("请求字符串非json格式！")
	}

	payload := strings.NewReader(string(data))
	req, _ := http.NewRequest("POST", path, payload)
	req.Header.Add("content-type", "application/json")

	if len(basicAuth) > 0 {
		if basicAuth[0].Username != "" && basicAuth[0].Password != "" {
			req.SetBasicAuth(basicAuth[0].Username, basicAuth[0].Password)
		}
	}

	for key, value := range header {
		req.Header.Add(key, value)
	}

	rsp, err := HttpClient.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return []byte(""), err
	}

	if rsp.StatusCode == 200 {
		return body, nil
	}

	return []byte(""), nil
}

func GET(path string, header map[string]string, basicAuth ...BasicAuth) ([]byte, error) {

	req, _ := http.NewRequest("GET", path, nil)

	for key, value := range header {
		req.Header.Add(key, value)
	}

	if len(basicAuth) > 0 {
		if basicAuth[0].Username != "" && basicAuth[0].Password != "" {
			req.SetBasicAuth(basicAuth[0].Username, basicAuth[0].Password)
		}
	}

	rsp, err := HttpClient.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return []byte(""), err
	}

	return body, nil
}
func POSTJson(path string, data []byte, header map[string]string, cli *http.Client) ([]byte, error) {
	if cli == nil {
		cli = HttpClient
	}

	payload := bytes.NewReader(data)
	req, err := http.NewRequest("POST", path, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")

	for key, value := range header {
		req.Header.Add(key, value)
	}

	rsp, err := cli.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode == 200 {
		switch rsp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err := gzip.NewReader(rsp.Body)
			if err != nil {
				return []byte(""), err
			}
			data, err = ioutil.ReadAll(reader)
			if err != nil {
				return []byte(""), err
			}
			return data, nil
		default:
			bodyByte, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				return []byte(""), err
			}
			return bodyByte, nil
		}
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return []byte(""), err
	}

	return body, nil
}
