CREATE TABLE policy_attachments (
    id          VARCHAR(64) PRIMARY KEY,
    project_id  VARCHAR(64) NOT NULL,
    user_id     VARCHAR(256) NOT NULL,
    policy_id   VARCHAR(64) NOT NULL,
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_pa_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_pa_user FOREIGN KEY (user_id) REFERENCES users(email),
    CONSTRAINT fk_pa_policy FOREIGN KEY (policy_id) REFERENCES policies(id),
    CONSTRAINT uq_pa_project_user_policy UNIQUE (project_id, user_id, policy_id)
);

CREATE INDEX idx_pa_project ON policy_attachments(project_id);
CREATE INDEX idx_pa_user ON policy_attachments(user_id);
CREATE INDEX idx_pa_policy ON policy_attachments(policy_id);
