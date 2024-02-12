CREATE TABLE IF NOT EXISTS prices (
    id bigserial PRIMARY KEY,
    price int not null,
    shop_id bigserial REFERENCES shops(id),
    product_id bigserial REFERENCES products(product_id),
    link TEXT NOT NULL
);