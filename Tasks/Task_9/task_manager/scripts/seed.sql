-- 20 users (user1 is the admin)
INSERT INTO users (username, password_hash, role)
SELECT 'user' || g,
       'not-a-real-hash',
       CASE WHEN g = 1 THEN 'admin' ELSE 'user' END
FROM generate_series(1, 20) AS g;

-- 10,000 tasks spread evenly over the 20 users (about 500 each)
WITH u AS (
    SELECT id, row_number() OVER (ORDER BY username) AS rn FROM users
)
INSERT INTO tasks (user_id, title, description, due_date, status)
SELECT u.id,
       'Task ' || g,
       'Seeded task number ' || g,
       now() + ((random() * 180 - 30) * interval '1 day'),
       (ARRAY['pending', 'in_progress', 'done'])[1 + floor(random() * 3)::int]
FROM generate_series(1, 10000) AS g
JOIN u ON u.rn = ((g - 1) % 20) + 1;

-- refresh planner statistics so EXPLAIN reflects the new data
ANALYZE users;
ANALYZE tasks;

SELECT count(*) AS users FROM users;
SELECT count(*) AS tasks FROM tasks;