CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "phone" varchar UNIQUE NOT NULL,
  "address" varchar NOT NULL,
  "province" bigserial NOT NULL,
  "role" bigserial NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "reset_password_token" varchar NOT NULL DEFAULT 'abc',
  "rspassword_token_expired_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "provinces" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL
);

CREATE TABLE "roles" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL
);

CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "product_name" varchar UNIQUE NOT NULL,
  "thumb" varchar NOT NULL,
  "price" float NOT NULL
);

CREATE TABLE "imgs_product" (
  "id" bigserial PRIMARY KEY,
  "product_id" bigserial NOT NULL,
  "image" varchar NOT NULL
);

CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL
);

CREATE TABLE "products_in_category" (
  "id" bigserial PRIMARY KEY,
  "category_id" bigserial NOT NULL,
  "product_id" bigserial NOT NULL
);

CREATE TABLE "descriptions_product" (
  "id" bigserial PRIMARY KEY,
  "product_id" bigserial NOT NULL,
  "gender" varchar NOT NULL,
  "material" varchar NOT NULL,
  "size" varchar NOT NULL,
  "size_of_model" varchar NOT NULL
);

CREATE TABLE "store" (
  "id" bigserial PRIMARY KEY,
  "product_id" bigserial NOT NULL,
  "size" varchar NOT NULL,
  "quantity" int NOT NULL
);

CREATE TABLE "promotions" (
  "id" bigserial PRIMARY KEY,
  "title" varchar UNIQUE NOT NULL,
  "description" varchar NOT NULL,
  "discount_percent" float NOT NULL,
  "start_date" timestamptz NOT NULL DEFAULT 'now()',
  "end_date" timestamptz NOT NULL DEFAULT '9999-01-01 00:00:00Z'
);

CREATE TABLE "orders" (
  "booking_id" varchar PRIMARY KEY,
  "user_booking" varchar NOT NULL,
  "promotion_id" varchar NOT NULL,
  "status" varchar NOT NULL DEFAULT 'validated',
  "booking_date" Date NOT NULL DEFAULT 'now()',
  "address" varchar NOT NULL,
  "province" bigserial NOT NULL,
  "tax" float NOT NULL,
  "amount" float NOT NULL,
  "payment_method" varchar NOT NULL DEFAULT 'SHIP COD'
);

CREATE TABLE "items_order" (
  "id" bigserial PRIMARY KEY,
  "booking_id" varchar NOT NULL,
  "product_id" bigserial NOT NULL,
  "quantity" int NOT NULL,
  "price" float NOT NULL
);

CREATE TABLE "feedbacks" (
  "id" bigserial PRIMARY KEY,
  "user_comment" varchar NOT NULL,
  "product_commented" bigserial NOT NULL,
  "rating" varchar NOT NULL,
  "commention" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "imgs_product" ("product_id");

CREATE INDEX ON "products_in_category" ("category_id");

CREATE UNIQUE INDEX ON "products_in_category" ("category_id", "product_id");

CREATE INDEX ON "store" ("product_id");

CREATE UNIQUE INDEX ON "store" ("product_id", "size");

CREATE INDEX ON "orders" ("user_booking");

CREATE INDEX ON "orders" ("promotion_id");

CREATE INDEX ON "items_order" ("booking_id");

CREATE INDEX ON "items_order" ("product_id");

CREATE INDEX ON "feedbacks" ("user_comment");

CREATE INDEX ON "feedbacks" ("product_commented");

CREATE INDEX ON "sessions" ("username");

COMMENT ON COLUMN "orders"."tax" IS 'must be positive';

ALTER TABLE "users" ADD FOREIGN KEY ("province") REFERENCES "provinces" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("role") REFERENCES "roles" ("id");

ALTER TABLE "imgs_product" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "products_in_category" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "products_in_category" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "descriptions_product" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "store" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_booking") REFERENCES "users" ("username");

ALTER TABLE "orders" ADD FOREIGN KEY ("promotion_id") REFERENCES "promotions" ("title");

ALTER TABLE "orders" ADD FOREIGN KEY ("province") REFERENCES "provinces" ("id");

ALTER TABLE "items_order" ADD FOREIGN KEY ("booking_id") REFERENCES "orders" ("booking_id");

ALTER TABLE "items_order" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "feedbacks" ADD FOREIGN KEY ("user_comment") REFERENCES "users" ("username");

ALTER TABLE "feedbacks" ADD FOREIGN KEY ("product_commented") REFERENCES "products" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
