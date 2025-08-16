CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Create "users" table
CREATE TABLE "public"."users" (
 "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
 "username" character varying(50) NOT NULL,
 "password" character varying(255) NOT NULL,
 "created_at" timestamptz NOT NULL DEFAULT now(),
 PRIMARY KEY ("id"),
 CONSTRAINT "users_username_key" UNIQUE ("username")
);
-- Create "urls" table
CREATE TABLE "public"."urls" (
 "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
 "original_url" text NOT NULL,
 "shortened_code" character varying(10) NOT NULL,
 "user_id" uuid NOT NULL,
 "created_at" timestamptz NOT NULL DEFAULT now(),
 PRIMARY KEY ("id"),
 CONSTRAINT "urls_shortened_code_key" UNIQUE ("shortened_code"),
 CONSTRAINT "urls_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_shortened_code" to table: "urls"
CREATE INDEX "idx_shortened_code" ON "public"."urls" ("shortened_code");
-- Create index "idx_user_id" to table: "urls"
CREATE INDEX "idx_user_id" ON "public"."urls" ("user_id");
