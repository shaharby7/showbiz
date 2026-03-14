CREATE TABLE connections (
    id          VARCHAR(64)  PRIMARY KEY,
    name        VARCHAR(128) NOT NULL,
    project_id  VARCHAR(64)  NOT NULL,
    provider    VARCHAR(64)  NOT NULL,
    credentials BLOB         NOT NULL,
    config      JSON,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_connections_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT uq_connections_project_name UNIQUE (project_id, name)
);

CREATE INDEX idx_connections_project ON connections(project_id);
CREATE INDEX idx_connections_provider ON connections(project_id, provider);
