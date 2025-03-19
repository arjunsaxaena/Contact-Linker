CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE contacts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NULL,
    phone_number TEXT UNIQUE NULL,
    linked_id UUID REFERENCES contacts(id) ON DELETE SET NULL,
    link_precedence TEXT CHECK (link_precedence IN ('primary', 'secondary')) NOT NULL DEFAULT 'primary',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

ALTER TABLE contacts DROP CONSTRAINT contacts_phone_number_key;
ALTER TABLE contacts DROP CONSTRAINT contacts_email_key;   