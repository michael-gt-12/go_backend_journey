-- pick one user and run the list query the API uses
SELECT id AS uid FROM users WHERE username = 'user2' \gset

EXPLAIN (ANALYZE, BUFFERS)
SELECT id, title, due_date, status
FROM tasks
WHERE user_id = :'uid'
ORDER BY due_date, id
LIMIT 10;