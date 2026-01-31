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
