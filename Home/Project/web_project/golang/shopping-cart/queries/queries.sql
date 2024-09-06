-- name: AddCartItem :exec
INSERT INTO cart_items (item_id, item_name, price, quantity) 
VALUES ($1, $2, $3, $4);

-- name: CountUniqueItemsInCart :one
SELECT COUNT(DISTINCT item_id) AS unique_items
FROM cart_items;


