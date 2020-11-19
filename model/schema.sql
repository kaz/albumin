CREATE TABLE IF NOT EXISTS photo (
	path TEXT PRIMARY KEY,
	hash TEXT,
	phash BIGINT UNSIGNED,
	timestamp DATETIME
);

CREATE INDEX photo_hash ON photo (hash);
CREATE INDEX photo_phash ON photo (phash);
