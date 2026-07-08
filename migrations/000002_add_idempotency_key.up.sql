ALTER TABLE emails
ADD COLUMN idempotency_key TEXT;

CREATE UNIQUE INDEX idx_emails_idempotency
ON emails(idempotency_key);