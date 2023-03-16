CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "hash_password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "changed_password_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00z',
  "created_at" timestamptz NOT NULL DEFAULT (Now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

CREATE INDEX ON "accounts" ("owner");

CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
