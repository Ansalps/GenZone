-- Rename the legacy order_items.price column to unit_price and add total_amount.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'order_items'
          AND column_name = 'price'
    )
       AND NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'order_items'
          AND column_name = 'unit_price'
    ) THEN
        ALTER TABLE "public"."order_items" RENAME COLUMN "price" TO "unit_price";
    END IF;
END $$;

ALTER TABLE "public"."order_items"
    ADD COLUMN IF NOT EXISTS "total_amount" numeric NULL DEFAULT 0;

UPDATE "public"."order_items"
SET "unit_price" = COALESCE("unit_price", 0),
    "total_amount" = COALESCE("total_amount", ("quantity" * COALESCE("unit_price", 0)) - COALESCE("total_discount", 0))
WHERE "unit_price" IS NULL OR "total_amount" IS NULL;

ALTER TABLE "public"."order_items"
    DROP COLUMN IF EXISTS "price";
