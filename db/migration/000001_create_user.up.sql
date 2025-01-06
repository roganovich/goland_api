CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "phone" varchar,
  "status" smallint NOT NULL DEFAULT 1,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE INDEX ON "users" ("email", "email");
CREATE UNIQUE INDEX users_email_unique ON users (email);

COMMENT ON COLUMN "users"."name" IS 'ФИО';
COMMENT ON COLUMN "users"."email" IS 'Email';
COMMENT ON COLUMN "users"."phone" IS 'Телефон';
COMMENT ON COLUMN "users"."status" IS 'Статус';
COMMENT ON COLUMN "users"."created_at" IS 'Дата создания';
COMMENT ON COLUMN "users"."updated_at" IS 'Дата изменения';
COMMENT ON COLUMN "users"."deleted_at" IS 'Дата удаления';