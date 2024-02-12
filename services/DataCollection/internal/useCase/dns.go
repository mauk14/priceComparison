package useCase

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	gt "github.com/bas24/googletranslatefree"
	"github.com/chromedp/chromedp"
	"net/http"
	errorsCFG "priceComp/pkg/errors"
	"priceComp/services/DataCollection/internal/domain"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

func dnsParseProducts(url string) (*domain.Products, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errorsCFG.ErrBadUrl
	}

	var product domain.Products

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var wg = sync.WaitGroup{}
	wg.Add(1)
	var err2 error
	go func() {
		defer wg.Done()
		var imageUrl []map[string]string

		if err = chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible(".tns-slide-active"),
			chromedp.AttributesAll(".tns-slide-active img", &imageUrl, chromedp.ByQueryAll)); err != nil {
			err2 = err
			return
		}

		imageData, err := downloadImage(imageUrl[0]["src"])
		if err != nil {
			err2 = err
			return
		}

		var image domain.Images
		image.Image_data = imageData
		image.Main_image = true
		product.Images = append(product.Images, image)

	}()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	description := strings.TrimSpace(doc.Find(".product-card-description-text").Text())
	description, err = gt.Translate(description, "ru", "en")
	if err != nil {
		return nil, err
	}
	product.Description = description

	doc.Find(".product-characteristics__spec").Each(func(i int, selection *goquery.Selection) {
		title := strings.TrimSpace(selection.Find(".product-characteristics__spec-title").Text())
		value := strings.TrimSpace(selection.Find(".product-characteristics__spec-value").Text())
		title, err = gt.Translate(title, "ru", "en")
		if err != nil {
			return
		}
		value, err = gt.Translate(value, "ru", "en")
		if err != nil {
			return
		}
		if strings.ToLower(title) == "model" {
			arr := strings.Split(value, " ")
			product.Brand = arr[0]
			str := ""
			for i, a := range arr {
				if i == 0 {
					continue
				}
				str += a + " "
			}
			product.ProductName = strings.TrimSpace(str)
		} else if strings.ToLower(title) == "type" {
			product.Category = value
		} else if strings.ToLower(title) == strings.ToLower("Seller/Manufacturer Warranty") ||
			strings.ToLower(title) == strings.ToLower("Manufacturer code") ||
			strings.Contains(strings.ToLower(title), "color") ||
			strings.ToLower(title) == strings.ToLower("Children's design") ||
			strings.ToLower(title) == strings.ToLower("Seller/Manufacturer Warranty") {

		} else {
			var Attribute domain.Attributes
			Attribute.AttributeName = title
			Attribute.AttributeValue = value
			product.Attributes = append(product.Attributes, Attribute)
		}
	})

	if err != nil {
		return nil, err
	}

	wg.Wait()

	if err2 != nil {
		return nil, err2
	}

	return &product, nil

}

func DnsParseManyProducts(url string) ([]*domain.Products, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errorsCFG.ErrBadUrl
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	products := make([]*domain.Products, 0, 5)
	var err2 error
	wg := &sync.WaitGroup{}
	doc.Find(".catalog-product").Each(func(i int, selection *goquery.Selection) {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				fmt.Printf("done %d", i)
			}()
			link, _ := selection.Find("a").Attr("href")
			link = "https://www.dns-shop.kz" + link + "characteristics/"
			product, err := dnsParseProducts(link)
			if err != nil {
				err2 = err
				return
			}
			products = append(products, product)
		}()
	})
	wg.Wait()
	if err2 != nil {
		return nil, err2
	}
	return products, nil
}

func dnsParsePrice(product *domain.Products, shop *domain.Shops) (*domain.Prices, error) {

	var url string
	if len(strings.Split(product.ProductName, " ")) <= 1 {
		str := product.Brand + "+" + strings.ReplaceAll(product.ProductName, " ", "+")
		url = fmt.Sprintf("https://www.dns-shop.kz/search/?q=%s", str)
	} else {
		url = fmt.Sprintf("https://www.dns-shop.kz/search/?q=%s", strings.ReplaceAll(product.ProductName, " ", "+"))
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errorsCFG.ErrBadUrl
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var err2 error
	var priceParse string
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible(".product-buy__price"),
			chromedp.Text(".product-buy__price", &priceParse, chromedp.ByQueryAll)); err != nil {
			err2 = err
		}
	}()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	productHtml := doc.Find(".catalog-product").First()
	link := ""
	isProduct := false
	if !strings.Contains(strings.ToLower(productHtml.Find(".catalog-product__name").Text()), strings.ToLower(product.ProductName)) {
		if !strings.Contains(strings.ToLower(doc.Find(".product-card-top__name").Text()), strings.ToLower(product.ProductName)) {
			fmt.Println(url)
			fmt.Println("Не сходятся: " + product.ProductName)
			return nil, fmt.Errorf("%s can't find", product.ProductName)
		}
		isProduct = true
	}

	if isProduct {
		link = res.Request.URL.String()
		fmt.Println(link)
		fmt.Println(doc.Find(".product-card-top__name").Text())
	} else {
		linkText, _ := productHtml.Find(".catalog-product__name").Attr("href")
		link = shop.Link + linkText
		fmt.Println(productHtml.Text())
	}

	var price domain.Prices
	price.Link = link
	price.Product_id = product.Product_id
	price.Shop = *shop
	priceStr := ""

	wg.Wait()

	if err2 != nil {
		return nil, err2
	}
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
		return nil, err
	}
	price.Price = priceInt
	fmt.Println(priceInt)
	fmt.Println(price)

	return &price, nil
}
