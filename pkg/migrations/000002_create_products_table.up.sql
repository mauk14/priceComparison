CREATE TABLE IF NOT EXISTS products (
    product_id bigserial PRIMARY KEY,
    productName VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    category VARCHAR(50) NOT NULL,
    brand VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS productAttributes (
    attributeID bigserial primary key ,
    attributeName VARCHAR(255) NOT NULL,
    attributeValue TEXT NOT NULL,
    product_id bigserial REFERENCES products(product_id)
)