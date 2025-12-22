-- Установка расширения UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA store;

-- ==========================
-- ENUM TYPES
-- ==========================
CREATE TYPE admin.admin_role AS ENUM ('senior admin', 'junior admin', 'ghost');
CREATE TYPE store.user_role AS ENUM ('buyer');
CREATE TYPE store.order_status AS ENUM ('pending', 'packaging', 'on the way', 'delivered');
CREATE TYPE store.payment_status AS ENUM ('pending', 'paid', 'failed', 'canceled');
CREATE TYPE store.payment_method AS ENUM ('bank card', 'crypto', 'other wallet');

-- ==========================
-- DOMAINS
-- ==========================
CREATE DOMAIN store.user_snp AS TEXT
CHECK (VALUE ~ '^[А-ЯЁа-яёA-Za-z-]+$' AND char_length(VALUE) BETWEEN 2 AND 50);

CREATE DOMAIN store.username_type AS TEXT
CHECK (VALUE ~ '^[a-z0-9_]+$' AND char_length(VALUE) BETWEEN 2 AND 50);

CREATE DOMAIN store.phone_type AS TEXT
CHECK (VALUE ~ '^[0-9]+$' AND char_length(VALUE) BETWEEN 5 AND 25);

CREATE DOMAIN store.email_type AS TEXT
CHECK (VALUE ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$' AND char_length(VALUE) <= 255);

CREATE DOMAIN store.money_type AS DECIMAL(10,2)
CHECK (VALUE >= 0);

CREATE DOMAIN store.postcode_type AS TEXT
CHECK (char_length(VALUE) BETWEEN 3 AND 9 AND VALUE ~ '^[0-9]+$');

CREATE DOMAIN store.country_type AS TEXT
CHECK (char_length(VALUE) BETWEEN 2 AND 100);

CREATE DOMAIN store.city_type AS TEXT
CHECK (char_length(VALUE) BETWEEN 2 AND 100);

CREATE DOMAIN store.street_type AS TEXT
CHECK (char_length(VALUE) BETWEEN 1 AND 100);

CREATE DOMAIN store.house_type AS TEXT
CHECK (char_length(VALUE) BETWEEN 1 AND 9);

CREATE DOMAIN store.apartment_type AS TEXT
CHECK (char_length(VALUE) BETWEEN 1 AND 9);

-- ==========================
-- TYPES
-- ==========================
CREATE TYPE store.shipping_address_type AS (
    postcode   store.postcode_type,
    country    store.country_type NOT NULL,
    city       store.city_type NOT NULL,
    street     store.street_type NOT NULL,
    house      store.house_type NOT NULL,
    apartment  store.apartment_type,
    phone      store.phone_type,
    surname    store.user_snp NOT NULL,
    name       store.user_snp NOT NULL,
    patronymic store.user_snp
);

-- ==========================
-- SCHEMA admin
-- ==========================
CREATE SCHEMA IF NOT EXISTS admin;

-- Таблица ролей администраторов
CREATE TABLE admin.admin_roles (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name admin.admin_role NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

-- Таблица администраторов
CREATE TABLE admin.admins (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username store.username_type NOT NULL UNIQUE,
    password_hash VARCHAR(512) NOT NULL,
    email store.email_type NOT NULL UNIQUE,
    role_id INT NOT NULL REFERENCES admin.admin_roles(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

-- ==========================
-- SCHEMA store
-- ==========================
CREATE SCHEMA IF NOT EXISTS store;

-- Таблица ролей пользователей
CREATE TABLE store.user_roles (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name store.user_role NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

-- Таблица пользователей
CREATE TABLE store.users (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    uuid UUID UNIQUE DEFAULT store.uuid_generate_v4(),
    surname store.user_snp NOT NULL,
    name store.user_snp NOT NULL,
    patronymic store.user_snp,
    username store.username_type NOT NULL UNIQUE,
    phone store.phone_type UNIQUE,
    email store.email_type NOT NULL UNIQUE,
    password_hash VARCHAR(512) NOT NULL,
    role_id INT NOT NULL REFERENCES store.user_roles(id) ON DELETE RESTRICT,
    token_version INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

-- Таблица категорий
CREATE TABLE store.categories (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(50) NOT NULL UNIQUE,
    parent_id INT REFERENCES store.categories(id) ON DELETE SET NULL,
    cat_code VARCHAR(3) NOT NULL UNIQUE CHECK (cat_code ~ '^[0-9]{3}$'),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

-- Таблица товаров
CREATE TABLE store.products (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(100) NOT NULL CHECK (name <> ''),
    description VARCHAR(512),
    base_price store.money_type NOT NULL,
    category_id INT NOT NULL REFERENCES store.categories(id) ON DELETE RESTRICT,
    product_code VARCHAR(5) NOT NULL UNIQUE CHECK (product_code ~ '^[0-9]{5}$'),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

-- Таблица изображений товаров
CREATE TABLE store.product_images (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    product_id INT NOT NULL REFERENCES store.products(id) ON DELETE CASCADE,
    image_key VARCHAR(255) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    position INT DEFAULT 0 CHECK (position >= 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Таблица размеров
CREATE TABLE store.sizes (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(25) NOT NULL UNIQUE CHECK (name <> ''),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

-- Таблица цветов
CREATE TABLE store.colors (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(25) NOT NULL UNIQUE CHECK (name <> ''),
    hex_code VARCHAR(7) UNIQUE CHECK (hex_code ~ '^#[0-9A-Fa-f]{6}$'),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

-- Таблица вариантов товаров
CREATE TABLE store.product_variants (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    product_id INT NOT NULL REFERENCES store.products(id) ON DELETE CASCADE,
    size_id INT NOT NULL REFERENCES store.sizes(id) ON DELETE RESTRICT,
    color_id INT NOT NULL REFERENCES store.colors(id) ON DELETE RESTRICT,
    price store.money_type NOT NULL,
    sku VARCHAR(11) NOT NULL CHECK (sku ~ '^[0-9]{11}$'),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    CONSTRAINT unique_variant UNIQUE (product_id, size_id, color_id)
);

-- Таблица корзин
CREATE TABLE store.carts (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INT NOT NULL UNIQUE REFERENCES store.users(id) ON DELETE CASCADE,
    item_count INT DEFAULT 0 CHECK (item_count >= 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

-- Таблица элементов корзины
CREATE TABLE store.cart_items (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    cart_id INT NOT NULL REFERENCES store.carts(id) ON DELETE CASCADE,
    variant_id INT NOT NULL REFERENCES store.product_variants(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    added_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Таблица заказов
CREATE TABLE store.orders (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INT NOT NULL REFERENCES store.users(id) ON DELETE RESTRICT,
    total_amount store.money_type NOT NULL,
    status store.order_status NOT NULL DEFAULT 'pending',
    shipping_address store.shipping_address_type NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

-- Таблица элементов заказа
CREATE TABLE store.order_items (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    order_id INT NOT NULL REFERENCES store.orders(id) ON DELETE CASCADE,
    variant_id INT NOT NULL REFERENCES store.product_variants(id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity > 0),
    price_at_purchase store.money_type NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Таблица платежей
CREATE TABLE store.payments (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    order_id INT NOT NULL UNIQUE REFERENCES store.orders(id) ON DELETE CASCADE,
    amount store.money_type NOT NULL,
    payment_method store.payment_method DEFAULT 'bank card',
    status store.payment_status DEFAULT 'pending',
    transaction_uuid UUID UNIQUE DEFAULT store.uuid_generate_v4(),
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

-- Таблица сессий (refresh tokens)
CREATE TABLE store.user_sessions (
    id UUID PRIMARY KEY DEFAULT store.uuid_generate_v4(),
    user_id INT NOT NULL REFERENCES store.users(id) ON DELETE CASCADE,
    refresh_token_hash VARCHAR(1024) NOT NULL,
    jti VARCHAR(256) NOT NULL UNIQUE,
    ip INET,
    device_info VARCHAR(256),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE
);

-- INITIAL DATA
INSERT INTO admin.admin_roles (name) VALUES ('senior admin'), ('junior admin'), ('ghost');
INSERT INTO store.user_roles (name) VALUES ('buyer');