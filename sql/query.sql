-- name: CreateUser :one
INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *;

-- name: GetUser :one
SELECT id, email FROM users WHERE email = $1 AND password = $2;

-- name: GetUserPassword :one
SELECT password FROM users WHERE email = $1;

-- name: CreateCustomer :one
INSERT INTO customers (name, logo) VALUES ($1, $2) RETURNING *;

-- name: UpdateCustomer :one
UPDATE customers SET name = $1, logo = $2 WHERE id = $3 RETURNING *;

-- name: GetCustomer :one
SELECT id, name, logo FROM customers WHERE id = $1;

-- name: GetCustomerLogo :one
SELECT logo FROM customers WHERE id = $1;

-- name: ListCustomers :many
SELECT id, name, logo FROM customers;

-- name: DeleteCustomer :exec
DELETE FROM customers WHERE id = $1;

-- name: DeleteCustomerWorkOrders :exec
DELETE FROM work_orders WHERE customer_id = $1;

-- name: DeleteWorkOrder :exec
DELETE FROM work_orders WHERE id = $1;

-- name: CreateWorkOrder :one
INSERT INTO work_orders (customer_id, status, type) VALUES ($1, $2, $3) RETURNING *;

-- name: GetWorkOrder :one
SELECT wo.id, wo.type, wo.status, c.name as customer FROM work_orders wo
	JOIN customers c ON c.id = wo.customer_id
	WHERE wo.id = $1;

-- name: GetWorkOrders :many
SELECT wo.id, c.name as customer, status, type FROM work_orders wo
	JOIN customers c ON c.id = wo.customer_id;

-- name: UpdateWorkOrder :one
UPDATE work_orders SET customer_id = $1, status = $2, type = $3 WHERE id = $4 RETURNING *;

-- name: GetFilteredWorkOrders :many
SELECT wo.id, c.name as customer, status, type FROM work_orders wo
	JOIN customers c ON c.id = wo.customer_id
	WHERE wo.status IN ($1, $2, $3, $4, $5, $6);

-- name: ListForecast :many
SELECT c.id, c.name, COUNT(wo.id) as count FROM customers c
	LEFT JOIN work_orders wo ON c.id = wo.customer_id AND wo.type = 'pm'
	GROUP BY c.id, c.name
	HAVING COUNT(wo.id) > 0;
