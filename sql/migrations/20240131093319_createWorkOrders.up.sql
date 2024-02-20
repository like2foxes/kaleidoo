CREATE TABLE IF NOT EXISTS work_orders (
	id SERIAL PRIMARY KEY ,
	customer_id int NOT NULL,
	status TEXT NOT NULL,
	type TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	FOREIGN KEY (customer_id) REFERENCES customers(id)
);
