CREATE TABLE policies (
    id              VARCHAR(64)  PRIMARY KEY,
    name            VARCHAR(128) NOT NULL,
    scope           VARCHAR(32)  NOT NULL,
    organization_id VARCHAR(64),
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_policies_org FOREIGN KEY (organization_id) REFERENCES organizations(id),
    CONSTRAINT chk_policies_scope CHECK (scope IN ('global', 'organization')),
    CONSTRAINT uq_policies_scope_org_name UNIQUE (scope, organization_id, name)
);

CREATE INDEX idx_policies_org ON policies(organization_id);
CREATE INDEX idx_policies_scope ON policies(scope);

CREATE TABLE policy_permissions (
    policy_id  VARCHAR(64) NOT NULL,
    permission VARCHAR(64) NOT NULL,

    PRIMARY KEY (policy_id, permission),
    CONSTRAINT fk_policy_perms_policy FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE
);

-- Seed the global Administrator policy
INSERT INTO policies (id, name, scope, organization_id) VALUES
    ('policy_global_admin', 'Administrator', 'global', NULL);

INSERT INTO policy_permissions (policy_id, permission) VALUES
    ('policy_global_admin', '*:*');
