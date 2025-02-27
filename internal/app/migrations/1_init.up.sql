CREATE TABLE gateways
(
    id                    SERIAL PRIMARY KEY,
    name                  VARCHAR(255) NOT NULL UNIQUE,
    data_format_supported VARCHAR(50)  NOT NULL,
    priority              INT          NOT NULL,
    created_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE gateways
    ADD CONSTRAINT check_gateway_priority CHECK (priority IN (0, 1, 2));

CREATE TABLE countries
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL UNIQUE,
    code       CHAR(2)      NOT NULL UNIQUE,
    currency   CHAR(3)      NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE gateway_countries
(
    gateway_id INT NOT NULL,
    country_id INT NOT NULL,
    PRIMARY KEY (gateway_id, country_id)
);

CREATE SEQUENCE IF NOT EXISTS transactions_tx_counter_seq;
CREATE TABLE transactions
(
    id         INTEGER PRIMARY KEY,
    amount     DECIMAL(10, 2) NOT NULL,
    type       VARCHAR(50)    NOT NULL,
    status     VARCHAR(50)    NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    gateway_id INT            NULL,
    country_id INT            NULL,
    user_id    INT            NOT NULL
);

ALTER TABLE transactions
    ADD CONSTRAINT check_transaction_type CHECK (type IN ('withdraw', 'deposit'));
ALTER TABLE transactions
    ADD CONSTRAINT check_transaction_status CHECK (status IN ('pending', 'completed', 'failed'));

CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(255) NOT NULL UNIQUE,
    email      VARCHAR(255) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    country_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

