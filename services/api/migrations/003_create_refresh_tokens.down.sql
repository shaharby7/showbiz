DROP INDEX idx_refresh_tokens_expires ON refresh_tokens;
DROP INDEX idx_refresh_tokens_user ON refresh_tokens;
DROP TABLE IF EXISTS refresh_tokens;
