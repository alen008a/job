package utils

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"io/ioutil"
)

//zlib压缩
func ZlibDo(src []byte) (string, error) {
	var out bytes.Buffer
	w := zlib.NewWriter(&out)
	_, err := w.Write(src)
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(out.Bytes()), nil
}

//进行zlib解压缩
func UnZlib(src string) (string, error) {
	decodeString, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	b := bytes.NewReader(decodeString)
	r, err := zlib.NewReader(b)
	if err = r.Close(); err != nil {
		return "", err
	}
	result, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
