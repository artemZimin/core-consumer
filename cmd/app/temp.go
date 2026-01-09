package main

import (
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

func main1() {
	browser := rod.New().ControlURL(
		launcher.New().Headless(true).MustLaunch(),
	).MustConnect()
	page := stealth.MustPage(browser)
	page.MustNavigate("https://www.wildberries.ru/catalog/0/search.aspx?page=1&sort=priceup&search=playstation+5+%D1%81+%D0%B4%D0%B8%D1%81%D0%BA%D0%BE%D0%B2%D0%BE%D0%B4%D0%BE%D0%BC&targeturl=ST")

	page.Timeout(60 * time.Second).Element(".product-card")
	page.MustScreenshot("1.png")
}
