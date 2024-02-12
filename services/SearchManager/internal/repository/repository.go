package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"priceComp/services/SearchManager/internal/domain"
)

type SearchRepository interface {
	GetAll(context.Context, string, string, string, string, int, int) ([]*domain.Products, error)
}

type searchRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) SearchRepository {
	return &searchRepository{db: db}
}

func (s *searchRepository) GetAll(ctx context.Context, searchKeyword, category, brand, sortBy string, limit, offset int) ([]*domain.Products, error) {
	var products []*domain.Products
	baseQuery := `SELECT prod.product_id, prod.productName, prod.category, prod.brand, prod.description, p.price FROM products prod join prices p on prod.product_id = p.product_id
WHERE 1=1`
	if searchKeyword != "" {
		baseQuery += fmt.Sprintf(" AND (productName ILIKE '%%%s%%' OR description ILIKE '%%%s%%')", searchKeyword, searchKeyword)
	}
	if category != "" {
		baseQuery += fmt.Sprintf(" AND category = '%s'", category)
	}
	if brand != "" {
		baseQuery += fmt.Sprintf(" AND brand = '%s'", brand)
	}
	switch sortBy {
	case "priceAsc":
		baseQuery += " ORDER BY price ASC"
	case "priceDesc":
		baseQuery += " ORDER BY price DESC"
	case "productName":
		baseQuery += " ORDER BY productname"
	}

	baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := s.db.Query(ctx, baseQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product domain.Products
		var priceQuery int
		err := rows.Scan(&product.Product_id, &product.ProductName, &product.Category, &product.Brand, &product.Description, &priceQuery)
		if err != nil {
			return nil, err
		}
		attrRows, err := s.db.Query(ctx, `SELECT attributeID, attributeName, attributeValue, product_id FROM productAttributes WHERE product_id = $1`, product.Product_id)
		if err != nil {
			return nil, err
		}
		defer attrRows.Close()
		product.Attributes = []domain.Attributes{}
		for attrRows.Next() {
			var attr domain.Attributes
			if err := attrRows.Scan(&attr.AttributeID, &attr.AttributeName, &attr.AttributeValue, &attr.Product_id); err != nil {
				return nil, err
			}
			product.Attributes = append(product.Attributes, attr)
		}

		imgRows, err := s.db.Query(ctx, `SELECT id, image_data, main_image, product_id FROM images WHERE product_id = $1`, product.Product_id)
		if err != nil {
			return nil, err
		}
		defer imgRows.Close()
		product.Images = []domain.Images{}
		for imgRows.Next() {
			var img domain.Images
			if err := imgRows.Scan(&img.Id, &img.Image_data, &img.Main_image, &img.Product_id); err != nil {
				return nil, err
			}
			product.Images = append(product.Images, img)
		}
		priceRows, err := s.db.Query(ctx, `SELECT id,price, shop_id, product_id, link FROM prices WHERE product_id = $1`, product.Product_id)
		if err != nil {
			fmt.Println("here1")
			return nil, err
		}
		defer priceRows.Close()
		product.Prices = []domain.Prices{}
		for priceRows.Next() {
			var price domain.Prices
			if err = priceRows.Scan(&price.Id, &price.Price, &price.Shop.Id, &price.Product_id, &price.Link); err != nil {
				fmt.Println("here2")
				return nil, err
			}
			err = s.db.QueryRow(ctx, `SELECT shopName, link from shops`).Scan(&price.Shop.ShopName, &price.Shop.Link)
			if err != nil {
				fmt.Println("here3")
				return nil, err
			}
			product.Prices = append(product.Prices, price)
		}
		products = append(products, &product)
	}
	return products, nil

}
