-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "deleted_at" timestamp NULL,
  "username" character varying(255) NOT NULL,
  "password" character varying(1000) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "users_username_key" UNIQUE ("username")
);
-- Create "posts" table
CREATE TABLE "public"."posts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "content" text NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "posts_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_posts_user_id" to table: "posts"
CREATE INDEX "idx_posts_user_id" ON "public"."posts" ("user_id");
-- Create "comments" table
CREATE TABLE "public"."comments" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "content" text NOT NULL,
  "user_id" uuid NOT NULL,
  "post_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "comments_post_id_fkey" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "comments_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_comments_post_id" to table: "comments"
CREATE INDEX "idx_comments_post_id" ON "public"."comments" ("post_id");
-- Create index "idx_comments_user_id" to table: "comments"
CREATE INDEX "idx_comments_user_id" ON "public"."comments" ("user_id");
