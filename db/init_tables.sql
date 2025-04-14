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

CREATE TABLE IF NOT EXISTS image (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country TEXT NOT NULL,
    city TEXT NOT NULL,
    street TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS product (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    price FLOAT NOT NULL,
    available_stock INT NOT NULL,
    last_update_date TIMESTAMP DEFAULT now(),
    supplier_id UUID NOT NULL,
    image_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS supplier (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    address_id UUID NOT NULL,
    phone_number TEXT NOT NULLs
);