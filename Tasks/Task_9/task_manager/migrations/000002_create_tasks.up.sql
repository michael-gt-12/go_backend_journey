CREATE TABLE tasks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title           TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    due_date        TIMESTAMPTZ NOT NULL,
    status          TEXT NOT NULL CHECK (status IN ('pending', 'in_progress', 'done')),
    reminder_status TEXT NOT NULL DEFAULT 'pending'
        CHECK (reminder_status IN ('pending', 'queued', 'sent', 'failed')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);