CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS address (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country TEXT NOT NULL,
    city TEXT NOT NULL,
    street TEXT NOT NULL,
    UNIQUE(country, city, street)
);

CREATE TABLE IF NOT EXISTS client (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    birthday TIMESTAMP,
    gender TEXT CHECK (gender IN ('male', 'female')) NOT NULL,
    registration_date TIMESTAMP DEFAULT now(),
    address_id UUID NULL,
    FOREIGN KEY (address_id) REFERENCES address (id)
);

CREATE TABLE IF NOT EXISTS image (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    data BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS supplier (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    address_id UUID NOT NULL,
    phone_number TEXT NOT NULL,
    FOREIGN KEY (address_id) REFERENCES address (id),
    UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS product (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    price FLOAT NOT NULL,
    available_stock INT NOT NULL,
    last_update_date TIMESTAMP DEFAULT now(),
    supplier_id UUID NOT NULL,
    image_id UUID NOT NULL,
    FOREIGN KEY (supplier_id) REFERENCES supplier (id),
    FOREIGN KEY (image_id) REFERENCES image (id)
);

ALTER TABLE product
ADD CONSTRAINT product_stock_nonnegative CHECK (available_stock >= 0);