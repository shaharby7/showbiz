CREATE TABLE users (
    email           VARCHAR(256) PRIMARY KEY,
    password_hash   VARCHAR(256) NOT NULL,
    organization_id VARCHAR(64)  NOT NULL,
    display_name    VARCHAR(256) NOT NULL,
    email_verified  BOOLEAN      NOT NULL DEFAULT FALSE,
    active          BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_users_org FOREIGN KEY (organization_id) REFERENCES organizations(id)
);

CREATE INDEX idx_users_org ON users(organization_id);
