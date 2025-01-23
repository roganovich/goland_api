CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL,
    "email" varchar NOT NULL,
    "password" text NOT NULL,
    "phone" varchar,
    "city" varchar,
    "logo" varchar,
    "media" jsonb not null default '[]'::jsonb,
    "status" smallint NOT NULL DEFAULT 1,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "deleted_at" timestamptz
);

CREATE INDEX IF NOT EXISTS idx_users_email_phone ON users (email, phone);
CREATE UNIQUE INDEX users_phone_unique ON users (phone);
CREATE UNIQUE INDEX users_email_unique ON users (email);

COMMENT ON COLUMN "users"."name" IS 'ФИО';
COMMENT ON COLUMN "users"."email" IS 'Email';
COMMENT ON COLUMN "users"."password" IS 'Пароль';
COMMENT ON COLUMN "users"."phone" IS 'Телефон';
COMMENT ON COLUMN "users"."city" IS 'Город';
COMMENT ON COLUMN "users"."logo" IS 'Логотип';
COMMENT ON COLUMN "users"."media" IS 'Медиа';
COMMENT ON COLUMN "users"."status" IS 'Статус';
COMMENT ON COLUMN "users"."created_at" IS 'Дата создания';
COMMENT ON COLUMN "users"."updated_at" IS 'Дата изменения';
COMMENT ON COLUMN "users"."deleted_at" IS 'Дата удаления';