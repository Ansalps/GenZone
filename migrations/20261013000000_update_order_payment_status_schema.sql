ALTER TABLE "public"."orders" RENAME COLUMN "order_status" TO "delivery_status";

ALTER TABLE "public"."orders"
  DROP CONSTRAINT IF EXISTS "chk_orders_order_status";

ALTER TABLE "public"."orders"
  ADD CONSTRAINT "chk_orders_delivery_status"
  CHECK ((delivery_status)::text = ANY ((ARRAY['pending'::character varying, 'shipped'::character varying, 'delivered'::character varying, 'cancelled'::character varying, 'failed'::character varying, 'returned'::character varying])::text[]));

ALTER TABLE "public"."orders"
  ADD COLUMN "order_status" character varying(20) NOT NULL DEFAULT 'pending';

ALTER TABLE "public"."orders"
  ADD CONSTRAINT "chk_orders_order_status"
  CHECK ((order_status)::text = ANY ((ARRAY['pending'::character varying, 'completed'::character varying, 'cancelled'::character varying])::text[]));

ALTER TABLE "public"."payments"
  DROP CONSTRAINT IF EXISTS "chk_payments_payment_status";

ALTER TABLE "public"."payments"
  ALTER COLUMN "payment_status" TYPE character varying(20);

ALTER TABLE "public"."payments"
  ADD CONSTRAINT "chk_payments_payment_status"
  CHECK ((payment_status)::text = ANY ((ARRAY['pending'::character varying, 'paid'::character varying, 'failed'::character varying, 'refunded'::character varying, 'expired'::character varying])::text[]));

DROP TABLE IF EXISTS "public"."temp_addresses";
