-- name: ListDelayReportsForVendor :many
-- Get all delay reports for a given vendor ID
SELECT delay_reports.*
FROM delay_reports
JOIN orders ON delay_reports.order_id = orders.id
WHERE orders.vendor_id = $1 AND delay_reports.created_at > $2;

-- name: CreateDelayReport :exec
-- Report a new delay for an order
INSERT INTO delay_reports (order_id) VALUES ($1);

-- name: CreateDelayCheck :exec
-- Create a new delay check
INSERT INTO delay_checks (agent_id, report_id) VALUES ($1, $2);

-- name: GetTripByOrderID :one
-- Get a trip record for a specific order_id
SELECT *
FROM trips
WHERE order_id = $1;

-- name: GetOrderByID :one
-- Get an order by ID
SELECT *
FROM orders
WHERE id = $1;

-- name: UpdateOrderDeliveryTime :exec
-- Updates the delivery time for a specific order by its ID
UPDATE orders
SET delivery_time = $2
WHERE id = $1;