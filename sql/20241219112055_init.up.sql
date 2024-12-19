CREATE TABLE "users" (
  "id" serial PRIMARY KEY,
  "email" varchar NOT NULL,
  "password" varchar NOT NULL,
  "name" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE TABLE "match" (
  "id" serial PRIMARY KEY,
  "swiper_id" int NOT NULL,
  "swiped_id" int NOT NULL,
  "status" varchar DEFAULT 'active',
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE TABLE "subscriptions" (
  "id" serial PRIMARY KEY,
  "user_id" int NOT NULL,
  "subscription_type" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz,
  "deleted_at" timestamptz,
  "expires_at" timestamptz
);

CREATE TABLE "configurations" (
  "key" varchar PRIMARY KEY,
  "value" jsonb,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE UNIQUE INDEX ON "users" ("email");

CREATE INDEX ON "match" ("swiper_id");

CREATE INDEX ON "match" ("swiped_id");

CREATE UNIQUE INDEX ON "match" ("swiper_id", "swiped_id");

CREATE INDEX ON "subscriptions" ("user_id", "subscription_type", "expires_at");

COMMENT ON COLUMN "match"."status" IS 'active,  inactive';

COMMENT ON COLUMN "subscriptions"."subscription_type" IS 'premium';
