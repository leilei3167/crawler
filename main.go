package main

import (
	"fmt"
	doubangroup "github.com/leilei3167/crawler/parse"
	"go.uber.org/zap"
	"time"

	"github.com/leilei3167/crawler/log"
	"go.uber.org/zap/zapcore"

	"github.com/leilei3167/crawler/collect"
	"github.com/leilei3167/crawler/proxy"
)

func main() {
	//plugin, c := log.NewFilePlugin("./log.txt", zapcore.InfoLevel)
	//defer c.Close()
	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)
	logger.Info("log init end")

	// proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8889"}
	proxyURLs := []string{}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		logger.Error("RoundRobinProxySwitcher failed")
	}

	// 按照翻页规则 进行基本的添加(广度优先)
	var worklist []*collect.Request
	for i := 0; i <= 100; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d&type=new", i)
		worklist = append(worklist, &collect.Request{
			URL:       str,
			ParseFunc: doubangroup.ParseURL,
			Cookie:    "",
		})
	}

	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   p,
	}

	// 广度优先
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			body, err := f.Get(item)
			time.Sleep(1 * time.Second)
			if err != nil {
				logger.Error("read content failed",
					zap.Error(err),
				)
				continue
			}
			// 解析数据,新获得的数据将被添加到该item的队列中
			res := item.ParseFunc(body, item)
			for _, item := range res.Items {
				logger.Info("result",
					zap.String("get url:", item.(string)))
			}
			worklist = append(worklist, res.Requests...)
		}
	}
}
