CREATE TABLE "carts" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "product_id" bigserial NOT NULL,
  "quantity" int NOT NULL,
  "price" float NOT NULL,
  "size" varchar NOT NULL
);

CREATE INDEX ON "carts" ("username");

ALTER TABLE "carts" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
ALTER TABLE "carts" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");