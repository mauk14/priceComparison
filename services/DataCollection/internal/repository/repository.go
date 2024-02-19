package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"priceComp/services/DataCollection/internal/domain"
)

type DataCollectionRep interface {
	InsertProduct(context.Context, *domain.Products) error
	GetProduct(context.Context, int64) (*domain.Products, error)
	UpdateProduct(context.Context, *domain.Products) error
	GetImage(context.Context, int64) (*domain.Images, error)
	GetAllProduct(context.Context) ([]*domain.Products, error)
	AddShop(context.Context, *domain.Shops) error
	GetShop(context.Context, string) (*domain.Shops, error)
	AddPrice(context.Context, *domain.Prices) error
}

type dataCollectionRep struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) DataCollectionRep {
	return &dataCollectionRep{db: db}
}

func (d *dataCollectionRep) InsertProduct(ctx context.Context, product *domain.Products) error {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return err
	}
	err = tx.QueryRow(ctx,
		`INSERT INTO products (productName, description, category, brand) VALUES ($1, $2, $3, $4) RETURNING product_id`,
		product.ProductName, product.Description, product.Category, product.Brand,
	).Scan(&product.Product_id)

	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	for i, attr := range product.Attributes {
		err = tx.QueryRow(ctx,
			`INSERT INTO productAttributes (attributeName, attributeValue, product_id) VALUES ($1, $2, $3) RETURNING attributeId, product_id`,
			attr.AttributeName, attr.AttributeValue, product.Product_id,
		).Scan(&product.Attributes[i].AttributeID, &product.Attributes[i].Product_id)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	for i, img := range product.Images {
		err = tx.QueryRow(ctx,
			`INSERT INTO images (image_data, main_image, product_id) VALUES ($1, $2, $3) RETURNING id, product_id`,
			img.Image_data, img.Main_image, product.Product_id,
		).Scan(&product.Images[i].Id, &product.Images[i].Product_id)
		if err != nil {
			tx.Rollback(ctx) // Roll back the transaction on error
			return err
		}
	}

	return tx.Commit(ctx)

}

func (d *dataCollectionRep) GetProduct(ctx context.Context, id int64) (*domain.Products, error) {
	var product domain.Products
	err := d.db.QueryRow(ctx, `SELECT product_id, productName, category, brand, description FROM products WHERE product_id = $1`, id).
		Scan(&product.Product_id, &product.ProductName, &product.Category, &product.Brand, &product.Description)
	if err != nil {
		return nil, err
	}
	rows, err := d.db.Query(ctx, `SELECT attributeID, attributeName, attributeValue, product_id FROM productAttributes WHERE product_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	product.Attributes = []domain.Attributes{}
	for rows.Next() {
		var attr domain.Attributes
		if err := rows.Scan(&attr.AttributeID, &attr.AttributeName, &attr.AttributeValue, &attr.Product_id); err != nil {
			return nil, err
		}
		product.Attributes = append(product.Attributes, attr)
	}

	imgRows, err := d.db.Query(ctx, `SELECT id, image_data, main_image, product_id FROM images WHERE product_id = $1`, id)
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
	priceRows, err := d.db.Query(ctx, `SELECT id,price, shop_id, product_id, link FROM prices WHERE product_id = $1`, id)
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
		err = d.db.QueryRow(ctx, `SELECT shopName, link from shops where id=$1`, price.Shop.Id).Scan(&price.Shop.ShopName, &price.Shop.Link)
		if err != nil {
			fmt.Println("here3")
			return nil, err
		}
		product.Prices = append(product.Prices, price)
	}
	return &product, nil
}

func (d *dataCollectionRep) UpdateProduct(ctx context.Context, product *domain.Products) error {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`UPDATE products SET productName = $1, description = $2, category = $3, brand = $4 WHERE product_id = $5`,
		product.ProductName, product.Description, product.Category, product.Brand, product.Product_id,
	)
	if err != nil {
		tx.Rollback(ctx) // Roll back the transaction on error
		return err
	}

	for _, attr := range product.Attributes {
		// Assuming an upsert operation: Update if exists, else insert
		_, err = tx.Exec(ctx,
			`INSERT INTO productAttributes (attributeID, attributeName, attributeValue, product_id) VALUES ($1, $2, $3, $4)
             ON CONFLICT (attributeID) DO UPDATE SET attributeName = EXCLUDED.attributeName, attributeValue = EXCLUDED.attributeValue`,
			attr.AttributeID, attr.AttributeName, attr.AttributeValue, product.Product_id,
		)
		if err != nil {
			tx.Rollback(ctx) // Roll back the transaction on error
			return err
		}
	}

	for _, img := range product.Images {
		_, err = tx.Exec(ctx,
			`INSERT INTO images (id, image_data, main_image, product_id) VALUES ($1, $2, $3, $4)
             ON CONFLICT (id) DO UPDATE SET image_data = EXCLUDED.image_data, main_image = EXCLUDED.main_image`,
			img.Id, img.Image_data, img.Main_image, product.Product_id,
		)
		if err != nil {
			tx.Rollback(ctx) // Roll back the transaction on error
			return err
		}
	}

	return tx.Commit(ctx)
}

func (d *dataCollectionRep) GetImage(ctx context.Context, id int64) (*domain.Images, error) {
	var image domain.Images
	image.Id = id
	err := d.db.QueryRow(ctx, "SELECT image_data, main_image, product_id FROM images WHERE id = $1", image.Id).Scan(&image.Image_data,
		&image.Main_image, &image.Product_id)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (d *dataCollectionRep) GetAllProduct(ctx context.Context) ([]*domain.Products, error) {
	products := make([]*domain.Products, 0, 5)
	rows, err := d.db.Query(ctx, `Select product_id, productName, category, brand, description FROM products`)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var product domain.Products
		if err := rows.Scan(&product.Product_id, &product.ProductName, &product.Category, &product.Brand, &product.Description); err != nil {
			return nil, err
		}
		products = append(products, &product)

	}
	return products, nil
}

func (d *dataCollectionRep) AddShop(ctx context.Context, shops *domain.Shops) error {
	return d.db.QueryRow(ctx,
		`INSERT INTO shops (shopName, link) VALUES ($1, $2) RETURNING id`,
		shops.ShopName, shops.Link).Scan(&shops.Id)
}

func (d *dataCollectionRep) GetShop(ctx context.Context, shopName string) (*domain.Shops, error) {
	var shop domain.Shops
	err := d.db.QueryRow(ctx, `Select id, shopName, link from shops where shopName = $1`, shopName).Scan(&shop.Id, &shop.ShopName, &shop.Link)
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

func (d *dataCollectionRep) AddPrice(ctx context.Context, price *domain.Prices) error {
	return d.db.QueryRow(ctx,
		`INSERT INTO prices (price, shop_id, product_id, link) VALUES ($1, $2, $3, $4) RETURNING id`,
		price.Price, price.Shop.Id, price.Product_id, price.Link,
	).Scan(&price.Id)
}
