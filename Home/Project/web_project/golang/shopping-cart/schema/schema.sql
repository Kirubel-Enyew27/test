CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    stock INT NOT NULL
);

CREATE TABLE cart_items (
    item_id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(product_id),
    quantity INT NOT NULL
);

CREATE TABLE cart_discounts (
    discount_id SERIAL PRIMARY KEY,
    percentage NUMERIC(5, 2) DEFAULT 0,
    flat_rate NUMERIC(10, 2) DEFAULT 0
);
