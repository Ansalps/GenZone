-- Modify "cart_items" table
ALTER TABLE "public"."cart_items" ADD COLUMN IF NOT EXISTS "product_variant_id" bigint NULL;
ALTER TABLE "public"."cart_items"
  ADD CONSTRAINT "fk_cart_items_product_variant"
  FOREIGN KEY ("product_variant_id") REFERENCES "public"."product_variants" ("id")
  ON UPDATE NO ACTION ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS "idx_cart_items_product_variant_id" ON "public"."cart_items" ("product_variant_id");
UPDATE "public"."cart_items" AS ci
SET "product_variant_id" = pv."id"
FROM "public"."product_variants" AS pv
WHERE ci."product_variant_id" IS NULL
  AND pv."product_id" = ci."product_id"
  AND pv."size" = ci."size";

-- Modify "order_items" table
ALTER TABLE "public"."order_items" ADD COLUMN IF NOT EXISTS "product_variant_id" bigint NULL;
ALTER TABLE "public"."order_items" ADD COLUMN IF NOT EXISTS "size" character varying(32) NULL;
ALTER TABLE "public"."order_items"
  ADD CONSTRAINT "fk_order_items_product_variant"
  FOREIGN KEY ("product_variant_id") REFERENCES "public"."product_variants" ("id")
  ON UPDATE NO ACTION ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS "idx_order_items_product_variant_id" ON "public"."order_items" ("product_variant_id");
UPDATE "public"."order_items" AS oi
SET "size" = COALESCE(oi."size", ''),
    "product_variant_id" = pv."id"
FROM "public"."product_variants" AS pv
WHERE oi."product_id" = pv."product_id"
  AND oi."product_variant_id" IS NULL
  AND pv."size" = COALESCE(oi."size", '');
