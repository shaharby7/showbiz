ALTER TABLE resources DROP FOREIGN KEY fk_resources_connection;
ALTER TABLE resources MODIFY connection_id VARCHAR(64) NOT NULL;
ALTER TABLE resources ADD CONSTRAINT fk_resources_connection FOREIGN KEY (connection_id) REFERENCES connections(id);
