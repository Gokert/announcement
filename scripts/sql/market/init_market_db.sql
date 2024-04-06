DROP TABLE IF EXISTS announcement CASCADE;
CREATE TABLE IF NOT EXISTS announcement (
                                            id SERIAL NOT NULL PRIMARY KEY,
                                            id_profile SERIAL NOT NULL,
                                            header TEXT NOT NULL DEFAULT '',
                                            photo_href TEXT NOT NULL DEFAULT '',
                                            info TEXT NOT NULL DEFAULT '',
                                            date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                            cost int NOT NULL DEFAULT 0
);