CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    city VARCHAR(255) NOT NULL,
    uniform_color VARCHAR(50),
    participant_count INT,
    responsible_id INT,
    disability_category VARCHAR(100),
    logo VARCHAR(255),
    media jsonb not null default '[]'::jsonb,
    "status" smallint NOT NULL DEFAULT 1,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "deleted_at" timestamptz,
    FOREIGN KEY (responsible_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_teams_name_city ON teams (name, city);
CREATE UNIQUE INDEX unique_team_name ON teams (name, city);

COMMENT ON COLUMN "teams"."name" IS 'Название';
COMMENT ON COLUMN "teams"."description" IS 'Описание';
COMMENT ON COLUMN "teams"."city" IS 'Город';
COMMENT ON COLUMN "teams"."uniform_color" IS 'Цвет формы';
COMMENT ON COLUMN "teams"."participant_count" IS 'Кол-во участников';
COMMENT ON COLUMN "teams"."responsible_id" IS 'Ответственный';
COMMENT ON COLUMN "teams"."disability_category" IS 'Категория инвалидности';
COMMENT ON COLUMN "teams"."logo" IS 'Логотип';
COMMENT ON COLUMN "teams"."media" IS 'Медиа';
COMMENT ON COLUMN "teams"."status" IS 'Статус';
COMMENT ON COLUMN "teams"."created_at" IS 'Дата создания';
COMMENT ON COLUMN "teams"."updated_at" IS 'Дата последнего обновления';
COMMENT ON COLUMN "teams"."deleted_at" IS 'Дата удаления';

