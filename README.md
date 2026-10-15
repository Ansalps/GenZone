# Genzone Backend API

Genzone is a robust backend API for an e-commerce platform, developed using Go, the Gin framework, and PostgreSQL as the database. This project implements essential features for both admins and users, providing seamless management and interaction capabilities.

## Features

### Admin Features
- **Product Management**: Create, edit, and delete products.
- **User Management**: Block or unblock users.
- **Order Management**: List and change the status of orders.
- **Sales Reports**: View sales reports (daily, monthly, yearly).
- **Best-Selling Products**: View the top products and categories by delivered order-item count, with month or year filters.

### User Features
- **Product Browsing**: View available products.
- **Cart and Wishlist**: Add products to cart or wishlist.
- **Order Management**: Place or cancel orders using:
  - Cash on Delivery (COD)
  - Online Payments (integrated with Razorpay)
  - Wallet
- **Invoice Generation**: Generate invoices for orders.

## Best-Selling Products

The admin endpoint `GET /admin/bestselling` returns the ten highest-ranked products and category rankings, based on delivered order items. It accepts these optional query parameters:

| Parameter | Behavior |
| --- | --- |
| `month` | Month number from `1` to `12`, covering that month in the current calendar year. |
| `year` | Calendar year from `1900` through the current year, covering January 1 through December 31. |
| `day` | A date in `YYYY-MM-DD` format, retained for compatibility; filters to that day. |

For example, `/admin/bestselling?month=3` filters sales to March of the current year, while `/admin/bestselling?year=2025` filters sales to all of 2025. The endpoint does not use a `period` filter.

The response shape is:

```json
{
  "best_selling_product": [],
  "best_selling_category": [],
  "message": "retrieved best selling products"
}
```

## Technical Highlights
- **Authentication/Authorization**: Implemented using JWT (JSON Web Tokens).
- **Database**: PostgreSQL for reliable data storage and management.
- **Deployment**: Deployed on AWS EC2. Created a Dockerfile and pushed the image on DockerHub.
