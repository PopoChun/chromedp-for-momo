package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const tragetUrl = "https://www.momoshop.com.tw/goods/GoodsDetail.jsp?i_code=9435324&Area=search&mdiv=403&oid=1_4&cid=index&kw=ps5"

func main() {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	if err := chromedp.Run(ctx, doTasks()); err != nil {
		log.Fatal(err)
	}
}

func doTasks() chromedp.Tasks {

	return chromedp.Tasks{
		// ===== Step 1: Initializing login and generating cookies =====
		chromedp.Navigate("https://www.momoshop.com.tw/"),
		chromedp.Click("#bt_0_150_01 > ul.rightMenu > li:nth-child(2) > span"),
		chromedp.SendKeys(`//input[@name="memId"]`, ""),
		chromedp.SetValue(`//input[@name="passwd"]`, ""),
		chromedp.Sleep(3 * time.Second),
		chromedp.Click("#loginForm > dl.leftArea > dd.loginBtn > input[type=image]"),
		// saveCookies(),
		chromedp.Sleep(3 * time.Second),

		// ===== Step 2: Buy it!!! =====
		// loadCookies(),

		chromedp.Navigate(tragetUrl),
		// chromedp.WaitVisible(`//a[@class="buynow"]`),
		chromedp.WaitVisible("#buy_yes > a > img"),
		chromedp.Click("#buy_yes > a"),
		chromedp.Click("#checkoutBar > tbody > tr > td.checkoutButton.selected"),
		chromedp.Click("#payment03"),
		chromedp.Click("#orderSave"),

		chromedp.Sleep(1 * time.Minute),
	}
}

func saveCookies() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		// if err = chromedp.WaitVisible(`#app`, chromedp.ByID).Do(ctx); err != nil {
		// 	log.Fatal(err.Error())
		// 	return
		// }

		cookies, err := network.GetAllCookies().Do(ctx)
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		cookiesData, err := network.GetAllCookiesReturns{Cookies: cookies}.MarshalJSON()
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		if err = ioutil.WriteFile("cookies.tmp", cookiesData, 0755); err != nil {
			log.Fatal(err.Error())
			return
		}
		return
	}
}

func loadCookies() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		if _, _err := os.Stat("cookies.tmp"); os.IsNotExist(_err) {
			log.Fatal(_err.Error())
			return
		}

		cookiesData, err := ioutil.ReadFile("cookies.tmp")
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		cookiesParams := network.SetCookiesParams{}
		if err = cookiesParams.UnmarshalJSON(cookiesData); err != nil {
			log.Fatal(err.Error())
			return
		}

		return network.SetCookies(cookiesParams.Cookies).Do(ctx)
	}
}
