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

-- name: ApplyDiscountToCart :exec
UPDATE cart_items
SET price = price - (price * $1 / 100)
WHERE price > 0;

-- name: ViewCart :many
SELECT item_id, item_name, price, quantity
FROM cart_items;

