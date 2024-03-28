DROP TABLE IF EXISTS profile CASCADE;
CREATE TABLE IF NOT EXISTS profile (
                                       id SERIAL NOT NULL PRIMARY KEY,
                                       login TEXT NOT NULL UNIQUE DEFAULT '',
                                       password bytea NOT NULL DEFAULT ''
);


DROP TABLE IF EXISTS announcement CASCADE;
CREATE TABLE IF NOT EXISTS announcement (
                                       id SERIAL NOT NULL PRIMARY KEY,
                                       id_profile SERIAL NOT NULL REFERENCES profile(id),
                                       header TEXT NOT NULL DEFAULT '',
                                       photo_href TEXT NOT NULL DEFAULT '',
                                       info TEXT NOT NULL DEFAULT '',
                                       date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                       cost int NOT NULL DEFAULT 0
);
