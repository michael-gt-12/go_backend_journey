CREATE INDEX idx_tasks_user_due ON tasks (user_id, due_date);

CREATE INDEX idx_tasks_reminder_pending ON tasks (due_date)
    WHERE reminder_status = 'pending';