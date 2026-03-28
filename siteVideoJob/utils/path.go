package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

// url域名正侧
const domainRegex string = "^(?U)(http(s)?):(.)+(com|cn|net|org|biz|info|cc|tv)"

// 拼接域名 相对路径
func BindUrl(hostname string, path ...string) string {
	if hostname == "" {
		return strings.TrimLeft(filepath.Join(path...), "/")
	}
	if len(path) == 0 {
		return hostname
	}
	pathx := &strings.Builder{}
	for _, v := range path {
		newV := strings.TrimRight(v, "/")
		if !strings.HasPrefix(newV, "/") {
			pathx.WriteString("/")
		}
		pathx.WriteString(newV)
	}
	pr := strings.TrimRight(hostname, "/")
	return pr + pathx.String()
}

// 替换资源域名
func ReplaceHost(src string, replaceHost string) string {
	if src == "" {
		return src
	}
	if !IsValidHost(src) {
		return src
	}
	reg, err := regexp.Compile(domainRegex)
	if err != nil {
		return src
	}
	return reg.ReplaceAllString(src, strings.TrimRight(replaceHost, "/"))
}

// 是否是合法的域名
func IsValidHost(hostPath string) bool {
	reg, err := regexp.Compile(domainRegex)
	if err != nil {
		return false
	}
	return reg.MatchString(hostPath)
}

/**
 *处理指定的资源路径
 * filed-待处理的路径
 * domain-给定的域名 带http开头的
 * filtersArr-不需要处理的域名白名单
 * 如果filed是绝对路径,且其域名 不在白名单内 则替换域名则使用新域名替换,否则不替换
 * 如果是相对路径,则拼接域名domain
 */

func BindOrReplacePathForNormal(filed, domain string, filtersArr []string) string {
	if filed == "" {
		return filed
	}

	spiltStr := ","
	urls := strings.Split(filed, spiltStr)
	var newFiled []string
	for _, v := range urls {
		if v == "" {
			continue
		}
		// 如果是绝对路径,则使用新域名替换
		if IsValidHost(v) {
			// 是否为白名单 不在白名单内 则替换域名
			if !ContainsAnyIgnoreCase(v, filtersArr...) {
				v = ReplaceHost(v, domain)
			}
		} else {
			// 相对路径 则拼接域名
			v = BindUrl(domain, v)
		}
		newFiled = append(newFiled, v)
	}
	return strings.Join(newFiled, spiltStr)
}
