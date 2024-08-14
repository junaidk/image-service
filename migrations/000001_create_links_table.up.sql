CREATE TABLE images_metadata (
    hash          TEXT PRIMARY KEY,
    type          TEXT NOT NULL,
    size_width    SMALLINT NOT NULL,
    size_height   SMALLINT NOT NULL,
    camera_model  TEXT,
    lens_model    TEXT,
    location_lat  DOUBLE PRECISION,
    location_long DOUBLE PRECISION
);

CREATE TABLE images (
    id            UUID PRIMARY KEY,
    name          TEXT NOT NULL,
    hash          TEXT NOT NULL,
    created_at     timestamp default current_timestamp,
CONSTRAINT fk_hash
FOREIGN KEY(hash)
REFERENCES images_metadata(hash)
ON DELETE NO ACTION
);

