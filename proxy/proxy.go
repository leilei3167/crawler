package proxy

import (
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
)

// ProxyFunc 满足Transport的字段要求
type ProxyFunc func(r *http.Request) (*url.URL, error)

func RoundRobinProxySwitcher(proxyURLS ...string) (ProxyFunc, error) {
	if len(proxyURLS) < 1 {
		return nil, errors.New("proxy URL is empty")
	}

	urls := make([]*url.URL, len(proxyURLS))
	for i, u := range proxyURLS {
		parsedU, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		urls[i] = parsedU
	}
	return (&roundRobinSwitcher{
		proxyURLs: urls,
		index:     0,
	}).GetProxy, nil
}

type roundRobinSwitcher struct {
	proxyURLs []*url.URL
	index     uint32
}

func (r *roundRobinSwitcher) GetProxy(pr *http.Request) (*url.URL, error) {
	// 取余轮询,每一次调用GetProxy将会使得其index加1(闭包技巧)
	index := atomic.AddUint32(&r.index, 1) - 1
	u := r.proxyURLs[index%uint32(len(r.proxyURLs))]
	return u, nil
}
