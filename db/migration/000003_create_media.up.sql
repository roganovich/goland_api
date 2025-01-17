CREATE TABLE "medias" (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL,
    "path" varchar NOT NULL,
    "ext" varchar NOT NULL,
    "size" int NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

COMMENT ON COLUMN "medias"."name" IS 'Имя файла';
COMMENT ON COLUMN "medias"."path" IS 'Путь к файлу';
COMMENT ON COLUMN "medias"."ext" IS 'Расширение';
COMMENT ON COLUMN "medias"."size" IS 'Размер';
COMMENT ON COLUMN "medias"."created_at" IS 'Дата создания';