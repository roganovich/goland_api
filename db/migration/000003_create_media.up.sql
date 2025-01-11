CREATE TABLE "files" (
     "id" bigserial PRIMARY KEY,
     "name" uuid NOT NULL,
     "path" varchar NOT NULL,
     "ext" varchar NOT NULL,
     "size" int NOT NULL,
     "created_at" timestamptz NOT NULL DEFAULT (now())
);

COMMENT ON COLUMN "files"."name" IS 'Имя файла';
COMMENT ON COLUMN "files"."path" IS 'Путь к файлу';
COMMENT ON COLUMN "files"."ext" IS 'Расширение';
COMMENT ON COLUMN "files"."size" IS 'Размер';
COMMENT ON COLUMN "files"."created_at" IS 'Дата создания';