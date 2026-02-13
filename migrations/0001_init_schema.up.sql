CREATE TABLE merchants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    primary_contact VARCHAR(255),
    timezone VARCHAR(100) DEFAULT 'Asia/Kolkata',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE stores (
    id SERIAL PRIMARY KEY,
    merchant_id INT NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(512),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    merchant_id INT NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_merchant_id ON users(merchant_id);
CREATE INDEX idx_users_phone ON users(phone);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    merchant_id INT NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    store_id INT NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    external_ref VARCHAR(255),
    amount BIGINT NOT NULL,
    currency VARCHAR(10) DEFAULT 'INR',
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    paid_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_merchant_store_created_at ON orders(merchant_id, store_id, created_at);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    merchant_id INT NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    store_id INT NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    channel VARCHAR(16) NOT NULL,
    amount BIGINT NOT NULL,
    currency VARCHAR(10) DEFAULT 'INR',
    time TIMESTAMPTZ NOT NULL,
    upi_ref VARCHAR(128),
    payer_vpa VARCHAR(255),
    payer_name VARCHAR(255),
    raw_message_id VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_merchant_store_time ON payments(merchant_id, store_id, time);
CREATE INDEX idx_payments_upi_ref ON payments(upi_ref);

CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    payment_id INT NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    confidence DOUBLE PRECISION NOT NULL,
    matched_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_matches_order_id ON matches(order_id);
CREATE INDEX idx_matches_payment_id ON matches(payment_id);

CREATE TABLE exceptions (
    id SERIAL PRIMARY KEY,
    merchant_id INT NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    store_id INT NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    order_id INT REFERENCES orders(id) ON DELETE SET NULL,
    payment_id INT REFERENCES payments(id) ON DELETE SET NULL,
    type VARCHAR(64) NOT NULL,
    reason VARCHAR(512),
    resolved BOOLEAN NOT NULL DEFAULT FALSE,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_exceptions_merchant_store_created_at ON exceptions(merchant_id, store_id, created_at);

