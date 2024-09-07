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



