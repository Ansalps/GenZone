-- Modify "addresses" table
ALTER TABLE "public"."addresses" DROP COLUMN "deleted_at";
-- Modify "admins" table
ALTER TABLE "public"."admins" DROP COLUMN "deleted_at";
-- Modify "order_items" table
ALTER TABLE "public"."order_items" DROP COLUMN "deleted_at", ALTER COLUMN "unit_price" SET DEFAULT 0, ALTER COLUMN "order_status" TYPE character varying(20), ALTER COLUMN "delivery_status" DROP NOT NULL;
-- Modify "orders" table
ALTER TABLE "public"."orders" ALTER COLUMN "delivery_status" TYPE character varying(20), ALTER COLUMN "order_status" DROP NOT NULL;
-- Modify "payments" table
ALTER TABLE "public"."payments" DROP CONSTRAINT "chk_payments_payment_status", DROP COLUMN "deleted_at";
-- Modify "users" table
ALTER TABLE "public"."users" DROP COLUMN "deleted_at";
-- Modify "wallet_transactions" table
ALTER TABLE "public"."wallet_transactions" DROP COLUMN "deleted_at";
-- Modify "wallets" table
ALTER TABLE "public"."wallets" DROP COLUMN "deleted_at";
-- Modify "wishlists" table
ALTER TABLE "public"."wishlists" DROP COLUMN "deleted_at";
