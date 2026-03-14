DELETE FROM policy_permissions WHERE policy_id = 'policy_global_admin';
DELETE FROM policies WHERE id = 'policy_global_admin';
DROP TABLE IF EXISTS policy_permissions;
DROP TABLE IF EXISTS policies;
