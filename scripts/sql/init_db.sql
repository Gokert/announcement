DROP TABLE IF EXISTS profile CASCADE;
CREATE TABLE IF NOT EXISTS profile (
                                       id SERIAL NOT NULL PRIMARY KEY,
                                       login TEXT NOT NULL UNIQUE DEFAULT '',
                                       password bytea NOT NULL DEFAULT ''
);
