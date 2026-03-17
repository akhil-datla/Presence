// Profile view: user profile management

registerRoute('/profile', async () => {
  app.innerHTML = renderShell(`
    <a class="back-link" onclick="navigateTo('#/dashboard')">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"/><polyline points="12 19 5 12 12 5"/></svg>
      Back to Sessions
    </a>
    <div id="profile-content">
      <div class="card profile-card">
        <div class="profile-header">
          <div class="skeleton skeleton-avatar" style="width:64px;height:64px"></div>
          <div>
            <div class="skeleton skeleton-title" style="width:150px"></div>
            <div class="skeleton skeleton-text" style="width:200px"></div>
          </div>
        </div>
      </div>
    </div>
  `);

  try {
    const user = await api.getProfile();
    setUser(user);
    renderProfile(user);
  } catch (err) {
    showToast(err.message, 'error');
  }
});

function renderProfile(user) {
  const initials = (user.first_name[0] + user.last_name[0]).toUpperCase();

  document.getElementById('profile-content').innerHTML = `
    <div class="card profile-card">
      <div class="profile-header">
        <div class="avatar">${initials}</div>
        <div class="profile-name">
          <h3>${escapeHtml(user.first_name)} ${escapeHtml(user.last_name)}</h3>
          <p>${escapeHtml(user.email)}</p>
        </div>
      </div>

      <form id="profile-form">
        <div class="form-row">
          <div class="form-group">
            <label for="first_name">First Name</label>
            <input type="text" id="first_name" value="${escapeHtml(user.first_name)}" required>
          </div>
          <div class="form-group">
            <label for="last_name">Last Name</label>
            <input type="text" id="last_name" value="${escapeHtml(user.last_name)}" required>
          </div>
        </div>
        <div class="form-group">
          <label for="email">Email</label>
          <input type="email" id="email" value="${escapeHtml(user.email)}" required>
        </div>
        <div class="form-group">
          <label for="password">New Password <span style="font-weight:400;color:var(--color-text-muted)">(leave blank to keep current)</span></label>
          <input type="password" id="password" placeholder="At least 8 characters" minlength="8">
        </div>
        <button type="submit" class="btn btn-primary" id="save-btn">Save Changes</button>
      </form>

      <div class="danger-zone">
        <h4>Danger Zone</h4>
        <p>Permanently delete your account and all associated data. This cannot be undone.</p>
        <button class="btn btn-danger btn-sm" onclick="confirmDeleteAccount()">Delete Account</button>
      </div>
    </div>
  `;

  document.getElementById('profile-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const btn = document.getElementById('save-btn');
    btn.disabled = true;
    btn.textContent = 'Saving...';

    try {
      const updates = {};
      const firstName = document.getElementById('first_name').value;
      const lastName = document.getElementById('last_name').value;
      const email = document.getElementById('email').value;
      const password = document.getElementById('password').value;

      if (firstName !== user.first_name) updates.first_name = firstName;
      if (lastName !== user.last_name) updates.last_name = lastName;
      if (email !== user.email) updates.email = email;
      if (password) updates.password = password;

      if (Object.keys(updates).length === 0) {
        showToast('No changes to save', 'info');
        btn.disabled = false;
        btn.textContent = 'Save Changes';
        return;
      }

      const updated = await api.updateProfile(updates);
      setUser(updated);
      showToast('Profile updated!', 'success');
      renderProfile(updated);
    } catch (err) {
      showToast(err.message, 'error');
      btn.disabled = false;
      btn.textContent = 'Save Changes';
    }
  });
}

function confirmDeleteAccount() {
  const overlay = document.createElement('div');
  overlay.className = 'modal-overlay';
  overlay.onclick = (e) => { if (e.target === overlay) overlay.remove(); };
  overlay.innerHTML = `
    <div class="modal">
      <h3>Delete Account</h3>
      <p>This will permanently delete your account, all your sessions, and all attendance records. This action <strong>cannot be undone</strong>.</p>
      <div class="form-group" style="margin-top:1rem">
        <label for="confirm-delete">Type <strong>DELETE</strong> to confirm</label>
        <input type="text" id="confirm-delete" placeholder="DELETE">
      </div>
      <div class="modal-actions">
        <button class="btn btn-outline" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
        <button class="btn btn-danger" id="delete-account-btn" disabled>Delete My Account</button>
      </div>
    </div>
  `;
  document.body.appendChild(overlay);

  document.getElementById('confirm-delete').addEventListener('input', (e) => {
    document.getElementById('delete-account-btn').disabled = e.target.value !== 'DELETE';
  });

  document.getElementById('delete-account-btn').addEventListener('click', async () => {
    try {
      await api.deleteProfile();
      overlay.remove();
      api.setToken(null);
      localStorage.removeItem('user');
      showToast('Account deleted', 'success');
      navigateTo('#/login');
    } catch (err) {
      showToast(err.message, 'error');
    }
  });
}
