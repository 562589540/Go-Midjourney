package services

import "os"

func OnProxy(http, https string) {
	// 设置代理
	os.Setenv("HTTP_PROXY", http)
	os.Setenv("HTTPS_PROXY", https)
}

func OffProxy() {
	// 关闭代理
	os.Unsetenv("HTTP_PROXY")
	os.Unsetenv("HTTPS_PROXY")
}
