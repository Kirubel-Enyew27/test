-- name: AddCartItem :exec
INSERT INTO cart_items (item_id, item_name, price, quantity) 
VALUES ($1, $2, $3, $4);

-- name: CountUniqueItemsInCart :one
SELECT COUNT(DISTINCT item_id) AS unique_items
FROM cart_items;

-- name: RemoveItem :exec
DELETE FROM cart_items WHERE item_id = $1;

-- name: RemoveAllItem :exec
DELETE FROM cart_items;

-- name: FindItemInCart :one
SELECT EXISTS (
    SELECT 1
    FROM cart_items
    WHERE item_id = $1
);

-- name: UpdateItemQuantity :exec
UPDATE cart_items
SET quantity = $2
WHERE item_id = $1;

-- name: ViewCart :many
SELECT item_id, item_name, price, quantity
FROM cart_items;

-- name: CheckoutCart :exec
DELETE FROM cart_items;

-- name: AddProduct :exec
INSERT INTO products (product_id, product_name, price, stock)
VALUES ($1, $2, $3, $4);

-- name: GetProductByID :one
SELECT product_id, product_name, price, stock
FROM products
WHERE product_id = $1;

-- name: UpdateProductStock :exec
UPDATE products
SET stock = $2
WHERE product_id = $1;

-- name: RemoveAllProduct :exec
DELETE FROM products;

-- name: ApplyFlatRateDiscountToCart :exec
UPDATE cart_items
SET price = GREATEST(price - $1, 0) 
WHERE price > 0;

-- name: ApplyPercentageDiscountToCart :exec
UPDATE cart_items
SET price = GREATEST(price - (price * $1 / 100), 0)
WHERE price > 0;
