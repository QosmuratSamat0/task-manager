Task Manager (Node.js + HTML)
=============================

This repository contains a simple Task Manager with a Node.js backend and a static HTML frontend.
The backend provides REST APIs for authentication, users, tasks, and projects, and the frontend
is a multi‑page Bootstrap UI that consumes those endpoints.

Folders
-------
- `task-m-b` — Node.js/Express backend (MongoDB + Mongoose)
- `task-m-f` — Static frontend (HTML/CSS/JS)

Features
--------
- User registration and login
- Admin pages for users and projects
- Task CRUD with user assignment
- Optional project linking for tasks
- Frontend compatible API responses

How to Run
----------
1) Backend:
   - `cd task-m-b`
   - `npm install`
   - create `.env` with:
     - `MONGODB_URI`
     - `JWT_SECRET`
   - `npm run dev`
2) Frontend:
   - open `task-m-f/index.html` in the browser
   - set the API base URL to `http://localhost:3000` in Settings if needed

Default Backend URL: `http://localhost:3000`

Notes
-----
- Admin navigation (Users/Projects) appears only when the logged‑in user has `role: "admin"`.
- The frontend stores session data in localStorage and in `window.name` to work when opened via `file://`.
