-- Create "click_counts" table
CREATE TABLE "public"."click_counts" (
 "url_id" uuid NOT NULL,
 "total_clicks" bigint NOT NULL DEFAULT 0,
 PRIMARY KEY ("url_id"),
 CONSTRAINT "click_counts_url_id_fkey" FOREIGN KEY ("url_id") REFERENCES "public"."urls" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
