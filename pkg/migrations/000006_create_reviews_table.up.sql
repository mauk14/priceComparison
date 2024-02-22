CREATE TABLE IF NOT EXISTS reviews(
    id bigserial primary key ,
    message TEXT NOT NULL,
    rating integer NOT NULL,
    user_id bigserial REFERENCES users(id),
    product_id bigserial REFERENCES products(product_id),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
)