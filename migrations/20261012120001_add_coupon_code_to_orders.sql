-- Modify "orders" table
ALTER TABLE "public"."orders" ADD COLUMN "coupon_code" character varying(100) NULL DEFAULT '';

-- Backfill the new column for historical orders that were created with a coupon.
UPDATE "public"."orders" AS o
SET "coupon_code" = c."code"
FROM "public"."coupons" AS c
WHERE o."coupon_id" = c."id"
  AND (o."coupon_code" IS NULL OR o."coupon_code" = '');
