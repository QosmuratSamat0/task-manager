// Simple API client for the Go backend
const TM_API = (function () {
  const DEFAULT_BASE = 'http://localhost:8080';

  const baseURL = () => {
    return localStorage.getItem('tm_api_base') || DEFAULT_BASE;
  };

  const headers = { 'Content-Type': 'application/json' };

  async function http(method, path, body) {
    const url = baseURL() + path;
    const opts = { method, headers };
    if (body) opts.body = JSON.stringify(body);
    const res = await fetch(url, opts);
    const text = await res.text();
    let data = null;
    try { data = text ? JSON.parse(text) : null; } catch (_) { data = text; }
    if (!res.ok) {
      const err = new Error((data && data.error) || (data && data.status) || res.statusText || 'Request failed');
      err.status = res.status;
      err.data = data;
      throw err;
    }
    return data;
  }

  // Users
  const getUser = (userName) => http('GET', `/users/${encodeURIComponent(userName)}`);
  const listUsers = () => http('GET', '/users/all');
  const createUser = (payload) => http('POST', '/users/', payload);
  const deleteUser = (userName) => http('DELETE', `/users/${encodeURIComponent(userName)}`);

  // Auth
  const login = (userName, password) => http('POST', '/auth/login', { user_name: userName, password });

  // Tasks
  const createTask = (payload) => http('POST', '/tasks/', payload);
  const listAllTasks = () => http('GET', '/tasks/all');
  const getTaskById = (id) => http('GET', `/tasks/${id}`); // backend param named user_id, but it is task id
  const deleteTaskById = (id) => http('DELETE', `/tasks/${id}`);
  const listTasksByUser = (userId) => http('GET', `/tasks/by-user/${encodeURIComponent(userId)}`);
  const updateTask = (id, payload) => http('PUT', `/tasks/${encodeURIComponent(id)}`, payload);

  return {
    baseURL,
    getUser,
    listUsers,
    createUser,
    deleteUser,
    login,
    createTask,
    listAllTasks,
    listTasksByUser,
    getTaskById,
    deleteTaskById,
    updateTask,
  };
})();
