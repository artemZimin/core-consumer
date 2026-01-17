package wbcatalognotification

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ParseStockParams struct {
	URL       string
	Proxy     string
	UserAgent string
	Cookie    string
}

type StockProduct struct {
	Price          int64
	Quantity       int64
	URL, Name, Img string
}

func (s *Service) ParseStock(params ParseStockParams) ([]StockProduct, error) {
	parts := strings.Split(params.Proxy, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid proxy format, expected ip:port@username:password")
	}

	authParts := strings.Split(parts[1], ":")
	if len(authParts) != 2 {
		return nil, fmt.Errorf("invalid username:password format")
	}

	proxy := parts[0]
	user := authParts[0]
	password := authParts[1]

	parsedURL, err := url.Parse(fmt.Sprintf("http://%s:%s@%s", user, password, proxy))
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить url")
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(parsedURL),
		DialContext: (&net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 2 * time.Second,
		}).DialContext,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   2 * time.Second,
	}
	req, err := http.NewRequest("GET", params.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("при создании запроса")
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "nl-NL,nl;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="143", "Chromium";v="143", "Not A(Brand";v="24"`)
	req.Header.Set("user-agent", params.UserAgent)
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("x-spa-version", "13.19.5")
	req.Header.Set("x-userid", "0")
	req.Header.Set("cookie", params.Cookie)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("при чтении тела ответа")
	}

	var catalogData catalog
	err = json.Unmarshal([]byte(bodyText), &catalogData)
	if err != nil {
		return nil, fmt.Errorf("при маршалинге ответа")
	}

	products := make([]StockProduct, 0)

	for _, p := range catalogData.Products {
		if p.Quantity == 0 || len(p.Sizes) == 0 {
			continue
		}

		strID := strconv.Itoa(int(p.ID))

		var product StockProduct

		product.Name = p.Name
		product.Price = p.Sizes[0].Price.Product / 100
		product.URL = fmt.Sprintf("https://www.wildberries.ru/catalog/%d/detail.aspx", p.ID)
		product.Quantity = p.Quantity
		product.Img = fmt.Sprintf(
			"https://spb-basket-cdn-02bl.geobasket.ru/vol%s/part%s/%s/images/c516x688/1.webp",
			strID[:4],
			strID[:6],
			strID,
		)

		products = append(products, product)
	}

	return products, nil
}
