-- Modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "created_at" TYPE timestamp, ALTER COLUMN "updated_at" TYPE timestamp, ALTER COLUMN "deleted_at" TYPE timestamp;
-- Modify "urls" table
ALTER TABLE "public"."urls" DROP CONSTRAINT "urls_user_id_fkey", ALTER COLUMN "user_id" DROP NOT NULL, ALTER COLUMN "created_at" TYPE timestamp, ALTER COLUMN "updated_at" TYPE timestamp, ALTER COLUMN "expire_time" TYPE timestamp, ALTER COLUMN "deleted_at" TYPE timestamp, ADD
 CONSTRAINT "urls_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
