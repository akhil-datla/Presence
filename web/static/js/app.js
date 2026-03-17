// SPA Router and App Shell
const app = document.getElementById('app');

// Toast icons
const TOAST_ICONS = {
  success: '<svg class="toast-icon" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/></svg>',
  error: '<svg class="toast-icon" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"/></svg>',
  info: '<svg class="toast-icon" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/></svg>',
};

function showToast(message, type = 'info') {
  let container = document.querySelector('.toast-container');
  if (!container) {
    container = document.createElement('div');
    container.className = 'toast-container';
    document.body.appendChild(container);
  }

  const toast = document.createElement('div');
  toast.className = `toast toast-${type}`;
  toast.innerHTML = `${TOAST_ICONS[type] || TOAST_ICONS.info}<span>${message}</span>`;
  container.appendChild(toast);

  setTimeout(() => toast.remove(), 3000);
}

// Simple router
const routes = {};

function registerRoute(hash, handler) {
  routes[hash] = handler;
}

function navigateTo(hash) {
  window.location.hash = hash;
}

async function handleRoute() {
  const hash = window.location.hash || '#/login';
  const [path, ...paramParts] = hash.slice(1).split('/').filter(Boolean);

  if (!api.isAuthenticated() && path !== 'login' && path !== 'register') {
    window.location.hash = '#/login';
    return;
  }

  if (api.isAuthenticated() && (path === 'login' || path === 'register')) {
    window.location.hash = '#/dashboard';
    return;
  }

  const fullPath = '/' + path;

  if (routes[fullPath]) {
    try {
      await routes[fullPath](paramParts);
      // Add page enter animation
      const main = app.querySelector('main') || app.querySelector('.auth-wrapper');
      if (main) main.classList.add('page-enter');
    } catch (err) {
      showToast(err.message, 'error');
    }
  } else {
    app.innerHTML = '<div class="container"><h2>Page not found</h2></div>';
  }
}

window.addEventListener('hashchange', handleRoute);

// Render app shell for authenticated pages
function renderShell(content) {
  const user = getUser();
  const initials = user ? (user.first_name[0] + user.last_name[0]).toUpperCase() : '?';

  return `
    <header class="app-header">
      <a href="#/dashboard" class="logo">
        <span class="logo-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <polyline points="16 11 18 13 22 9"/>
          </svg>
        </span>
        Presence
      </a>
      <nav class="header-nav">
        <a href="#/dashboard">Sessions</a>
        <a href="#/profile">Profile</a>
        <a href="#" onclick="logout(); return false;">Logout</a>
        <div class="nav-avatar" onclick="navigateTo('#/profile')" title="${user ? user.first_name + ' ' + user.last_name : ''}">${initials}</div>
      </nav>
    </header>
    <main class="container">${content}</main>
  `;
}

function logout() {
  api.setToken(null);
  localStorage.removeItem('user');
  window.location.hash = '#/login';
}

function setUser(user) {
  localStorage.setItem('user', JSON.stringify(user));
}

function getUser() {
  try {
    return JSON.parse(localStorage.getItem('user'));
  } catch {
    return null;
  }
}

function formatDate(dateStr) {
  if (!dateStr) return '\u2014';
  const d = new Date(dateStr);
  return d.toLocaleDateString('en-US', {
    month: 'short', day: 'numeric', year: 'numeric',
    hour: 'numeric', minute: '2-digit'
  });
}

function formatTime(dateStr) {
  if (!dateStr) return '\u2014';
  const d = new Date(dateStr);
  return d.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit', second: '2-digit' });
}

function getGreeting() {
  const h = new Date().getHours();
  if (h < 12) return 'Good morning';
  if (h < 17) return 'Good afternoon';
  return 'Good evening';
}

function escapeHtml(str) {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}

function getInitials(name) {
  return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2);
}

document.addEventListener('DOMContentLoaded', () => {
  handleRoute();
});
