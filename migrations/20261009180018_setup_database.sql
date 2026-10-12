-- Create "admins" table
CREATE TABLE "public"."admins" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "email" text NULL,
  "password" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_admins_deleted_at" to table: "admins"
CREATE INDEX "idx_admins_deleted_at" ON "public"."admins" ("deleted_at");
-- Create "offers" table
CREATE TABLE "public"."offers" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "product_id" bigint NOT NULL,
  "discount_percentage" numeric NOT NULL,
  "start_at" timestamptz NULL,
  "end_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create "otps" table
CREATE TABLE "public"."otps" (
  "email" text NULL,
  "otp" text NULL,
  "otp_expiry" timestamptz NULL
);
-- Create "payments" table
CREATE TABLE "public"."payments" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NULL,
  "order_id" bigint NULL,
  "order_item_id" bigint NULL,
  "total_amount" numeric NULL,
  "transaction_id" text NULL,
  "payment_date" text NULL,
  "payment_type" character varying(10) NULL DEFAULT 'COD',
  "payment_status" character varying(10) NULL DEFAULT 'pending',
  "description" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_payments_deleted_at" to table: "payments"
CREATE INDEX "idx_payments_deleted_at" ON "public"."payments" ("deleted_at");
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "first_name" text NULL,
  "last_name" text NULL,
  "email" text NULL,
  "password" text NULL,
  "phone" text NULL,
  "profile_picture" text NULL,
  "status" character varying(10) NULL DEFAULT 'Active',
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_email" UNIQUE ("email")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- Create "wallet_transactions" table
CREATE TABLE "public"."wallet_transactions" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "amount" numeric NOT NULL,
  "transaction_type" character varying(50) NOT NULL,
  "description" character varying(255) NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_wallet_transactions_deleted_at" to table: "wallet_transactions"
CREATE INDEX "idx_wallet_transactions_deleted_at" ON "public"."wallet_transactions" ("deleted_at");
-- Create "wallets" table
CREATE TABLE "public"."wallets" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "balance" numeric(10,2) NULL DEFAULT 0,
  PRIMARY KEY ("id")
);
-- Create index "idx_wallets_deleted_at" to table: "wallets"
CREATE INDEX "idx_wallets_deleted_at" ON "public"."wallets" ("deleted_at");
-- Create "wishlists" table
CREATE TABLE "public"."wishlists" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "product_name" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_wishlists_deleted_at" to table: "wishlists"
CREATE INDEX "idx_wishlists_deleted_at" ON "public"."wishlists" ("deleted_at");
-- Create "addresses" table
CREATE TABLE "public"."addresses" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NULL,
  "country" text NULL,
  "state" text NULL,
  "city" text NULL,
  "street_name" text NULL,
  "pin_code" text NULL,
  "phone" text NULL,
  "default" boolean NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_addresses_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_addresses_deleted_at" to table: "addresses"
CREATE INDEX "idx_addresses_deleted_at" ON "public"."addresses" ("deleted_at");
-- Create "categories" table
CREATE TABLE "public"."categories" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "category_name" text NOT NULL,
  "description" text NOT NULL,
  "image_url" text NOT NULL,
  "parent_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_categories_category_name" UNIQUE ("category_name"),
  CONSTRAINT "fk_categories_children" FOREIGN KEY ("parent_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_categories_parent_id" to table: "categories"
CREATE INDEX "idx_categories_parent_id" ON "public"."categories" ("parent_id");
-- Create "products" table
CREATE TABLE "public"."products" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "category_id" bigint NOT NULL,
  "product_name" text NOT NULL,
  "description" text NOT NULL,
  "image_url" text NOT NULL,
  "price" numeric(10,2) NOT NULL,
  "popular" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_categories_products" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_products_category_id" to table: "products"
CREATE INDEX "idx_products_category_id" ON "public"."products" ("category_id");
-- Create "carts" table
CREATE TABLE "public"."carts" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_carts_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_carts_user_id" to table: "carts"
CREATE UNIQUE INDEX "idx_carts_user_id" ON "public"."carts" ("user_id");
-- Create "cart_items" table
CREATE TABLE "public"."cart_items" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "cart_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "size" character varying(32) NOT NULL,
  "quantity" bigint NOT NULL DEFAULT 1,
  "unit_price" numeric(10,2) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_cart_items_product" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_carts_cart_items" FOREIGN KEY ("cart_id") REFERENCES "public"."carts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_cart_items_cart_id" to table: "cart_items"
CREATE INDEX "idx_cart_items_cart_id" ON "public"."cart_items" ("cart_id");
-- Create index "idx_cart_items_product_id" to table: "cart_items"
CREATE INDEX "idx_cart_items_product_id" ON "public"."cart_items" ("product_id");
-- Create index "idx_cart_product_size" to table: "cart_items"
CREATE UNIQUE INDEX "idx_cart_product_size" ON "public"."cart_items" ("cart_id", "product_id", "size");
-- Create "coupons" table
CREATE TABLE "public"."coupons" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "code" text NOT NULL,
  "discount" numeric(5,2) NOT NULL,
  "min_purchase" numeric(10,2) NULL,
  "start_at" timestamptz NULL,
  "end_at" timestamptz NULL,
  "is_active" boolean NULL,
  PRIMARY KEY ("id")
);
-- Create "orders" table
CREATE TABLE "public"."orders" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "user_id" bigint NULL,
  "address_id" bigint NULL,
  "total_amount" numeric NULL,
  "payment_method" character varying(10) NULL DEFAULT 'COD',
  "order_status" character varying(10) NULL DEFAULT 'pending',
  "coupon_id" bigint NULL,
  "offer_discount" numeric NULL,
  "coupon_discount" numeric NULL,
  "total_discount_amount" numeric(10,2) NULL DEFAULT 0,
  "final_amount" numeric(10,2) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_orders_address" FOREIGN KEY ("address_id") REFERENCES "public"."addresses" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_orders_coupon" FOREIGN KEY ("coupon_id") REFERENCES "public"."coupons" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_orders_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_orders_order_status" CHECK ((order_status)::text = ANY ((ARRAY['pending'::character varying, 'shipped'::character varying, 'delivered'::character varying, 'cancelled'::character varying, 'failed'::character varying])::text[]))
);
-- Create "order_items" table
CREATE TABLE "public"."order_items" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "order_id" bigint NULL,
  "product_id" bigint NULL,
  "price" numeric NULL,
  "order_status" character varying(10) NULL DEFAULT 'pending',
  "payment_method" character varying(10) NULL DEFAULT 'COD',
  "coupon_discount" numeric NULL DEFAULT 0,
  "offer_discount" numeric NULL DEFAULT 0,
  "total_discount" numeric NULL DEFAULT 0,
  "paid_amount" numeric NULL DEFAULT 0,
  "delivered_date" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_order_items_order" FOREIGN KEY ("order_id") REFERENCES "public"."orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_order_items_product" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_order_items_order_status" CHECK ((order_status)::text = ANY ((ARRAY['pending'::character varying, 'shipped'::character varying, 'delivered'::character varying, 'cancelled'::character varying, 'failed'::character varying, 'return'::character varying])::text[]))
);
-- Create index "idx_order_items_deleted_at" to table: "order_items"
CREATE INDEX "idx_order_items_deleted_at" ON "public"."order_items" ("deleted_at");
-- Create "product_variants" table
CREATE TABLE "public"."product_variants" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "product_id" bigint NOT NULL,
  "size" character varying(10) NOT NULL,
  "stock" bigint NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_variants" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_product_size" to table: "product_variants"
CREATE UNIQUE INDEX "idx_product_size" ON "public"."product_variants" ("product_id", "size");
-- Create index "idx_product_variants_product_id" to table: "product_variants"
CREATE INDEX "idx_product_variants_product_id" ON "public"."product_variants" ("product_id");
