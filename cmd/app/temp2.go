package main

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

func main2() {
	ch := make(chan int)

	browser := rod.New().ControlURL(
		launcher.New().Headless(false).MustLaunch(),
	).MustConnect()
	page := stealth.MustPage(browser)
	page.MustNavigate("https://www.easemate.ai/flux-2-ai-image-generator")

	<-ch
}
