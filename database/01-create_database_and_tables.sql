CREATE DATABASE walks;

\c walks;        

CREATE TABLE walks (
    id character varying(50) PRIMARY KEY,
    walk_date timestamp with time zone NOT NULL UNIQUE,
    duration numeric NOT NULL,
    rate_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE rate (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    start timestamp with time zone NOT NULL,
    end timestamp with time zone,
    amount numeric NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE google_token(
    access_token character varying(255),
    refresh_token character varying(255),
    token_type character varying(50),
    expires_at timestamp with time zone,
);

INSERT INTO rate (start, amount) VALUES (
    423705600, 34,
);