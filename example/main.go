package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/GiaoGiaoCat/zhima"
	"github.com/chromedp/chromedp"
)

func main() {
	var buf []byte
	// 全国，所有城市，线路不限，不去重，端口4位，稳定时长 5-25 分钟，协议 HTTP
	options := zhima.Options{Pro: 0, City: 0, YYS: 0, MR: 3, PB: 4, Time: 1, Port: 1}
	proxy, err := zhima.GetProxy(options)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	speed, status, err := zhima.TestProxy(proxy)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Printf("speed %d ms, status %d\n", speed, status)

	proxyAddr := fmt.Sprintf("http://%s:%d", proxy.IP, proxy.Port)
	// println(proxyAddr)
	// return

	// user proxy
	o := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer(proxyAddr),
	)

	cx, cancel := chromedp.NewExecAllocator(context.Background(), o...)
	defer cancel()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		cx,
		// context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	err = chromedp.Run(ctx,
		chromedp.Navigate(`https://www.123cha.com/`),
		chromedp.Sleep(10*time.Second),
		chromedp.CaptureScreenshot(&buf),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("fullScreenshot.png", buf, 0600); err != nil {
		log.Fatal(err)
	}
}
