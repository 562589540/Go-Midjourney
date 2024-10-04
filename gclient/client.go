package gclient

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	once           sync.Once
	clientInstance *Client
)

// Client 是封装了 HTTP 客户端的结构体
type Client struct {
	httpClient       *http.Client
	defaultTransport *http.Transport // 默认的 Transport，没有代理
	mu               sync.Mutex      // 确保并发安全
}

func GetGclient() *Client {
	once.Do(func() {
		clientInstance = &Client{
			httpClient: &http.Client{
				Transport: &http.Transport{}, // 使用默认 Transport
				Timeout:   20 * time.Second,  // 设置超时
			},
			defaultTransport: &http.Transport{}, // 初始化默认 Transport
		}
	})
	return clientInstance

}

// GetHTTPClient 获取 HTTP 客户端的实例
func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}

// SetProxy 设置代理的方法
func (c *Client) SetProxy(proxyURL string) (*url.URL, error) {
	c.mu.Lock()         // 锁定，确保并发安全
	defer c.mu.Unlock() // 函数结束时解锁

	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	// 设置带有代理的 Transport
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
	c.httpClient.Transport = transport

	return proxy, nil
}

// ClearProxy 取消代理的方法，恢复默认的 Transport
func (c *Client) ClearProxy() {
	c.mu.Lock()         // 锁定，确保并发安全
	defer c.mu.Unlock() // 函数结束时解锁

	c.httpClient.Transport = c.defaultTransport
}
