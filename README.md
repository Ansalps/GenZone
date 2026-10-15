# GeZOne

GeZOne is a full-stack e-commerce platform with a Go/Gin backend API and a Next.js admin storefront frontend. The project includes product and order management, user authentication, payment flows, reporting, and admin dashboards.

## Overview

This repository contains:

- Go backend API in the root project
- Next.js frontend in `front-end/`
- PostgreSQL support with SQLite fallback for local development
- AWS S3 asset handling
- Razorpay payment integration
- JWT-based authentication and Google OAuth support
- Admin tools for sales reporting, products, users, coupons, and offers

## Tech Stack

- Backend: Go, Gin, GORM, PostgreSQL
- Frontend: Next.js 16, React 19, TypeScript, Tailwind CSS
- Auth: JWT + Google OAuth
- Payments: Razorpay
- Storage: AWS S3
- Mail: SMTP email OTP delivery

## Project Structure

```bash
.
├── admin/             # Admin handlers and endpoints
├── database/         # DB initialization and seeding
├── front-end/        # Next.js admin/front-end app
├── helper/           # Email, S3, validation utilities
├── middleware/       # JWT middleware helpers
├── migrations/       # SQL migrations
├── models/           # Database models
├── public/           # Public endpoints and product/user APIs
├── requestmodels/    # Request validation models
├── responsemodels/   # API response payloads
├── route/            # Route registration
├── templates/        # HTML templates
├── user/             # User flows: carts, orders, checkout, wallet, auth
├── utils/            # Shared helper logic
├── main.go           # Gin server bootstrap
├── go.mod            # Go dependencies
├── atlas.hcl         # Atlas schema configuration
├── README.md         # Project documentation
└── .env.example      # Optional example environment file
```

## Features

### Admin Features
- Product catalog management
- Product offer and stock updates
- User management, blocking, and listing
- Sales reports by date, month, and year
- Best-selling products and category analytics
- Coupon, order, and invoice administration

### User Features
- User signup and login
- OTP-based email verification
- Google authentication
- Product browsing and category search
- Cart and wishlist support
- Checkout flow with COD, wallet, and Razorpay
- Order cancellation and retry logic
- Invoice generation

## Prerequisites

- Go 1.25+
- Node.js 20 LTS or newer
- pnpm (recommended for the frontend)
- PostgreSQL (optional for local dev; SQLite works as a fallback when `DATABASE_URL` is not set)

## Environment Variables

Create a `.env` file in the project root with values similar to:

```env
DATABASE_URL=postgres://username:password@localhost:5432/genzone
JWT_SECRET=your_jwt_secret
EMAIL_USER=your_email@example.com
EMAIL_PASSWORD=your_email_app_password
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
RAZORPAY_KEY_ID=your_razorpay_key_id
RAZORPAY_KEY_SECRET=your_razorpay_key_secret
AWS_REGION=ap-south-1
S3_BUCKET_NAME=your-s3-bucket-name
```

Notes:
- If `DATABASE_URL` is missing, the app falls back to a local SQLite database file named `genzone.db`.
- The backend loads `.env` automatically via `godotenv`.
- AWS, email, Google, and payment values are required for the features that use those integrations.

## Running the Backend

From the project root:

```bash
go mod download
go run .
```

The API server runs on:

- http://localhost:8080

The backend also registers HTML templates from `templates/` and sets CORS for `http://localhost:3000`.

## Running the Frontend

From the `front-end` directory:

```bash
cd front-end
pnpm install
pnpm dev
```

The frontend runs on:

- http://localhost:3000

## Admin and API Notes

The project contains multiple admin endpoints such as:

- product and category management
- user listing and moderation
- sales reports and filters
- best-selling analytics
- order status updates

A representative endpoint is:

```http
GET /admin/bestselling?month=3
```

This returns a best-selling product/category summary for the specified sales period.

## Database Setup

The app is built around GORM models and includes SQL migrations in `migrations/`.

For local development:

- Set `DATABASE_URL` to a PostgreSQL database, or
- Let the app create and use a local SQLite database automatically

## Deployment

The project is structured for deployment as a backend service plus a frontend application. The backend uses AWS services for media storage and the frontend is configured as a standard Next.js app.

## License

This project is currently distributed without a formal license declaration in the repository. If needed, add a license file before open-sourcing or sharing it externally.

## Contributing

Contributions are welcome. For local changes, test the backend and frontend separately before merging, especially around auth, checkout, and reports.

