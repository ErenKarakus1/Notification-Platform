CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL,
    recipient TEXT NOT NULL,
    channel TEXT NOT NULL,
    subject TEXT,
    body TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);