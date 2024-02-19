package useCase

import (
	"context"
	"errors"
	"fmt"
	"priceComp/services/DataCollection/internal/domain"
	"priceComp/services/DataCollection/internal/repository"
	"strings"
	"sync"
	"time"
)

type DataCollectionService interface {
	AddNewProduct(context.Context, string, string) (*domain.Products, error)
	AddNewManyProduct(context.Context, string, string) ([]*domain.Products, error)
	GetProduct(context.Context, int64) (*domain.Products, error)
	GetImage(context.Context, int64) (*domain.Images, error)
	AddShop(context.Context, string, string) (*domain.Shops, error)
	AddPrice(context.Context, string) ([]*domain.Prices, error)
}

type dataCollectionService struct {
	rep repository.DataCollectionRep
}

func New(repository repository.DataCollectionRep) DataCollectionService {
	return &dataCollectionService{rep: repository}
}

func (d *dataCollectionService) AddNewProduct(ctx context.Context, url, shop string) (*domain.Products, error) {
	var product *domain.Products
	if shop == "dns" {
		pr, err := dnsParseProducts(url)
		if err != nil {
			return nil, err
		}
		product = pr

	}

	err := d.rep.InsertProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (d *dataCollectionService) GetProduct(ctx context.Context, id int64) (*domain.Products, error) {
	product, err := d.rep.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (d *dataCollectionService) GetImage(ctx context.Context, id int64) (*domain.Images, error) {
	image, err := d.rep.GetImage(ctx, id)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (d *dataCollectionService) AddNewManyProduct(ctx context.Context, url, shop string) ([]*domain.Products, error) {
	var products []*domain.Products
	if shop == "dns" {
		pr, err := DnsParseManyProducts(url)
		if err != nil {
			return nil, err
		}
		products = pr

	}
	for i := range products {
		err := d.rep.InsertProduct(ctx, products[i])
		if err != nil {
			return nil, err
		}
	}
	return products, nil
}

func (d *dataCollectionService) AddPrice(ctx context.Context, shopName string) ([]*domain.Prices, error) {
	products, err := d.rep.GetAllProduct(ctx)
	if err != nil {
		return nil, err
	}
	shop, err := d.rep.GetShop(ctx, shopName)
	if err != nil {
		return nil, err
	}
	wg := &sync.WaitGroup{}
	prices := make([]*domain.Prices, 0, 5)
	var err2 error
	if shopName == "dns" {
		for i := range products {
			go func(wg *sync.WaitGroup, i int) {
				wg.Add(1)
				defer wg.Done()
				fmt.Println(products[i].ProductName)
				price, err := dnsParsePrice(products[i], shop)
				if err != nil && !strings.Contains(err.Error(), "can't find") {
					err2 = err
					return
				}
				prices = append(prices, price)
			}(wg, i)
		}
		wg.Wait()
	} else if shopName == "kaspi" {
		for i := range products {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()
			price, err := kaspiParsePrice(ctx, products[i], shop)
			if err != nil {
				switch {
				case strings.Contains(err.Error(), "can't find"):
					continue
				case errors.Is(err, context.DeadlineExceeded):
					continue
				default:
					return nil, err
				}
			}
			prices = append(prices, price)
		}
	}

	if err2 != nil {
		return nil, err2
	}

	for i := range prices {
		if prices[i] == nil {
			continue
		}
		err = d.rep.AddPrice(ctx, prices[i])
		if err != nil {
			return nil, err
		}
	}
	return prices, nil
}

func (d *dataCollectionService) AddShop(ctx context.Context, shopName string, url string) (*domain.Shops, error) {
	var shop domain.Shops
	shop.ShopName = shopName
	shop.Link = url
	err := d.rep.AddShop(ctx, &shop)
	if err != nil {
		return nil, err
	}
	return &shop, nil
}
