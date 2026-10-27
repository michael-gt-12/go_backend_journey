# MongoDB Notes (Day 9)

Hands-on notes for the Task Manager database. Everything was run in mongosh (inside Compass) against a local MongoDB server.

- Database: `taskmanager`
- Collection: `tasks`
- Document shape: `{ _id, title, description, due_date, status }`

## 1. Setup

```javascript
use taskmanager
db.tasks.insertMany([ /* 10 sample tasks */ ])
db.tasks.countDocuments()
```

Result: `10` documents were inserted. MongoDB created the database and the collection automatically on the first insert, and added an `_id` of type `ObjectId` to every document.

Sample data (title, due date, status):

| #   | Title                  | Due date   | Status      |
| --- | ---------------------- | ---------- | ----------- |
| 1   | Finish Task 3 docs     | 2026-10-28 | in_progress |
| 2   | Set up MongoDB         | 2026-10-27 | done        |
| 3   | Learn ObjectID         | 2026-10-27 | in_progress |
| 4   | Write mongo notes      | 2026-10-28 | pending     |
| 5   | Connect Go to MongoDB  | 2026-10-29 | pending     |
| 6   | Use context timeouts   | 2026-10-29 | pending     |
| 7   | Persist the Task API   | 2026-10-30 | pending     |
| 8   | Handle ObjectID errors | 2026-10-30 | pending     |
| 9   | Post Day 8 update      | 2026-10-26 | done        |
| 10  | Tag release v0.1       | 2026-10-26 | done        |

## 2. The five queries

### Query 1: Filter by status

```javascript
db.tasks.find({ status: "pending" });
```

What it does: returns only the documents whose `status` equals `"pending"`. The object passed to `find` is the filter.

Result: 5 documents:

- Write mongo notes
- Connect Go to MongoDB
- Use context timeouts
- Persist the Task API
- Handle ObjectID errors

API equivalent: `GET /tasks?status=pending` (a filtered list).

### Query 2: Sort by due date

```javascript
db.tasks.find().sort({ due_date: 1 });
```

What it does: returns all documents ordered by `due_date`. `1` means ascending (earliest first) and `-1` means descending.

Result: all 10 documents, earliest due date first:

| Due date   | Documents                                    |
| ---------- | -------------------------------------------- |
| 2026-10-26 | Post Day 8 update, Tag release v0.1          |
| 2026-10-27 | Set up MongoDB, Learn ObjectID               |
| 2026-10-28 | Finish Task 3 docs, Write mongo notes        |
| 2026-10-29 | Connect Go to MongoDB, Use context timeouts  |
| 2026-10-30 | Persist the Task API, Handle ObjectID errors |

Documents that share a due date can appear in either order, because MongoDB does not guarantee the order of ties.

### Query 3: Count by status

```javascript
db.tasks.aggregate([{ $group: { _id: "$status", count: { $sum: 1 } } }]);
```

What it does: groups documents by the value of `status` and counts how many are in each group.

- `_id: "$status"` means "group by the status field". The `$` means "the value of this field".
- `count: { $sum: 1 }` adds 1 for every document in the group.

Result (the order of the groups can vary):

```javascript
{ _id: 'pending', count: 5 }
{ _id: 'in_progress', count: 2 }
{ _id: 'done', count: 3 }
```

SQL equivalent: `SELECT status, COUNT(*) FROM tasks GROUP BY status`.

### Query 4: Update many

```javascript
db.tasks.updateMany({ status: "pending" }, { $set: { status: "in_progress" } });
```

What it does: finds every document matching the filter (first argument) and changes the fields in `$set` (second argument). `$set` changes only the named fields and leaves the rest untouched.

Result:

```javascript
{
  acknowledged: true,
  insertedId: null,
  matchedCount: 5,
  modifiedCount: 5,
  upsertedCount: 0
}
```

All 5 pending tasks became `in_progress`, so the collection now has 7 `in_progress` and 3 `done` documents.

Note: `updateOne` or `updateMany` with a plain document and no `$set` throws an error. To replace a whole document on purpose, use `replaceOne`, which keeps the same `_id`.

### Query 5: Delete one

```javascript
db.tasks.deleteOne({ title: "Tag release v0.1" });
```

What it does: removes the first document that matches the filter.

Result:

```javascript
{ acknowledged: true, deletedCount: 1 }
```

The collection now has 9 documents. `db.tasks.countDocuments()` returns `9`.

## 3. Filter operators learned

| Operator       | Meaning                              | Example                                           |
| -------------- | ------------------------------------ | ------------------------------------------------- |
| `$lt` / `$lte` | less than / less than or equal       | `{ due_date: { $lt: new Date("2026-10-29") } }`   |
| `$gt` / `$gte` | greater than / greater than or equal | `{ due_date: { $gte: new Date("2026-10-28") } }`  |
| `$ne`          | not equal                            | `{ status: { $ne: "done" } }`                     |
| `$in`          | matches any value in a list          | `{ status: { $in: ["pending", "in_progress"] } }` |

An empty filter `{}` matches every document, so `deleteMany({})` would delete the whole collection.

## 4. CRUD mapped to the API

| Operation | MongoDB                   | Task Manager endpoint          |
| --------- | ------------------------- | ------------------------------ |
| Create    | `insertOne`, `insertMany` | `POST /tasks`                  |
| Read      | `find`, `findOne`         | `GET /tasks`, `GET /tasks/:id` |
| Update    | `updateOne` with `$set`   | `PUT /tasks/:id`               |
| Delete    | `deleteOne`               | `DELETE /tasks/:id`            |

## 5. ObjectID vs a normal string ID

Every document has an `_id`. By default MongoDB generates it as an `ObjectId`:

```javascript
ObjectId("671e4a2b9f1c3d5e7a8b9c0d");
```

An `ObjectId` is a 12-byte value shown as 24 hex characters. It contains a timestamp, some random bytes and a counter, which is why this works:

```javascript
db.tasks.findOne()._id.getTimestamp();
```

|                        | ObjectId                         | Normal string ID          |
| ---------------------- | -------------------------------- | ------------------------- |
| Type                   | Special MongoDB type             | Plain text                |
| Generated              | Automatically by MongoDB         | Manually by the developer |
| Contains creation time | Yes                              | No                        |
| Filter syntax          | `find({ _id: ObjectId("...") })` | `find({ _id: "..." })`    |

Key point: the ID in a URL such as `GET /tasks/671e...` arrives as a string. Searching with that string finds nothing, because the stored value is an `ObjectId`. The string must be converted to an `ObjectId` first, and a string that is not a valid ObjectID should return `400 Bad Request`. I will handle this on Day 11 when I apply the driver in the API.

## 6. What I learned

1. A document is one record (like one `Task`) and a collection is a group of documents.
2. `_id` is an auto-generated `ObjectId`, not a string, so lookups need a conversion.
3. Filters choose which documents are affected, and `$set` changes only the named fields.
4. `aggregate` with `$group` is how to count or summarise data per group.

Next up (Day 10): connect to MongoDB from Go with the official driver, using `context.WithTimeout` on every call.
