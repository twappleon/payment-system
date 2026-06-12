CREATE TABLE merchant (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    merchant_no VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    api_key VARCHAR(128) NOT NULL,
    api_secret VARCHAR(256) NOT NULL,
    status TINYINT NOT NULL DEFAULT 1,
    callback_url VARCHAR(512),
    ip_whitelist TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE payment_order (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    merchant_no VARCHAR(64) NOT NULL,
    merchant_order_no VARCHAR(128) NOT NULL,
    platform_order_no VARCHAR(128) NOT NULL UNIQUE,
    channel_code VARCHAR(64),
    channel_order_no VARCHAR(128),
    amount DECIMAL(20, 4) NOT NULL,
    currency VARCHAR(16) NOT NULL,
    status VARCHAR(32) NOT NULL,
    notify_status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
    subject VARCHAR(256),
    client_ip VARCHAR(64),
    expire_time DATETIME,
    paid_time DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    UNIQUE KEY uk_merchant_order (merchant_no, merchant_order_no),
    KEY idx_status (status),
    KEY idx_channel_order_no (channel_order_no),
    KEY idx_created_at (created_at)
);

CREATE TABLE payout_order (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    merchant_no VARCHAR(64) NOT NULL,
    merchant_order_no VARCHAR(128) NOT NULL,
    platform_order_no VARCHAR(128) NOT NULL UNIQUE,
    channel_code VARCHAR(64),
    channel_order_no VARCHAR(128),
    amount DECIMAL(20, 4) NOT NULL,
    currency VARCHAR(16) NOT NULL,
    account_name VARCHAR(128),
    account_no VARCHAR(256),
    bank_code VARCHAR(64),
    status VARCHAR(32) NOT NULL,
    notify_status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
    paid_time DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    UNIQUE KEY uk_merchant_order (merchant_no, merchant_order_no)
);

CREATE TABLE payment_channel (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    channel_code VARCHAR(64) NOT NULL UNIQUE,
    channel_name VARCHAR(128) NOT NULL,
    channel_type VARCHAR(32) NOT NULL,
    api_url VARCHAR(512) NOT NULL,
    merchant_id VARCHAR(128),
    secret_key TEXT,
    public_key TEXT,
    private_key TEXT,
    status TINYINT NOT NULL DEFAULT 1,
    weight INT DEFAULT 100,
    daily_limit DECIMAL(20, 4),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE callback_log (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    channel_code VARCHAR(64) NOT NULL,
    platform_order_no VARCHAR(128),
    channel_order_no VARCHAR(128),
    raw_body TEXT NOT NULL,
    headers TEXT,
    sign_valid TINYINT NOT NULL,
    process_status VARCHAR(32) NOT NULL,
    error_message TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    KEY idx_platform_order_no (platform_order_no),
    KEY idx_channel_order_no (channel_order_no)
);

CREATE TABLE merchant_notify_log (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    merchant_no VARCHAR(64) NOT NULL,
    platform_order_no VARCHAR(128) NOT NULL,
    notify_url VARCHAR(512) NOT NULL,
    request_body TEXT,
    response_body TEXT,
    http_status INT,
    retry_count INT NOT NULL DEFAULT 0,
    status VARCHAR(32) NOT NULL,
    next_retry_time DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    KEY idx_order (platform_order_no),
    KEY idx_next_retry (next_retry_time)
);

CREATE TABLE audit_log (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    operator_id BIGINT NOT NULL,
    operator_name VARCHAR(128),
    action VARCHAR(128) NOT NULL,
    resource_type VARCHAR(64),
    resource_id VARCHAR(128),
    before_data JSON,
    after_data JSON,
    ip VARCHAR(64),
    user_agent VARCHAR(512),
    created_at DATETIME NOT NULL
);

CREATE TABLE outbox_event (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    event_id VARCHAR(128) NOT NULL UNIQUE,
    event_type VARCHAR(64) NOT NULL,
    aggregate_type VARCHAR(64) NOT NULL,
    aggregate_id VARCHAR(128) NOT NULL,
    payload JSON NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
    retry_count INT NOT NULL DEFAULT 0,
    next_retry_time DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    KEY idx_status_next_retry (status, next_retry_time)
);
