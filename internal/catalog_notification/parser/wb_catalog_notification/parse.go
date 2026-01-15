package wbcatalognotification

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

type Product struct {
	URL   string
	Img   string
	Price int64
	Name  string
}

type ParseParams struct {
	NotificationID int64
	URL            string
	Proxy          string
	UserAgent      string
	MaxPrice       int64
}

func (s *Service) Parse(params ParseParams) ([]Product, error) {
	browser, err := s.browserStorage.GetOrCreate(1, params.Proxy)
	if err != nil {
		return nil, err
	}

	page, err := stealth.Page(browser)
	if err != nil {
		return nil, err
	}
	defer page.Close()

	page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: params.UserAgent,
	})
	err = page.Navigate(params.URL)
	if err != nil {
		s.browserStorage.Remove(1)
		return nil, fmt.Errorf("ошибка навигации")
	}

	_, err = page.Timeout(60 * time.Second).Element(".catalog-page__content .product-card")
	if err != nil {
		return nil, fmt.Errorf("карточка товара не найдена")
	}

	cardEls, err := page.Elements(".catalog-page__content .product-card")
	if err != nil {
		return nil, fmt.Errorf("карточки товаров не найдены")
	}

	var products []Product

	for _, cardEl := range cardEls {
		linkEl, err := cardEl.Element("a")
		if err != nil {
			return nil, fmt.Errorf("ссылка не найдена")
		}

		href, err := linkEl.Attribute("href")
		if err != nil {
			return nil, fmt.Errorf("аттрибут href не найден")
		}

		priceEl, err := cardEl.Element(".price__lower-price")
		if err != nil {
			return nil, fmt.Errorf("цена не найдена")
		}

		priceText, err := priceEl.Text()
		if err != nil {
			return nil, fmt.Errorf("не удалось получить текст цены")
		}

		priceText = strings.ReplaceAll(priceText, " ", "")
		priceText = strings.ReplaceAll(priceText, "\xc2\xa0", "")
		priceText = strings.ReplaceAll(priceText, "\u00a0", "")
		priceText = strings.ReplaceAll(priceText, "₽", "")
		priceFloat, err := strconv.ParseFloat(priceText, 64)
		if err != nil {
			return nil, err
		}
		price := int64(priceFloat)

		if price > params.MaxPrice {
			continue
		}

		imgEl, err := cardEl.Element("img")
		if err != nil {
			return nil, fmt.Errorf("картинка не найдена")
		}

		img, err := imgEl.Attribute("src")
		if err != nil {
			return nil, fmt.Errorf("аттрибут src не найден")
		}

		nameEl, err := cardEl.Element(".product-card__name")
		if err != nil {
			return nil, fmt.Errorf("имя товара не найдено")
		}

		nameText, err := nameEl.Text()
		if err != nil {
			return nil, fmt.Errorf("невозможно получить текст имени товара")
		}

		products = append(products, Product{
			URL:   *href,
			Img:   *img,
			Price: price,
			Name:  nameText,
		})
	}

	return products, nil
}
