-- Modify "urls" table
ALTER TABLE "public"."urls" ADD COLUMN "expire_time" timestamptz NULL, ADD COLUMN "is_deleted" boolean NULL DEFAULT false, ADD COLUMN "is_active" boolean NULL DEFAULT true, ADD COLUMN "deleted_at" timestamptz NULL;
-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "is_deleted" boolean NULL DEFAULT false, ADD COLUMN "is_active" boolean NULL DEFAULT true, ADD COLUMN "updated_at" timestamptz NULL, ADD COLUMN "deleted_at" timestamptz NULL;
