ALTER TABLE users
    ADD COLUMN role_id int REFERENCES roles(id) DEFAULT 1;