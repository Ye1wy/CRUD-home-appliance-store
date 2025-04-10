CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS client (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    birthday TIMESTAMP,
    gender TEXT CHECK (gender IN ('male', 'female')) NOT NULL,
    registation_date TIMESTAMP DEFAULT now(),
    address_id UUID NOT NULL
);