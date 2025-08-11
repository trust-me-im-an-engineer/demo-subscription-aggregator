CREATE TABLE IF NOT EXISTS subscriptions
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL,
    price        INTEGER      NOT NULL CHECK (price > 0), -- monthly price
    user_id      UUID         NOT NULL,
    start_date   DATE         NOT NULL,
    end_date     DATE,
    UNIQUE (service_name, user_id)
);
