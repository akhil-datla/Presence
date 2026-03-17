// Dashboard view: Session list and creation

registerRoute('/dashboard', async () => {
  const user = getUser();
  const greeting = user ? `${getGreeting()}, ${escapeHtml(user.first_name)}` : getGreeting();

  app.innerHTML = renderShell(`
    <div class="page-header">
      <div>
        <h2>Your Sessions</h2>
        <p class="greeting">${greeting}</p>
      </div>
      <button class="btn btn-primary" onclick="showCreateSessionModal()">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        New Session
      </button>
    </div>
    <div id="sessions-list">
      <div class="sessions-grid">
        <div class="skeleton skeleton-card"></div>
        <div class="skeleton skeleton-card"></div>
        <div class="skeleton skeleton-card"></div>
      </div>
    </div>
  `);

  await loadSessions();
});

async function loadSessions() {
  try {
    const sessions = await api.listSessions();
    const container = document.getElementById('sessions-list');

    if (!sessions || sessions.length === 0) {
      container.innerHTML = `
        <div class="empty-state">
          <div class="empty-state-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
              <line x1="16" y1="2" x2="16" y2="6"/>
              <line x1="8" y1="2" x2="8" y2="6"/>
              <line x1="3" y1="10" x2="21" y2="10"/>
            </svg>
          </div>
          <h3>No sessions yet</h3>
          <p>Create your first attendance session to start tracking.</p>
          <button class="btn btn-primary" onclick="showCreateSessionModal()">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
            New Session
          </button>
        </div>
      `;
      return;
    }

    container.innerHTML = `
      <div class="sessions-grid">
        ${sessions.map(s => `
          <div class="card card-clickable session-card" onclick="navigateTo('#/session/${s.id}')">
            <div class="session-info">
              <h3>${escapeHtml(s.name)}</h3>
              <p>Created ${formatDate(s.created_at)}</p>
            </div>
            <div class="session-actions">
              <button class="btn btn-ghost btn-sm" onclick="event.stopPropagation(); showEditSessionModal('${s.id}', '${escapeHtml(s.name)}')" title="Edit">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
              </button>
              <button class="btn btn-ghost btn-sm" onclick="event.stopPropagation(); confirmDeleteSession('${s.id}', '${escapeHtml(s.name)}')" title="Delete" style="color:var(--color-danger)">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
              </button>
            </div>
          </div>
        `).join('')}
      </div>
    `;
  } catch (err) {
    showToast(err.message, 'error');
  }
}

function showCreateSessionModal() {
  const overlay = document.createElement('div');
  overlay.className = 'modal-overlay';
  overlay.onclick = (e) => { if (e.target === overlay) overlay.remove(); };
  overlay.innerHTML = `
    <div class="modal">
      <h3>Create New Session</h3>
      <form id="create-session-form">
        <div class="form-group">
          <label for="session-name">Session Name</label>
          <input type="text" id="session-name" placeholder="e.g. Team Standup, CS 101 Lecture" required autofocus>
        </div>
        <div class="modal-actions">
          <button type="button" class="btn btn-outline" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
          <button type="submit" class="btn btn-primary">Create Session</button>
        </div>
      </form>
    </div>
  `;
  document.body.appendChild(overlay);

  document.getElementById('create-session-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const name = document.getElementById('session-name').value;
    try {
      await api.createSession({ name });
      overlay.remove();
      showToast('Session created!', 'success');
      await loadSessions();
    } catch (err) {
      showToast(err.message, 'error');
    }
  });
}

function showEditSessionModal(id, currentName) {
  const overlay = document.createElement('div');
  overlay.className = 'modal-overlay';
  overlay.onclick = (e) => { if (e.target === overlay) overlay.remove(); };
  overlay.innerHTML = `
    <div class="modal">
      <h3>Edit Session</h3>
      <form id="edit-session-form">
        <div class="form-group">
          <label for="edit-session-name">Session Name</label>
          <input type="text" id="edit-session-name" value="${currentName}" required autofocus>
        </div>
        <div class="modal-actions">
          <button type="button" class="btn btn-outline" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
          <button type="submit" class="btn btn-primary">Save Changes</button>
        </div>
      </form>
    </div>
  `;
  document.body.appendChild(overlay);

  document.getElementById('edit-session-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const name = document.getElementById('edit-session-name').value;
    try {
      await api.updateSession(id, { name });
      overlay.remove();
      showToast('Session updated!', 'success');
      await loadSessions();
    } catch (err) {
      showToast(err.message, 'error');
    }
  });
}

function confirmDeleteSession(id, name) {
  const overlay = document.createElement('div');
  overlay.className = 'modal-overlay';
  overlay.onclick = (e) => { if (e.target === overlay) overlay.remove(); };
  overlay.innerHTML = `
    <div class="modal">
      <h3>Delete Session</h3>
      <p>Are you sure you want to delete <strong>${name}</strong>? This will also delete all attendance records. This action cannot be undone.</p>
      <div class="modal-actions">
        <button class="btn btn-outline" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
        <button class="btn btn-danger" id="confirm-delete-btn">Delete Session</button>
      </div>
    </div>
  `;
  document.body.appendChild(overlay);

  document.getElementById('confirm-delete-btn').addEventListener('click', async () => {
    try {
      await api.deleteSession(id);
      overlay.remove();
      showToast('Session deleted', 'success');
      await loadSessions();
    } catch (err) {
      showToast(err.message, 'error');
    }
  });
}
