ALTER TABLE files_metadata
ADD COLUMN checksum BYTEA NOT NULL,
ADD COLUMN encryption_key BYTEA;
