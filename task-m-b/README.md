Task Manager Backend (Node.js)
==============================

This is the backend API for the Task Manager project. It is built with Node.js, Express, and MongoDB, and provides
endpoints for authentication, users, tasks, and projects. The legacy frontend in `task-m-f` is supported out of the box.

Features
--------
- User registration and login
- Admin area for managing users and projects
- Task CRUD with user assignment and optional project linking
- REST API responses compatible with the existing frontend

Tech Stack
----------
- Node.js + Express
- MongoDB + Mongoose
- JWT (issued on login, not required by the legacy frontend)

Setup
-----
1) Install dependencies:
   - `npm install`
2) Create `.env`:
   - `MONGODB_URI=<your mongo connection string>`
   - `JWT_SECRET=<your secret>`
   - `JWT_EXPIRES_IN=1d` (optional)
   - `BCRYPT_SALT_ROUNDS=10` (optional)
3) Run the server:
   - `npm run dev` or `npm start`

Default URL: `http://localhost:3000`

API Overview
------------
Auth
- `POST /auth/register`
- `POST /auth/login`

Users (frontend compatibility)
- `GET /users/all`
- `GET /users/:userName`
- `POST /users`
- `DELETE /users/:userName`

Tasks
- `GET /tasks`
- `GET /tasks/all`
- `GET /tasks/by-user/:userId`
- `GET /tasks/:id`
- `POST /tasks`
- `PUT /tasks/:id`
- `DELETE /tasks/:id`

Projects
- `GET /projects`
- `GET /projects/:id`
- `POST /projects`
- `PUT /projects/:id`
- `DELETE /projects/:id`

Formatters
----------
The backend uses utility formatters (see `src/utils/formatters.js`) to ensure consistent API responses:
- `toUserResponse(user)`: Converts a user document to a clean API response (id, user_name, email, role, created_at, updated_at).
- `toProjectResponse(project)`: Converts a project document to a clean API response (id, name, description, created_at, updated_at).
- `toTaskResponse(task)`: Converts a task document to a clean API response (id, user_id, project_id, title, description, status, priority, deadline, created_at, updated_at).

These formatters:
- Convert MongoDB's `_id` to `id` and use snake_case keys
- Handle both populated and unpopulated references (user/project)
- Hide internal fields from API responses
- Are used before sending any user, project, or task data to the frontend

Screenshots
-----------

**User Registration and Login:**
![Register and Login](../screenshots/register-login.png)

**Project Creation:**
![Create Project](../screenshots/create-project.png)

**Task Creation with Project Selection:**
![Create Task with Project](../screenshots/create-task.png)

**Frontend Home:**
![Frontend Home](../screenshots/frontend-home.png)

**User List:**
![User List](../screenshots/user-list.png)

**Project List:**
![Project List](../screenshots/project-list.png)
