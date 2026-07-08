DROP INDEX IF EXISTS idx_emails_idempotency;

ALTER TABLE emails
DROP COLUMN IF EXISTS idempotency_key;