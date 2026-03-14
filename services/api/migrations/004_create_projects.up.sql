CREATE TABLE projects (
    id              VARCHAR(64)  PRIMARY KEY,
    name            VARCHAR(128) NOT NULL,
    organization_id VARCHAR(64)  NOT NULL,
    description     TEXT,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_projects_org FOREIGN KEY (organization_id) REFERENCES organizations(id),
    CONSTRAINT uq_projects_org_name UNIQUE (organization_id, name)
);

CREATE INDEX idx_projects_org ON projects(organization_id);
