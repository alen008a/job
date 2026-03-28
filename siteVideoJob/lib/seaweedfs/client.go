package seaweedfs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"sync"

	"siteVideoJob/internal/glog"
	"siteVideoJob/lib/httpclient"
	"siteVideoJob/mdata"
)

type Client struct {
	config *Config
}

type Config struct {
	Env    string // 环境
	Host   string // ip
	Port   string // 端口
	Volume string // 上传目录
	Domain string // 域名
}

var gg *Client
var once sync.Once

//func NewClient(conf *config.SeaWeedsConfig) *Client {
//	var env = config.GetEnv()
//	if env == "dev" {
//		glog.Infof("seaweeds上传服务获取配置 -- host=%v -- post=%v -- domain=%v -- volume=%v",
//			conf.SeaWeedsHost,
//			conf.SeaWeedsPort,
//			conf.SeaWeedDomain,
//			conf.SeaWeedUploadVolume,
//		)
//	}
//
//	once.Do(func() {
//		gg = new(Client)
//		gg.config = &Config{
//			Env:    env,
//			Host:   conf.SeaWeedsHost,
//			Port:   conf.SeaWeedsPort,
//			Domain: conf.SeaWeedDomain,
//			Volume: conf.SeaWeedUploadVolume,
//		}
//	})
//
//	return gg
//}

/**
 * 单个图片文件上传
 */
func (t *Client) UploadFile(
	filename, mineType string,
	file io.Reader,
) (r *UploadResp, err error) {

	formData, contentType, err := t.makeFormData(filename, mineType, file)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%v:%v%v", t.config.Host, t.config.Port, t.config.Volume)

	if t.config.Env == "local" {
		url = fmt.Sprintf("%v%v", t.config.Domain, t.config.Volume)
	}

	glog.Infof(">>> seaweedfs -- url:%v", url)

	resp, err := t.Upload(url, contentType, formData)
	if err != nil {
		return nil, err
	}

	resp.FileUrl = url
	return resp, nil
}

func (t *Client) makeFormData(
	filename,
	mimeType string,
	context io.Reader,
) (formData io.Reader, contentType string, err error) {

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, err := t.CreateFormFile(writer, "file", filename, mimeType)
	if err != nil {
		glog.Errorf(">>>Client->makeFormData->createFormFile -- %v", err)
		return
	}
	_, err = io.Copy(part, context)
	if err != nil {
		glog.Errorf(">>> Client->makeFormData->Copy -- err=%v", err)
		return
	}
	formData = buf
	contentType = writer.FormDataContentType()
	_ = writer.Close()
	return
}

func (t *Client) CreateFormFile(writer *multipart.Writer, fieldname, filename, mime string) (io.Writer, error) {
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			t.EscapeQuotes(fieldname), t.EscapeQuotes(filename)))

	if len(mime) == 0 {
		mime = "application/octet-stream"
	}
	header.Set("Content-Type", mime)
	return writer.CreatePart(header)
}

/**
 * 上传post请求处理
 */
func (t *Client) Upload(
	url, contentType string,
	formData io.Reader) (r *UploadResp, err error) {
	// 连接 增加代理
	var client = httpclient.HttpClient
	resp, err := client.Post(url, contentType, formData)
	if err != nil {
		glog.Errorf(">>> Client->upload -- err=%v", err)
		return nil, err
	}
	defer resp.Body.Close()

	uploadResp := new(UploadResp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf(">>> Client->upload->ReadAll -- err=%v", err)
		return nil, err
	}
	err = mdata.Cjson.Unmarshal(body, &uploadResp)
	if err != nil {
		glog.Errorf(">>> err=%v", err)
		return nil, err
	}

	// glog.Errorf(">>> resp=%v -- uploadResp=%v", resp, uploadResp)
	return uploadResp, nil
}

func (t *Client) EscapeQuotes(s string) string {
	if s != "" {
		return quoteEscaper.Replace(s)
	}
	return s
}
