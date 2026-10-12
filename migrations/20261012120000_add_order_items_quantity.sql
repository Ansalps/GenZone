-- Add quantity to order_items so each order item can reflect the purchased count.
ALTER TABLE "public"."order_items"
    ADD COLUMN IF NOT EXISTS "quantity" bigint NOT NULL DEFAULT 1;
