# Day 13 – Authentication, Password Hashing and JWT Notes

## Authentication vs Authorization

- **Authentication** answers "who are you?" It happens at login, when the user proves their identity with a username and password.
- **Authorization** answers "what are you allowed to do?" It happens on every protected request, by checking the user's role.
- Failing authentication returns **401 Unauthorized** (missing or invalid credentials or token).
- Failing authorization returns **403 Forbidden** (valid user, wrong role).

## Why passwords are never stored in plaintext

- If the database leaks, every user's password is exposed, and people reuse passwords on other sites.
- We store only a **hash**, a one-way transformation that cannot be reversed.
- **bcrypt** is used instead of plain SHA-256 because it is deliberately slow (adjustable cost) and adds a random salt automatically. Slowness makes brute-force attacks expensive. The salt means the same password gives a different hash each time, which defeats precomputed tables.
- To verify a login, `bcrypt.CompareHashAndPassword` re-hashes the attempt using the salt stored inside the hash string.
- bcrypt only uses the first 72 bytes of a password, so longer passwords are rejected.
- The hash is never returned in any API response.

## JWT structure

A JWT is three base64url parts separated by dots: `header.payload.signature`.

1. **Header** – the algorithm (HS256) and the token type.
2. **Payload (claims)** – data about the user (user ID, username, role) plus registered claims such as `exp` (expiry) and `iat` (issued at).
3. **Signature** – HMAC-SHA256 of `header.payload` using the secret key.

Important points:

- The payload is **encoded, not encrypted**. Anyone can read it, so no passwords or sensitive data go inside.
- The signature proves the token was not changed. If someone edits the payload (for example role: admin), the signature no longer matches and the token is rejected (demonstrated in the scratch program).
- Only the server knows the secret, so only the server can create valid tokens.

## Where the secret lives

- In an **environment variable** (`JWT_SECRET`), never in the code and never committed to Git (`.env` is in `.gitignore`).
- If the secret leaks, anyone can forge valid tokens, including admin ones. The fix is to rotate the secret, which invalidates all existing tokens.

## Token expiry

- Every token has an `exp` claim. After it passes, validation fails with an expired error and the API returns 401.
- The client then logs in again to get a new token.
- Access tokens are short-lived (for example 15 minutes) so a stolen token is only useful for a short time.

## Safe parsing checklist

- Verify the signing method is the one we expect (HMAC / HS256) inside the key function.
- Use `jwt.WithValidMethods` to restrict allowed algorithms.
- Always check the error from parsing before trusting any claims.

## What I demonstrated in the scratch program

- bcrypt: correct password accepted, wrong password rejected, same password gives different hashes.
- JWT creation and parsing with custom claims.
- Rejected tokens: expired, tampered payload, and wrong secret.
