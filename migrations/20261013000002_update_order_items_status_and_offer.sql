ALTER TABLE "public"."orders"
  DROP COLUMN IF EXISTS "offer_id";

ALTER TABLE "public"."order_items"
  DROP CONSTRAINT IF EXISTS "chk_order_items_delivery_status",
  DROP CONSTRAINT IF EXISTS "chk_order_items_order_status";

ALTER TABLE "public"."order_items"
  ADD COLUMN IF NOT EXISTS "offer_id" bigint,
  ADD COLUMN IF NOT EXISTS "delivery_status" character varying(20) NOT NULL DEFAULT 'pending',
  ADD COLUMN IF NOT EXISTS "order_status" character varying(20) NOT NULL DEFAULT 'pending';

UPDATE "public"."order_items"
SET "delivery_status" = CASE
  WHEN "order_status" = 'return' THEN 'returned'
  WHEN "order_status" IN ('pending','shipped','delivered','cancelled','failed') THEN "order_status"
  ELSE 'pending'
END,
    "order_status" = CASE
      WHEN "order_status" IN ('pending','completed','cancelled') THEN "order_status"
      WHEN "order_status" IN ('shipped','delivered') THEN 'completed'
      ELSE 'pending'
    END
WHERE "order_status" IS NOT NULL;

ALTER TABLE "public"."order_items"
  ADD CONSTRAINT "chk_order_items_order_status"
  CHECK ((order_status)::text = ANY ((ARRAY['pending'::character varying, 'completed'::character varying, 'cancelled'::character varying])::text[]));

ALTER TABLE "public"."order_items"
  ADD CONSTRAINT "chk_order_items_delivery_status"
  CHECK ((delivery_status)::text = ANY ((ARRAY['pending'::character varying, 'shipped'::character varying, 'delivered'::character varying, 'cancelled'::character varying, 'failed'::character varying, 'returned'::character varying])::text[]));
