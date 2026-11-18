-- table to store the user information 
CREATE TABLE users (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 username TEXT NOT NULL UNIQUE,
 password_hash TEXT NOT NULL,
 role TEXT NOT NULL CHECK (role IN ('admin','user')),
 created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- table to store the tasks and link each task to specific user
CREATE TABLE tasks (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 title TEXT NOT NULL,
 description TEXT NOT NULL DEFAULT '',
 due_date TIMESTAMPTZ NOT NULL,
 status TEXT NOT NULL CHECK (status IN ('pending','in_progress','done')),
 reminder_status TEXT NOT NULL DEFAULT 'pending'
 CHECK (reminder_status IN ('pending','queued','sent','failed')),
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 3 users (fake placeholder hashes, not real bcrypt output)
INSERT INTO users (username, password_hash, role) VALUES
  ('alice', '$2a$10$placeholderhashalice000000000000000000000000000000000', 'admin'),
  ('bob',   '$2a$10$placeholderhashbob00000000000000000000000000000000000', 'user'),
  ('carol', '$2a$10$placeholderhashcarol0000000000000000000000000000000000', 'user');

-- 10 tasks: alice 4, bob 6, carol 0
INSERT INTO tasks (user_id, title, description, due_date, status, reminder_status) VALUES
  ((SELECT id FROM users WHERE username = 'alice'), 'Review pull requests',  'Check open PRs on the API repo', now() + interval '2 days',  'pending',     'pending'),
  ((SELECT id FROM users WHERE username = 'alice'), 'Write sprint report',   'Summarise this sprint',          now() + interval '5 days',  'in_progress', 'pending'),
  ((SELECT id FROM users WHERE username = 'alice'), 'Update documentation',  'Refresh the API docs',           now() - interval '1 day',   'done',        'sent'),
  ((SELECT id FROM users WHERE username = 'alice'), 'Plan team meeting',     'Prepare the agenda',             now() + interval '10 minutes', 'pending',  'pending'),

  ((SELECT id FROM users WHERE username = 'bob'), 'Fix login bug',           'Wrong error on bad password',    now() + interval '1 day',   'in_progress', 'pending'),
  ((SELECT id FROM users WHERE username = 'bob'), 'Write unit tests',        'Cover the usecase layer',        now() + interval '3 days',  'pending',     'pending'),
  ((SELECT id FROM users WHERE username = 'bob'), 'Deploy to staging',       'Run the deploy script',          now() - interval '2 days',  'pending',     'failed'),
  ((SELECT id FROM users WHERE username = 'bob'), 'Refactor repository',     'Split Mongo and Postgres code',  now() + interval '7 days',  'pending',     'pending'),
  ((SELECT id FROM users WHERE username = 'bob'), 'Set up CI',               'GitHub Actions workflow',        now() - interval '3 days',  'done',        'sent'),
  ((SELECT id FROM users WHERE username = 'bob'), 'Backup database',         'Weekly backup',                  now() + interval '30 minutes', 'pending',  'queued');


SELECT u.username, COUNT(t.id) AS task_count
FROM users u
LEFT JOIN tasks t ON t.user_id = u.id
GROUP BY u.id, u.username
ORDER BY task_count DESC;

SELECT * FROM users;

-- testing the DB 

-- this query gives this error
-- ERROR:  duplicate key value violates unique constraint "users_username_key"
-- Key (username)=(alice) already exists. 
-- SQL state: 23505
-- Detail: Key (username)=(alice) already exists.

INSERT INTO users (username, password_hash, role) VALUES
  ('alice', '$2a$10$placeholderhashalice000000000000000000000000000000000', 'admin');

-- this query gives this error
-- ERROR:  new row for relation "tasks" violates check constraint "tasks_status_check"
-- Failing row contains (35de0bfa-5904-4926-96be-8e7194d32497, 7eba2936-159f-4ec9-ae99-50b17bfb879e, Review pull requests, Check open PRs on the API repo, 2026-11-20 21:34:15.206872+03, wakanda, pending, 2026-11-18 21:34:15.206872+03, 2026-11-18 21:34:15.206872+03). 
-- SQL state: 23514
-- Detail: Failing row contains (35de0bfa-5904-4926-96be-8e7194d32497, 7eba2936-159f-4ec9-ae99-50b17bfb879e, Review pull requests, Check open PRs on the API repo, 2026-11-20 21:34:15.206872+03, wakanda, pending, 2026-11-18 21:34:15.206872+03, 2026-11-18 21:34:15.206872+03).

INSERT INTO tasks (user_id, title, description, due_date, status, reminder_status) VALUES
  ((SELECT id FROM users WHERE username = 'alice'), 'Review pull requests',  'Check open PRs on the API repo', now() + interval '2 days',  'wakanda',     'pending');