-- Restore soft-delete columns for active catalog and auth models
ALTER TABLE "public"."addresses" ADD COLUMN IF NOT EXISTS "deleted_at" TIMESTAMPTZ NULL;
ALTER TABLE "public"."admins" ADD COLUMN IF NOT EXISTS "deleted_at" TIMESTAMPTZ NULL;
ALTER TABLE "public"."users" ADD COLUMN IF NOT EXISTS "deleted_at" TIMESTAMPTZ NULL;
ALTER TABLE "public"."categories" ADD COLUMN IF NOT EXISTS "deleted_at" TIMESTAMPTZ NULL;
ALTER TABLE "public"."products" ADD COLUMN IF NOT EXISTS "deleted_at" TIMESTAMPTZ NULL;
ALTER TABLE "public"."product_variants" ADD COLUMN IF NOT EXISTS "deleted_at" TIMESTAMPTZ NULL;
ALTER TABLE "public"."coupons" ADD COLUMN IF NOT EXISTS "deleted_at" TIMESTAMPTZ NULL;
ALTER TABLE "public"."offers" ADD COLUMN IF NOT EXISTS "deleted_at" TIMESTAMPTZ NULL;
