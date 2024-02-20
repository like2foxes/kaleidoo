create table if not exists customers (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	logo VARCHAR(255),
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
