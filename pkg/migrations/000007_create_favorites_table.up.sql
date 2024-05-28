CREATE TABLE IF NOT EXISTS favorites(
    id bigserial primary key ,
    user_id bigserial REFERENCES users(id),
    product_id bigserial REFERENCES products(product_id),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
)