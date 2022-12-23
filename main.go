package main

import (
	"time"

	"github.com/leilei3167/crawler/log"
	"go.uber.org/zap/zapcore"

	"github.com/leilei3167/crawler/collect"
	"github.com/leilei3167/crawler/proxy"
)

func main() {
	plugin, c := log.NewFilePlugin("./log.txt", zapcore.InfoLevel)
	defer c.Close()
	logger := log.NewLogger(plugin)
	logger.Info("log init end")

	// proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8889"}
	proxyURLs := []string{}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		logger.Error("RoundRobinProxySwitcher failed")
	}
	url := "https://google.com"
	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   p,
	}

	body, err := f.Get(url)
	if err != nil {
		logger.Sugar().Errorf("read content failed:%v", err)
		return
	}
	logger.Warn(string(body))
}
