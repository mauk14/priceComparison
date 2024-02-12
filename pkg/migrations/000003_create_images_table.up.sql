CREATE TABLE IF NOT EXISTS images (
    id bigserial PRIMARY KEY,
    image_data BYTEA NOT NULL,
    main_image BOOLEAN NOT NULL,
    product_id bigserial REFERENCES products(product_id)
);