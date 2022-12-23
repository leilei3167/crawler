package collect

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/leilei3167/crawler/proxy"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Fetcher interface {
	Get(url *Request) ([]byte, error)
}

// BaseFetch 实现最最基本的一次网页数据爬取
type BaseFetch struct{}

func (BaseFetch) Get(req *Request) ([]byte, error) {
	resp, err := http.Get(req.URL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code:%d\n", resp.StatusCode)
	}
	bodyReader := bufio.NewReader(resp.Body)
	e := DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return io.ReadAll(utf8Reader)
}

func DetermineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}

// BrowserFetch 实现伪装User-Agent的一次爬取
type BrowserFetch struct {
	Timeout time.Duration
	Proxy   proxy.ProxyFunc
}

func (b BrowserFetch) Get(request *Request) ([]byte, error) {
	client := &http.Client{
		Timeout: b.Timeout,
	}

	if b.Proxy != nil {
		transport := http.DefaultTransport.(*http.Transport)
		transport.Proxy = b.Proxy
		client.Transport = transport
	}

	req, err := http.NewRequest("GET", request.URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	if len(request.Cookie) > 0 {
		req.Header.Set("Cookie", request.Cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return io.ReadAll(utf8Reader)
}

// 操作浏览器来实现对数据的获取,如某些数据被折叠需要鼠标点击才能展开,一般是通过操作Chrome DevTools 协议来实现,Go中最知名的库是chromedp
