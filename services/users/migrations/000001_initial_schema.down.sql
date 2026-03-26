-- Drop triggers first to avoid dependency issues
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_user_addresses_updated_at ON user_addresses;
DROP TRIGGER IF EXISTS users_audit_trigger ON users;
DROP TRIGGER IF EXISTS user_addresses_audit_trigger ON user_addresses;
DROP TRIGGER IF EXISTS ensure_single_default_address_trigger ON user_addresses;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS log_user_changes();
DROP FUNCTION IF EXISTS ensure_single_default_address();

-- Drop tables
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS email_verification_tokens;
DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS user_addresses;
DROP TABLE IF EXISTS users;

-- Drop types
DROP TYPE IF EXISTS user_role;
DROP TYPE IF EXISTS user_status;
