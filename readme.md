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

### Book Endpoints

**Get list of books**
- **URL:** `/api/book`
- **Method:** `GET`
- **Description:** Retrieves a list of books available in the store.
- **Request Body:** N/A
- **Response:**
  - Returns a list of books with details like title, author, price, etc.
  - Returns an error message if the list of books cannot be fetched.

### Order Endpoints

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

