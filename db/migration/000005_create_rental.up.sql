CREATE TABLE "rentals" (
     "id" bigserial PRIMARY KEY,
     "field_id" INT,
     "team_id" INT,
     "user_id" INT,
     "comment" varchar,
     "start_date" timestamptz NOT NULL DEFAULT (now()),
     "end_date" timestamptz NOT NULL DEFAULT (now()),
     "duration" INT,
     "status" smallint NOT NULL DEFAULT 1,
     "created_at" timestamptz NOT NULL DEFAULT (now()),
     "updated_at" timestamptz NOT NULL DEFAULT (now()),
     "deleted_at" timestamptz,
     FOREIGN KEY (field_id) REFERENCES fields(id),
     FOREIGN KEY (team_id) REFERENCES teams(id),
     FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_rentals_field_team_user ON rentals (field_id, team_id, user_id);

COMMENT ON COLUMN "rentals"."field_id" IS 'Идентификатор площадки';
COMMENT ON COLUMN "rentals"."team_id" IS 'Идентификатор команды';
COMMENT ON COLUMN "rentals"."user_id" IS 'Идентификатор пользователя';
COMMENT ON COLUMN "rentals"."comment" IS 'Комментарий';
COMMENT ON COLUMN "rentals"."start_date" IS 'Дата начала аренды';
COMMENT ON COLUMN "rentals"."end_date" IS 'Дата завершения аренды';
COMMENT ON COLUMN "rentals"."duration" IS 'Длительность аренды в секундах';
COMMENT ON COLUMN "rentals"."status" IS 'Статус аренды';
COMMENT ON COLUMN "rentals"."created_at" IS 'Дата создания';
COMMENT ON COLUMN "rentals"."updated_at" IS 'Дата изменения';
COMMENT ON COLUMN "rentals"."deleted_at" IS 'Дата удаления';

