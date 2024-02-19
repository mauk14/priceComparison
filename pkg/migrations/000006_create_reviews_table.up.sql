CREATE TABLE IF NOT EXISTS productAttributes (
    attributeID bigserial primary key ,
    message TEXT NOT NULL,
    rating integer NOT NULL,
    user_id bigserial REFERENCES users(id),
    product_id bigserial REFERENCES products(product_id)
)