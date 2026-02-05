// Injects a shared navbar and footer into pages
(function () {
  // Ensure jQuery is available for animated theme toggling
  if (typeof window.jQuery === 'undefined') {
    const s = document.createElement('script');
    s.src = 'https://code.jquery.com/jquery-3.7.1.min.js';
    s.async = true;
    document.head.appendChild(s);
  }

  const user = (window.tmGetUser ? window.tmGetUser() : JSON.parse(localStorage.getItem('tm_user') || 'null'));
  const isAuthed = !!user;
  const isAdmin = !!(user && user.role === 'admin');

  const nav = `
  <nav class="navbar navbar-expand-lg bg-body-tertiary tm-navbar">
    <div class="container">
      <a class="navbar-brand fw-bold" href="index.html"><span class="brand-gradient">TaskManager</span></a>
      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#tmNav" aria-controls="tmNav" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="tmNav">
        <ul class="navbar-nav me-auto mb-2 mb-lg-0">
          ${isAuthed ? '<li class="nav-item"><a class="nav-link" href="index.html">Home</a></li>' : ''}
          ${isAuthed ? '<li class="nav-item"><a class="nav-link" href="tasks.html">Tasks</a></li>' : ''}
          ${isAuthed ? '<li class="nav-item"><a class="nav-link" href="completed.html">Completed</a></li>' : ''}
          ${isAdmin ? '<li class="nav-item"><a class="nav-link" href="admin.html">📊 Admin</a></li>' : ''}
          ${isAdmin ? '<li class="nav-item"><a class="nav-link" href="users.html">Users</a></li>' : ''}
          ${isAdmin ? '<li class="nav-item"><a class="nav-link" href="projects.html">Projects</a></li>' : ''}
          ${isAuthed ? '<li class="nav-item"><a class="nav-link" href="settings.html">Settings</a></li>' : ''}
        </ul>
        <div class="d-flex gap-2 align-items-center">
          ${isAuthed ? `<span class="text-muted small">Signed in as</span> <span class="badge text-bg-primary">${user.user_name}</span>` : ''}
          ${isAuthed ? `<button id="tmLogoutBtn" class="btn btn-outline-danger btn-sm">Logout</button>` : `
            <a class="btn btn-outline-primary btn-sm" href="login.html">Login</a>
            <a class="btn btn-primary btn-sm" href="register.html">Register</a>
          `}
        </div>
      </div>
    </div>
  </nav>`;

  const footer = `
  <footer class="py-4 mt-5 tm-footer">
    <div class="container d-flex flex-column flex-md-row justify-content-between align-items-center gap-2">
      <div class="small">© ${new Date().getFullYear()} TaskManager · Astana</div>
      <div class="small">
        Owner: Samat Kosmurat ·
        <a class="tm-link" href="https://t.me/qosmurat" target="_blank" rel="noopener">Telegram</a> ·
        <a class="tm-link" href="settings.html">Toggle Theme</a>
      </div>
    </div>
  </footer>`;

  const headerMount = document.getElementById('app-header');
  const footerMount = document.getElementById('app-footer');
  if (headerMount) headerMount.innerHTML = nav;
  if (footerMount) footerMount.innerHTML = footer;

  // Logout handler
  document.addEventListener('click', (e) => {
    const el = e.target;
    if (el && el.id === 'tmLogoutBtn') {
      e.preventDefault();
      localStorage.removeItem('tm_user');
      window.name = '';
      window.location.href = 'index.html';
    }
  });
})();
