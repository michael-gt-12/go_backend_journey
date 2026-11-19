# Performance notes: list query and the (user_id, due_date) index

## Setup

- PostgreSQL (local), table `tasks` seeded with 10,000 rows via `generate_series`
  (20 users, about 500 tasks each; random due dates and statuses).
- `ANALYZE` run after seeding so the planner has fresh statistics.
- Query (the one behind `GET /tasks`, first page):

  SELECT id, title, due_date, status
  FROM tasks
  WHERE user_id = $1
  ORDER BY due_date, id
  LIMIT 10;

## Before the index (migrations 000001-000002 only)

- Plan: Seq Scan on tasks, followed by Sort using top-N heapsort
- Execution time: 0.948 ms (best of 3 runs)

## After the index (migration 000003: `idx_tasks_user_due` on (user_id, due_date))

- Plan: Index Scan using `idx_tasks_user_due`, followed by Incremental Sort
- Execution time: 0.079 ms (best of 3 runs)

## Conclusion

- Without the index, Postgres reads every one of the 10,000 rows, filters by user,
  and sorts. With the index it jumps straight to that user's rows, already ordered
  by due_date.
- Speed-up: about 12.0x on 10,000 rows. The gap grows as the table grows, because a
  sequential scan is O(N) while an index lookup is O(log N).
- Cost: every INSERT/UPDATE of user_id or due_date must also update the index, and
  the index takes disk space. It is worth it here because listing a user's tasks is
  the most frequent read.
- The partial index `idx_tasks_reminder_pending` keeps the reminder-claim query
  fast without indexing rows that are already sent or failed.
