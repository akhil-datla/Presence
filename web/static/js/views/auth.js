// Auth views: Login and Register with split layout

const AUTH_BRAND_PANEL = `
  <div class="auth-brand">
    <div class="auth-brand-content">
      <h1>Track attendance,<br>effortlessly.</h1>
      <p>A simple, self-hosted solution for managing attendance sessions. No complex setup, no external dependencies.</p>
      <div class="auth-features">
        <div class="auth-feature">
          <div class="auth-feature-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
          </div>
          Single binary — zero config needed
        </div>
        <div class="auth-feature">
          <div class="auth-feature-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
          </div>
          JWT authentication built-in
        </div>
        <div class="auth-feature">
          <div class="auth-feature-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          </div>
          Export attendance as CSV
        </div>
        <div class="auth-feature">
          <div class="auth-feature-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
          </div>
          RESTful API for integrations
        </div>
      </div>
    </div>
  </div>
`;

registerRoute('/login', async () => {
  app.innerHTML = `
    <div class="auth-wrapper">
      ${AUTH_BRAND_PANEL}
      <div class="auth-form-side">
        <div class="auth-card">
          <h2>Welcome back</h2>
          <p class="subtitle">Sign in to your account to continue</p>
          <form id="login-form">
            <div class="form-group">
              <label for="email">Email</label>
              <input type="email" id="email" placeholder="you@example.com" required autofocus>
            </div>
            <div class="form-group">
              <label for="password">Password</label>
              <input type="password" id="password" placeholder="Enter your password" required>
            </div>
            <button type="submit" class="btn btn-primary btn-block" id="login-btn">Sign In</button>
          </form>
          <div class="auth-toggle">
            Don't have an account? <a onclick="navigateTo('#/register')">Create one</a>
          </div>
        </div>
      </div>
    </div>
  `;

  document.getElementById('login-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const btn = document.getElementById('login-btn');
    btn.disabled = true;
    btn.textContent = 'Signing in...';

    try {
      const data = await api.login({
        email: document.getElementById('email').value,
        password: document.getElementById('password').value,
      });
      api.setToken(data.token);
      setUser(data.user);
      showToast('Welcome back!', 'success');
      navigateTo('#/dashboard');
    } catch (err) {
      showToast(err.message, 'error');
      btn.disabled = false;
      btn.textContent = 'Sign In';
    }
  });
});

registerRoute('/register', async () => {
  app.innerHTML = `
    <div class="auth-wrapper">
      ${AUTH_BRAND_PANEL}
      <div class="auth-form-side">
        <div class="auth-card">
          <h2>Create account</h2>
          <p class="subtitle">Get started with Presence in seconds</p>
          <form id="register-form">
            <div class="form-row">
              <div class="form-group">
                <label for="first_name">First Name</label>
                <input type="text" id="first_name" placeholder="Jane" required autofocus>
              </div>
              <div class="form-group">
                <label for="last_name">Last Name</label>
                <input type="text" id="last_name" placeholder="Doe" required>
              </div>
            </div>
            <div class="form-group">
              <label for="email">Email</label>
              <input type="email" id="email" placeholder="you@example.com" required>
            </div>
            <div class="form-group">
              <label for="password">Password</label>
              <input type="password" id="password" placeholder="At least 8 characters" required minlength="8">
            </div>
            <button type="submit" class="btn btn-primary btn-block" id="register-btn">Create Account</button>
          </form>
          <div class="auth-toggle">
            Already have an account? <a onclick="navigateTo('#/login')">Sign in</a>
          </div>
        </div>
      </div>
    </div>
  `;

  document.getElementById('register-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const btn = document.getElementById('register-btn');
    btn.disabled = true;
    btn.textContent = 'Creating account...';

    try {
      const data = await api.register({
        first_name: document.getElementById('first_name').value,
        last_name: document.getElementById('last_name').value,
        email: document.getElementById('email').value,
        password: document.getElementById('password').value,
      });
      api.setToken(data.token);
      setUser(data.user);
      showToast('Account created!', 'success');
      navigateTo('#/dashboard');
    } catch (err) {
      showToast(err.message, 'error');
      btn.disabled = false;
      btn.textContent = 'Create Account';
    }
  });
});
