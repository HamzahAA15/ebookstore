# Online Book Store Service

Online Book Store Service is a service that enables you to buy books online, creat customer account, send notification. It uses an persistence postgresql data source.

**Note:** View this readme file using preview mode (Ctrl-K + V) or (Command-K + V) for better readability.

## Features

- bcrypt hash & JWT Token to secure user's password & user access.
- Strong Password regulation, at least 8 characters in length contains at least one lowercase letter, one uppercase letter, one digit, and one special character
- Username regulation, between 4 and 16 characters in length contains only alphanumeric characters or underscores
- Email should use uniq and your actual email, so you can receive the notification :)
- Asynchronous email notification feature using Gmail SMTP server (using feature flag).
- Modular project structure with dependency injection on the repository, service & controller layers.
- Using Docker Compose to ease the experience of using this service
- Using inMemory Cache to improve api latency for books & order
- `Ready to use DB, since we run migration and data seeder only for you :)`

## Installation

To enable the email notification feature using Gmail SMTP server, follow these steps:
- Set `CONFIG_EMAIL_SERVICE` to `true`.
- Set `CONFIG_AUTH_PASSWORD` with the app secret of your Gmail account.

Set other `CONFIG` variables as required in `config.go`.

## Usage

Run the following command:

```bash
docker-compose build
docker-compose up

accesss through `http://localhost:8080/`
```

import the JSON collection of request from the attachment of email into API platform such as Postman

Create Account first through `{host}/api/customer/register`
 
## API LIST

### Customer Endpoints

<details>

**Register a new customer**
- **URL:** `/api/customer/register`
- **Method:** `POST`
- **Description:** Registers a new customer account.
- **Request Body:**
  - Requires user details like email, password, etc.
- **Response:**
  - Returns a success message & JWT token upon successful registration.
  - Returns an error message if registration fails.

**Login as a customer**
- **URL:** `/api/customer/login`
- **Method:** `POST`
- **Description:** Allows a customer to log in to their account.
- **Request Body:**
  - Requires user credentials like email and password.
- **Response:**
  - Returns a JWT token upon successful login.
  - Returns an error message if login fails.

</details>

### Book Endpoints
<details>

**Get list of books**
- **URL:** `/api/book`
- **Method:** `GET`
- **Description:** Retrieves a list of books available in the store.
- **Request Body:** N/A
- **Response:**
  - Returns a list of books with details like title, author, price, etc.
  - Returns an error message if the list of books cannot be fetched.

</details>

### Order Endpoints
<details>

**Create a new order**
- **URL:** `/api/order`
- **Method:** `POST`
- **Description:** Allows a customer to place a new order.
- **Authorization:** Requires authentication bearer token.
- **Request Body:**
  - Requires details of the order like book IDs, quantity, etc.
- **Response:**
  - Returns a success message upon successful order creation.
  - Returns an error message if order creation fails.

**Get order history of a customer**
- **URL:** `/api/order/order-history`
- **Method:** `GET`
- **Description:** Retrieves the order history of a customer.
- **Authorization:** Requires authentication bearer token.
- **Request Body:** N/A
- **Response:**
  - Returns a list of past orders made by the customer.
  - Returns an error message if the order history cannot be fetched.

</details>

## DB Schema
<details>
<summary>Click to toggle the database schema</summary>

```sql
CREATE TABLE Customers (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL
);

CREATE TABLE Categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE Books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    category_id INTEGER REFERENCES Categories(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE Orders (
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

CREATE TABLE Items (
    id SERIAL PRIMARY KEY,
    book_id INTEGER REFERENCES Books(id),
    quantity INTEGER NOT NULL,
    order_id INTEGER REFERENCES Orders(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

</details>
```
