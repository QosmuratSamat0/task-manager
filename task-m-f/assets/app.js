// Global app script: theme, helpers, validation
(function () {
  // Theme init from localStorage
  const theme = localStorage.getItem('tm_theme') || 'light';
  document.documentElement.setAttribute('data-theme', theme);

  // Theme animation mode: 'crossfade' | 'radial' | 'none'
  let themeAnim = localStorage.getItem('tm_theme_anim') || 'radial';
  window.tmSetThemeAnimation = function (mode) {
    themeAnim = (mode === 'crossfade' || mode === 'radial' || mode === 'none') ? mode : 'radial';
    localStorage.setItem('tm_theme_anim', themeAnim);
    return themeAnim;
  };

  // Remember last click position to use as origin for radial reveal
  let lastClick = { x: null, y: null };
  document.addEventListener('click', function (e) {
    lastClick = { x: e.clientX, y: e.clientY };
  }, true);

  // Theme toggle handler with selectable animation
  window.tmToggleTheme = function () {
    const cur = document.documentElement.getAttribute('data-theme') || 'light';
    const next = cur === 'light' ? 'dark' : 'light';

    if (themeAnim === 'radial') {
      // Radial reveal: shrink overlay of OLD theme to reveal the NEW theme underneath
      const oldBG = getComputedStyle(document.body).backgroundColor;
      // Switch to new theme under the overlay
      document.documentElement.setAttribute('data-theme', next);
      localStorage.setItem('tm_theme', next);

      const originX = (lastClick.x == null || lastClick.y == null)
        ? window.innerWidth / 2
        : lastClick.x;
      const originY = (lastClick.x == null || lastClick.y == null)
        ? window.innerHeight / 2
        : lastClick.y;

      const overlay = document.createElement('div');
      overlay.id = 'tm-theme-overlay-radial';
      Object.assign(overlay.style, {
        position: 'fixed',
        inset: '0',
        pointerEvents: 'none',
        zIndex: '2147483646',
        backgroundColor: oldBG,
        // start fully covered: large circle that covers viewport
        clipPath: `circle(150% at ${originX}px ${originY}px)`,
        WebkitClipPath: `circle(150% at ${originX}px ${originY}px)`,
        transition: 'clip-path 450ms ease-out',
        WebkitTransition: '-webkit-clip-path 450ms ease-out'
      });
      document.body.appendChild(overlay);
      // next frame: shrink to center
      requestAnimationFrame(() => {
        overlay.style.clipPath = `circle(0% at ${originX}px ${originY}px)`;
        overlay.style.WebkitClipPath = `circle(0% at ${originX}px ${originY}px)`;
      });
      overlay.addEventListener('transitionend', () => overlay.remove(), { once: true });
      return;
    }

    if (themeAnim === 'crossfade' && window.jQuery) {
      const $ = window.jQuery;
      const $overlay = $('<div/>')
        .attr('id', 'tm-theme-overlay')
        .css({
          position: 'fixed', inset: 0, pointerEvents: 'none', zIndex: 2147483646,
          backgroundColor: getComputedStyle(document.body).backgroundColor, opacity: 0
        });
      $('body').append($overlay);
      $overlay.animate({ opacity: 1 }, 150, function () {
        document.documentElement.setAttribute('data-theme', next);
        localStorage.setItem('tm_theme', next);
        $overlay.animate({ opacity: 0 }, 300, function () { $overlay.remove(); });
      });
      return;
    }

    // No animation
    document.documentElement.setAttribute('data-theme', next);
    localStorage.setItem('tm_theme', next);
  };

  // Bootstrap validation helper
  window.tmValidate = function (form) {
    if (!form) return false;
    if (!form.checkValidity()) {
      form.classList.add('was-validated');
      return false;
    }
    return true;
  };

  // Utility: format date for backend (RFC3339)
  window.tmToISO = function (value) {
    if (!value) return null;
    // value from <input type="datetime-local"> is local time without zone
    const dt = new Date(value);
    return dt.toISOString();
  };

  function readUserFromWindowName() {
    const raw = window.name || '';
    if (!raw || raw.indexOf('tm_user:') !== 0) return null;
    try {
      const payload = raw.slice('tm_user:'.length);
      return JSON.parse(payload);
    } catch (_) {
      return null;
    }
  }

  window.tmGetUser = function () {
    const fromStorage = localStorage.getItem('tm_user');
    if (fromStorage) {
      try { return JSON.parse(fromStorage); } catch (_) { /* ignore */ }
    }
    const fromWindow = readUserFromWindowName();
    if (fromWindow) {
      localStorage.setItem('tm_user', JSON.stringify(fromWindow));
      return fromWindow;
    }
    return null;
  };

  window.tmSetUser = function (user) {
    if (!user) return;
    localStorage.setItem('tm_user', JSON.stringify(user));
    window.name = 'tm_user:' + JSON.stringify(user);
  };

  // Utility: get current user or redirect
  window.tmRequireUser = function () {
    const user = window.tmGetUser ? window.tmGetUser() : null;
    if (!user) {
      window.location.href = 'login.html';
      return null;
    }
    return user;
  };

  // Toast notifications (Bootstrap-like alerts)
  function ensureToastRoot() {
    let root = document.getElementById('tm-toast-root');
    if (!root) {
      root = document.createElement('div');
      root.id = 'tm-toast-root';
      root.style.position = 'fixed';
      root.style.top = '1rem';
      root.style.right = '1rem';
      root.style.zIndex = '1080';
      document.body.appendChild(root);
    }
    return root;
  }

  window.tmToast = function (message, type = 'info', timeout = 2500) {
    const colors = {
      success: 'alert-success',
      error: 'alert-danger',
      warn: 'alert-warning',
      info: 'alert-primary',
    };
    const cls = colors[type] || colors.info;
    const root = ensureToastRoot();
    const el = document.createElement('div');
    el.className = `alert ${cls} shadow-sm mb-2`;
    el.textContent = message;
    root.appendChild(el);
    setTimeout(() => {
      el.remove();
    }, timeout);
  };
})();
