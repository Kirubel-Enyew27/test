-- name: CreateProduct :one
INSERT INTO products (name, price, stock) 
VALUES ($1, $2, $3)
RETURNING product_id;

-- name: AddCartItem :one
INSERT INTO cart_items (product_id, quantity) 
VALUES ($1, $2)
RETURNING item_id;

-- name: RemoveCartItem :exec
DELETE FROM cart_items 
WHERE product_id = $1;

-- name: UpdateCartItemQuantity :exec
UPDATE cart_items 
SET quantity = $2 
WHERE product_id = $1;

-- name: GetCartItems :many
SELECT ci.item_id, p.name, p.price, ci.quantity
FROM cart_items ci
JOIN products p ON ci.product_id = p.product_id;

-- name: CreateDiscount :one
INSERT INTO cart_discounts (percentage, flat_rate) 
VALUES ($1, $2)
RETURNING discount_id;

-- name: GetCartTotal :one
SELECT SUM(p.price * ci.quantity) AS total
FROM cart_items ci
JOIN products p ON ci.product_id = p.product_id;

-- name: ClearCart :exec
DELETE FROM cart_items;
