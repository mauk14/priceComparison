package useCase

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"priceComp/services/DataCollection/internal/domain"
	"strconv"
	"strings"
	"unicode"
)

func kaspiParsePrice(con context.Context, product *domain.Products, shop *domain.Shops) (*domain.Prices, error) {
	url := ""
	if len(strings.Split(product.ProductName, " ")) <= 1 {
		str := product.Brand + "%20" + strings.ReplaceAll(product.ProductName, " ", "%20")
		url = fmt.Sprintf("https://kaspi.kz/shop/search/?text=%s", str)
	} else {
		url = fmt.Sprintf("https://kaspi.kz/shop/search/?text=%s", strings.ReplaceAll(product.ProductName, " ", "%20"))
	}
	fmt.Println(url)
	ctx, cancel := chromedp.NewContext(con)
	var name, priceParse string
	var productUrl []map[string]string
	defer cancel()
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(".item-card"),
		chromedp.Text(".item-card__name", &name, chromedp.ByQueryAll),
		chromedp.Text(".item-card__prices-price", &priceParse, chromedp.ByQueryAll),
		chromedp.AttributesAll(".item-card__name-link", &productUrl, chromedp.ByQueryAll)); err != nil {
		return nil, err
	}

	if !strings.Contains(strings.ToLower(name), strings.ToLower(product.ProductName)) {
		return nil, fmt.Errorf("%s can't find", product.ProductName)
	}
	priceStr := ""
	link := productUrl[0]["href"]
	for _, c := range priceParse {
		if c == '₸' {
			break
		}
		if !unicode.IsDigit(c) {
			continue
		}
		priceStr += string(c)
	}
	priceInt, err := strconv.Atoi(priceStr)
	if err != nil {
		fmt.Println(url)
		fmt.Println(name, priceStr)
		return nil, err
	}
	fmt.Println(name, priceInt, link)
	var price domain.Prices
	price.Price = priceInt
	price.Shop = *shop
	price.Link = link
	price.Product_id = product.Product_id

	return &price, nil
}
