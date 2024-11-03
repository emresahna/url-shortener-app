-- Modify "urls" table
ALTER TABLE "public"."urls" ALTER COLUMN "id" SET DEFAULT public.uuid_generate_v4();
-- Modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "id" SET DEFAULT public.uuid_generate_v4();
