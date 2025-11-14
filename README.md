# TaskManager

TaskManager is a full-stack project that combines a Go-based REST backend with a multi-page, Bootstrap-powered frontend for managing users and their tasks. It demonstrates a clean separation between HTTP handlers, storage, and configuration on the server while offering a responsive, theme-aware web UI that talks to the API through a small client.

## Highlights

- **Backend**: Go 1.21+, chi router, PostgreSQL storage, layered packages under `internal/*`.
- **Frontend**: Static HTML/JS/CSS served from `frontend/`, reusable navbar/footer, toast notifications, dark/light mode with animated transitions.
- **Auth & Users**: User creation, login via `/auth/login`, localStorage session persistence, unique usernames enforced in DB.
- **Tasks**: CRUD endpoints, per-user filtering, status/priority management, completed-task view.
- **Dev UX**: Makefile targets, Docker Compose stack (app + db + migrations), health checks, environment-specific config.

## Repository Layout

```
backend/
  cmd/app/               Main HTTP server
  internal/config        Config loader
  internal/handlers      auth, user, task route handlers
  internal/storage       Interfaces + PostgreSQL implementation
  migrations/            SQL migrations (migrate-compatible)
  Makefile, docker-compose.yaml, Dockerfile
frontend/
  assets/                JS/CSS shared by all pages
  *.html                 Landing + auth + task/user/settings views
```

## Backend Architecture

- **Entry point**: `backend/cmd/app/main.go` loads YAML config, configures slog logger, opens the Postgres storage (see `internal/storage/postgre`), and mounts routes on chi.
- **Middleware**: Request ID, real IP, structured logging, panic recovery, URL formatting, and permissive CORS to allow local/static frontend hosts.
- **Handlers** (`internal/handlers/...`):
  - `auth.Login` hashes compare credentials and returns the user on success.
  - `user` handlers cover create (`POST /users/`), list (`GET /users/all`), get by username, and delete.
  - `task` handlers cover create, delete, update, list all, list per user, and fetch by ID.
- **Storage layer**: Interface-driven, with PostgreSQL implementation split into files for save/get/list/update/delete operations.
- **Config**: YAML files under `backend/config/` (e.g., `local.yaml`, `container.yaml`). Settings include HTTP address/timeouts and `database` DSN.
- **Health endpoint**: `GET /healthz` pings the DB with a timeout before returning `ok`.

## Database

- PostgreSQL 16 (see `backend/docker-compose.yaml`).
- Migrations managed via the `migrate` Docker image. Initial migration (`1_init`) creates `users`, `tasks`, `project` tables; `2_add_password_and_unique_username` adds password hashes and enforces unique usernames.
- Default DSN in `config/local.yaml`: `postgres://postgres:ayatuly@localhost:5432/postgres`.

## Frontend Overview

- **Shared utilities** (`frontend/assets/app.js`):
  - Theme toggle with radial or crossfade animations (stored under `tm_theme` in localStorage).
  - Form validation helper (`tmValidate`), toast notifications (`tmToast`), datetime formatting helper (`tmToISO`), and `tmRequireUser` gatekeeper.
- **Shared components** (`assets/components.js`): Injects navbar/footer into each page, shows login/register buttons or the signed-in user badge, and wires logout.
- **API client** (`assets/api.js`): Minimal wrapper around `fetch` with `TM_API.*` methods for users, auth, and tasks. The base URL is configurable via Settings and stored as `tm_api_base`.
- **Pages**:
  - `index.html`: marketing landing page, shows the currently configured API base.
  - `login.html` / `register.html`: forms that call the backend, save `tm_user` to localStorage, and redirect to tasks.
  - `tasks.html`: authed view for creating tasks, listing current user tasks, updating status/priority, and deleting tasks.
  - `completed.html`: filters the user’s tasks for `status === done`, allows reopening or deleting.
  - `settings.html`: change theme, configure API base URL, clear local session.
  - `users.html`: admin-style table that lists and deletes users (left accessible even though the navbar no longer links to it).
- **Styling**: `assets/style.css` defines CSS variables for light/dark themes, transitions, reveal animations, and small UX touches.

## Running Locally

### Prerequisites

- Go 1.21+
- Docker & Docker Compose (if you prefer containers)
- PostgreSQL 16 (local or via Docker)

### Option A: Go + local Postgres

1. Ensure a PostgreSQL instance is running and matches the DSN in `backend/config/local.yaml`. Apply migrations:
   ```bash
   cd backend
   migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" up
   ```
   (Or use any migration tool that supports standard SQL scripts.)
2. Start the backend:
   ```bash
   make run
   ```
   This honors `CONFIG_PATH=config/local.yaml` and serves on `localhost:8020`.
3. Serve the frontend (any static server works):
   ```bash
   cd ../frontend
   python -m http.server 5173
   ```
   Update the API base in `settings.html` (or via Settings page) to match `http://localhost:8020`.

### Option B: Docker Compose

```bash
cd backend
cp .env.example .env   # provide DB_USER, DB_PASS, DB_NAME
make up                # builds app image, runs db + migrate + app
```

- App accessible on `http://localhost:8080`.
- Database exposed on `localhost:5432`.
- Tear down with `make down`.

## Useful Make Targets

- `make build` – cross-compile static binary to `bin/app`.
- `make build-local` – build for your platform.
- `make test` – run Go unit tests.
- `make docker-build` – build the backend image.
- `make up` / `make down` – start/stop the full stack.
- `make migrate` – run the migration container once.
- `make db-shell` – open `psql` inside the DB container.

## API Cheat Sheet

| Method | Path                     | Description                          |
|--------|------------------------- |--------------------------------------|
| POST   | `/auth/login`            | Login with `user_name` & `password`. |
| POST   | `/users/`                | Create user.                         |
| GET    | `/users/all`             | List all users.                      |
| GET    | `/users/{user_name}`     | Fetch single user.                   |
| DELETE | `/users/{user_name}`     | Remove user.                         |
| POST   | `/tasks/`                | Create task for a user.              |
| GET    | `/tasks/{id}`            | Get task by ID.                      |
| PUT    | `/tasks/{id}`            | Update task status/priority/etc.     |
| DELETE | `/tasks/{id}`            | Delete task.                         |
| GET    | `/tasks/all`             | List every task (admin).             |
| GET    | `/tasks/by-user/{id}`    | List tasks for a specific user.      |
| GET    | `/healthz`               | Health probe (DB ping).              |

All responses are JSON and adopt a `{ data: ..., status: ... }` pattern in handlers.

## Testing

- Unit tests exist for handler logic under `backend/internal/handlers/*` (e.g., `task/save_test.go`). Run with `make test`.
- Frontend testing is manual; use browser devtools or add your preferred tooling.

## Frontend Tips

- The API base defaults to `http://localhost:8080`. If you run the backend on another port, open Settings and update the base.
- Theme preference (`tm_theme`) and animation mode (`tm_theme_anim`) persist in localStorage. A radial wipe animation is enabled by default.
- Authed pages call `tmRequireUser()` – if `tm_user` is absent, the user is redirected to login.

## Deployment Notes

- Update `backend/config/container.yaml` (or add `prod.yaml`) with production DSNs and addresses.
- Configure HTTPS termination (e.g., via reverse proxy) in front of the Go app; the server currently serves plain HTTP.
- Ensure CORS is restricted to allowed origins in production by editing the `cors.Handler` options in `main.go`.

## Roadmap Ideas

- Replace localStorage sessions with JWTs or cookie-based auth.
- Add pagination/filtering to `/users/all` and `/tasks/all`.
- Implement project endpoints (scaffolded in migrations but commented out in router).
- Add automated frontend tests or a modern framework build if necessary.

Happy hacking! Feel free to open issues or PRs if you build on top of TaskManager.

