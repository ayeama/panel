-- +up
-- CREATE TABLE IF NOT EXISTS servers (
--     id VARCHAR(64) NOT NULL PRIMARY KEY,
--     name VARCHAR(256) NOT NULL,
--     status VARCHAR(32) NOT NULL
-- );
CREATE TABLE IF NOT EXISTS servers (
    id VARCHAR(64) NOT NULL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    status VARCHAR(32) NOT NULL,
    container_id VARCHAR(64)
);

-- -- +down
-- DROP TABLE servers;
-- -- +seed
-- INSERT INTO servers (id) VALUES ("84464dd3d057341ac22b6c87a44ee48594c81a1ae93632b1ca15b075ee4231fb");
