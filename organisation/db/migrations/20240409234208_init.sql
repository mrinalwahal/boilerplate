-- +goose Up
-- create "organisations" table
CREATE TABLE "public"."organisations" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "title" text NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_organisations_title" CHECK (length(title) > 0)
);

-- +goose Down
-- reverse: create "organisations" table
DROP TABLE "public"."organisations";
