CREATE TABLE resources (
    id             VARCHAR(512) PRIMARY KEY,
    name           VARCHAR(128) NOT NULL,
    project_id     VARCHAR(64)  NOT NULL,
    connection_id  VARCHAR(64)  NOT NULL,
    resource_type  VARCHAR(64)  NOT NULL,
    values_json    JSON         NOT NULL,
    status         VARCHAR(32)  NOT NULL DEFAULT 'creating',
    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_resources_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_resources_connection FOREIGN KEY (connection_id) REFERENCES connections(id),
    CONSTRAINT uq_resources_project_name UNIQUE (project_id, name)
);

CREATE INDEX idx_resources_project ON resources(project_id);
CREATE INDEX idx_resources_connection ON resources(connection_id);
CREATE INDEX idx_resources_type ON resources(project_id, resource_type);
CREATE INDEX idx_resources_status ON resources(project_id, status);
