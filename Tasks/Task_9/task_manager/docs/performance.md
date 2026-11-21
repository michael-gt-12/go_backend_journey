# Performance notes: list query, the (user_id, due_date) index, and pagination

## Setup

- PostgreSQL 18.7 (local, Windows), table `tasks` seeded with 10,000 rows via `generate_series`
  (20 users, about 500 tasks each; random due dates and statuses).
- `ANALYZE` run after seeding so the planner has fresh statistics.
- Every number is the best of 3 runs of `EXPLAIN (ANALYZE, BUFFERS)`.
- The data is small and fully cached in memory, so absolute times are tiny. The
  point is the ratio and the shape of the plans, not the milliseconds.

## 1. Effect of the (user_id, due_date) index

Query (the one behind `GET /tasks`, first page):

    SELECT id, title, due_date, status
    FROM tasks
    WHERE user_id = $1
    ORDER BY due_date, id
    LIMIT 10;

### Before the index (migrations 000001-000002 only)

- Plan: Seq Scan on tasks (500 rows kept, 9,512 removed by the filter), followed by
  a Sort using top-N heapsort
- Buffers: 197 shared hits
- Execution time: 1.268 ms (best of 3 runs: 1.894, 2.313, 1.268)

### After the index (migration 000003: `idx_tasks_user_due` on (user_id, due_date))

- Plan: Index Scan using `idx_tasks_user_due`, followed by an Incremental Sort
  (the index already returns rows sorted by due_date, so only ties on due_date
  need sorting by id)
- Buffers: 13 shared hits
- Execution time: 0.132 ms (best of 3 runs: 0.132, 0.140, 0.210)

### Conclusion

- Without the index, Postgres reads every one of the 10,000 rows, filters by user,
  and sorts. With the index it jumps straight to that user's rows, already ordered
  by due_date.
- Speed-up: about 9.6x on 10,000 rows (1.268 ms to 0.132 ms), and 15x fewer buffers
  (197 to 13). The gap grows as the table grows, because a sequential scan is O(N)
  while an index lookup is O(log N).
- Cost: every INSERT/UPDATE of user_id or due_date must also update the index, and
  the index takes disk space. It is worth it here because listing a user's tasks is
  the most frequent read.
- The partial index `idx_tasks_reminder_pending` keeps the reminder-claim query
  fast without indexing rows that are already sent or failed.## 2. OFFSET vs keyset pagination

Both queries fetch the same 10 rows: page 1000 of the whole table (page size 10),
ordered by `(due_date, id)`, with a temporary index on `(due_date, id)`.

    -- OFFSET
    SELECT id, title, due_date FROM tasks ORDER BY due_date, id LIMIT 10 OFFSET 9990;

    -- keyset (the cursor is the last row of the previous page)
    SELECT id, title, due_date FROM tasks
    WHERE (due_date, id) > ($cursor_due_date, $cursor_id)
    ORDER BY due_date, id LIMIT 10;

| Method                           | Plan                            | Rows processed            | Buffers (hits) | Execution time (best of 3) |
| -------------------------------- | ------------------------------- | ------------------------- | -------------- | -------------------------- |
| `LIMIT 10 OFFSET 9990`           | Seq Scan, then Sort of all rows | 10,012 then 10,000 sorted | 197            | 5.321 ms                   |
| keyset `(due_date, id) > cursor` | Index Scan on `(due_date, id)`  | 10                        | 12             | 0.071 ms                   |

### Conclusion

- Keyset was about 75x faster here (0.071 ms vs 5.321 ms; about 120x on the median
  of the 3 runs, 0.083 ms vs 10.1 ms). It also touched 12 buffers instead of 197.
- OFFSET must produce and throw away every skipped row before it returns the 10 it
  needs. For this deep page the planner did not even use the index: it chose a Seq
  Scan plus a full sort, because it expected to read nearly the whole table anyway.
  The cost therefore grows with the page number and the table size. Keyset jumps
  straight to the cursor in the index and stays flat however deep the page is.
- Honest caveat: with 10,000 rows the real-world difference is a few milliseconds.
  It becomes important on tables with millions of rows or for clients that page deep.
- Trade-off: keyset cannot jump to "page N" and needs a stable unique sort key
  (`due_date, id`) plus a cursor from the previous page. It also does not return a
  total count without an extra query.

## Decision

`GET /tasks` uses OFFSET (`?page=&limit=`) because the spec asks for page numbers
and `X-Page`/`X-Total-Count` headers, and one user's task list is small (the
`(user_id, due_date)` index keeps the first pages fast). Keyset pagination is the
upgrade path if lists become very large or clients need infinite scroll.
