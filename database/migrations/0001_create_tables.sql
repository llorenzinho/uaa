-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    
    password_hash TEXT NOT NULL,
    
    -- is_active BOOLEAN DEFAULT TRUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);


CREATE TABLE jwk_keys (
    kid VARCHAR(64) PRIMARY KEY,
    private_key_pem TEXT NOT NULL,
    public_key_pem TEXT NOT NULL,
    algorithm VARCHAR(10) NOT NULL DEFAULT 'RS256',
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);


CREATE TABLE authorization_codes (
    code VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL, 
    client_id VARCHAR(100) NOT NULL,
    redirect_uri TEXT NOT NULL,
    scope TEXT,
    
    -- Campi per la sicurezza PKCE (Code Challenge)
    code_challenge VARCHAR(255),
    code_challenge_method VARCHAR(20) DEFAULT 'S256',
    
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_auth_codes_expires_at ON authorization_codes(expires_at);
CREATE INDEX idx_auth_codes_user_id ON authorization_codes(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE jwk_keys;
DROP TABLE authorization_codes;
-- +goose StatementEnd