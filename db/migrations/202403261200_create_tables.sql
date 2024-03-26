-- Migration for creating Customers table if not exists
CREATE TABLE IF NOT EXISTS Customers (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL
);

-- Migration for creating Categories table if not exists
CREATE TABLE IF NOT EXISTS Categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- Migration for creating Books table if not exists
CREATE TABLE IF NOT EXISTS Books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    category_id INTEGER REFERENCES Categories(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Migration for creating Orders table if not exists
CREATE TABLE IF NOT EXISTS Orders (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES Customers(id),
    customer_reference VARCHAR(255),
    receiver_name VARCHAR(255),
    address VARCHAR(255),
    city VARCHAR(255),
    district VARCHAR(255),
    postal_code VARCHAR(20),
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    shipper VARCHAR(255) NOT NULL,
    airwaybill_number VARCHAR(255) NOT NULL,
    total_item INTEGER,
    total_price DECIMAL(10, 2),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Migration for creating Items table if not exists
CREATE TABLE IF NOT EXISTS Items (
    id SERIAL PRIMARY KEY,
    book_id INTEGER REFERENCES Books(id),
    quantity INTEGER NOT NULL,
    order_id INTEGER REFERENCES Orders(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
