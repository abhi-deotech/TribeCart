-- Create enum types
CREATE TYPE user_role AS ENUM (
    'USER_ROLE_UNSPECIFIED',
    'USER_ROLE_CUSTOMER',
    'USER_ROLE_SELLER',
    'USER_ROLE_ADMIN',
    'USER_ROLE_SUPER_ADMIN'
);

CREATE TYPE user_status AS ENUM (
    'USER_STATUS_UNSPECIFIED',
    'USER_STATUS_ACTIVE',
    'USER_STATUS_PENDING',
    'USER_STATUS_SUSPENDED',
    'USER_STATUS_DELETED'
);

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    role user_role NOT NULL DEFAULT 'USER_ROLE_CUSTOMER',
    status user_status NOT NULL DEFAULT 'USER_STATUS_ACTIVE',
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for users table
CREATE INDEX idx_users_email ON users(LOWER(email));
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Create password reset tokens table
CREATE TABLE password_reset_tokens (
    token_hash VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    used BOOLEAN NOT NULL DEFAULT FALSE
);

-- Create indexes for password_reset_tokens
CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
CREATE INDEX idx_password_reset_tokens_expires_at ON password_reset_tokens(expires_at);

-- Create email verification tokens table
CREATE TABLE email_verification_tokens (
    token_hash VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    used BOOLEAN NOT NULL DEFAULT FALSE
);

-- Create indexes for email_verification_tokens
CREATE INDEX idx_email_verification_tokens_user_id ON email_verification_tokens(user_id);
CREATE INDEX idx_email_verification_tokens_expires_at ON email_verification_tokens(expires_at);

-- Create user_sessions table
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash VARCHAR(255) NOT NULL,
    user_agent TEXT,
    ip_address INET,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for user_sessions
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_refresh_token_hash ON user_sessions(refresh_token_hash);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);

-- Create user_addresses table
CREATE TABLE user_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    label VARCHAR(100) NOT NULL,
    line1 VARCHAR(255) NOT NULL,
    line2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_user_address_label UNIQUE (user_id, label)
);

-- Create indexes for user_addresses
CREATE INDEX idx_user_addresses_user_id ON user_addresses(user_id);
CREATE INDEX idx_user_addresses_is_default ON user_addresses(is_default);

-- Create audit_logs table
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for audit_logs
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- Create function to update updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_addresses_updated_at
BEFORE UPDATE ON user_addresses
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create function to log user changes
CREATE OR REPLACE FUNCTION log_user_changes()
RETURNS TRIGGER AS $$
DECLARE
    v_old_data JSONB;
    v_new_data JSONB;
    v_user_id UUID;
BEGIN
    -- Determine the user ID for the audit log
    v_user_id := NULL;
    
    -- If this is a web request, get the user ID from the JWT token
    -- This requires the application to set the user ID in a transaction variable
    -- For example: `SET LOCAL app.current_user_id = 'user-id'`
    IF current_setting('app.current_user_id', TRUE) ~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$' THEN
        v_user_id := current_setting('app.current_user_id')::UUID;
    END IF;
    
    -- For INSERT operations
    IF TG_OP = 'INSERT' THEN
        v_new_data := to_jsonb(NEW);
        
        -- Don't log sensitive information
        v_new_data := v_new_data - 'password';
        
        INSERT INTO audit_logs (
            user_id, action, entity_type, entity_id, new_values
        ) VALUES (
            v_user_id, 'CREATE', TG_TABLE_NAME, NEW.id, v_new_data
        );
        
        RETURN NEW;
    
    -- For UPDATE operations
    ELSIF TG_OP = 'UPDATE' THEN
        v_old_data := to_jsonb(OLD);
        v_new_data := to_jsonb(NEW);
        
        -- Don't log sensitive information
        v_old_data := v_old_data - 'password';
        v_new_data := v_new_data - 'password';
        
        -- Only log if there are changes
        IF v_old_data != v_new_data THEN
            INSERT INTO audit_logs (
                user_id, action, entity_type, entity_id, old_values, new_values
            ) VALUES (
                v_user_id, 'UPDATE', TG_TABLE_NAME, NEW.id, 
                (SELECT jsonb_object_agg(key, value) FROM jsonb_each(v_old_data) 
                 WHERE (v_old_data->key) IS DISTINCT FROM (v_new_data->key)),
                (SELECT jsonb_object_agg(key, value) FROM jsonb_each(v_new_data) 
                 WHERE (v_old_data->key) IS DISTINCT FROM (v_new_data->key))
            );
        END IF;
        
        RETURN NEW;
    
    -- For DELETE operations
    ELSIF TG_OP = 'DELETE' THEN
        v_old_data := to_jsonb(OLD);
        
        -- Don't log sensitive information
        v_old_data := v_old_data - 'password';
        
        INSERT INTO audit_logs (
            user_id, action, entity_type, entity_id, old_values
        ) VALUES (
            v_user_id, 'DELETE', TG_TABLE_NAME, OLD.id, v_old_data
        );
        
        RETURN OLD;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for audit logging
CREATE TRIGGER users_audit_trigger
AFTER INSERT OR UPDATE OR DELETE ON users
FOR EACH ROW EXECUTE FUNCTION log_user_changes();

CREATE TRIGGER user_addresses_audit_trigger
AFTER INSERT OR UPDATE OR DELETE ON user_addresses
FOR EACH ROW EXECUTE FUNCTION log_user_changes();

-- Create function to ensure only one default address per user
CREATE OR REPLACE FUNCTION ensure_single_default_address()
RETURNS TRIGGER AS $$
BEGIN
    -- Only run the trigger if is_default is being set to true
    IF NEW.is_default = TRUE THEN
        -- Set all other addresses for this user to not be default
        UPDATE user_addresses
        SET is_default = FALSE
        WHERE user_id = NEW.user_id
        AND id != NEW.id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for single default address
CREATE TRIGGER ensure_single_default_address_trigger
BEFORE INSERT OR UPDATE OF is_default ON user_addresses
FOR EACH ROW
WHEN (NEW.is_default = TRUE)
EXECUTE FUNCTION ensure_single_default_address();
