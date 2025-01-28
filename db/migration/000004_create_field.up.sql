CREATE TABLE "fields" (
     "id" bigserial PRIMARY KEY,
     "name" varchar NOT NULL,
     "description" text,
     "city" varchar NOT NULL,
     "address" text NOT NULL,
     "location" jsonb not null default '[]'::jsonb,
     "square" INT,
     "info" text,
     "places" INT,
     "dressing" bool,
     "toilet" bool,
     "display" bool,
     "parking" bool,
     "for_disabled" bool,
     "logo" varchar,
     "media" jsonb not null default '[]'::jsonb,
     "status" smallint NOT NULL DEFAULT 1,
     "created_at" timestamptz NOT NULL DEFAULT (now()),
     "updated_at" timestamptz NOT NULL DEFAULT (now()),
     "deleted_at" timestamptz
);

CREATE INDEX IF NOT EXISTS idx_fields_name_city ON fields (name, city);
CREATE UNIQUE INDEX unique_fields_name ON fields (name, city);

COMMENT ON COLUMN "fields"."name" IS 'Название';
COMMENT ON COLUMN "fields"."description" IS 'Описание';
COMMENT ON COLUMN "fields"."city" IS 'Город';
COMMENT ON COLUMN "fields"."address" IS 'Адрес';
COMMENT ON COLUMN "fields"."location" IS 'Координаты';
COMMENT ON COLUMN "fields"."square" IS 'Площадь';
COMMENT ON COLUMN "fields"."info" IS 'Административная информация';
COMMENT ON COLUMN "fields"."places" IS 'Кол во места';
COMMENT ON COLUMN "fields"."dressing" IS 'Наличие раздевалки';
COMMENT ON COLUMN "fields"."toilet" IS 'Наличие туалета';
COMMENT ON COLUMN "fields"."display" IS 'Наличие цифрового табло';
COMMENT ON COLUMN "fields"."parking" IS 'Наличие парковки';
COMMENT ON COLUMN "fields"."for_disabled" IS 'Подходит для инвалидов';
COMMENT ON COLUMN "fields"."logo" IS 'Логотип';
COMMENT ON COLUMN "fields"."media" IS 'Медиа';
COMMENT ON COLUMN "fields"."status" IS 'Статус';
COMMENT ON COLUMN "fields"."created_at" IS 'Дата создания';
COMMENT ON COLUMN "fields"."updated_at" IS 'Дата изменения';
COMMENT ON COLUMN "fields"."deleted_at" IS 'Дата удаления';