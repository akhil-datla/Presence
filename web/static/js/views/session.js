// Session detail view: attendance management

registerRoute('/session', async (params) => {
  const sessionId = params[0];
  if (!sessionId) { navigateTo('#/dashboard'); return; }

  app.innerHTML = renderShell(`
    <a class="back-link" onclick="navigateTo('#/dashboard')">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"/><polyline points="12 19 5 12 12 5"/></svg>
      Back to Sessions
    </a>
    <div id="session-content">
      <div class="skeleton skeleton-title"></div>
      <div class="stats-bar">
        <div class="skeleton skeleton-card" style="height:72px"></div>
        <div class="skeleton skeleton-card" style="height:72px"></div>
        <div class="skeleton skeleton-card" style="height:72px"></div>
      </div>
      <div class="skeleton skeleton-card" style="height:200px"></div>
    </div>
  `);

  try {
    const session = await api.getSession(sessionId);
    await renderSessionDetail(session);
  } catch (err) {
    document.getElementById('session-content').innerHTML = `
      <div class="empty-state">
        <div class="empty-state-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
        </div>
        <h3>Session not found</h3>
        <p>${escapeHtml(err.message)}</p>
        <button class="btn btn-primary" onclick="navigateTo('#/dashboard')">Back to Dashboard</button>
      </div>
    `;
  }
});

async function renderSessionDetail(session) {
  const content = document.getElementById('session-content');

  content.innerHTML = `
    <div class="session-detail-header">
      <h2>${escapeHtml(session.name)}</h2>
      <div class="session-detail-actions">
        <button class="btn btn-success" id="checkin-btn" onclick="doCheckIn('${session.id}')">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
          Check In
        </button>
        <button class="btn btn-outline" id="checkout-btn" onclick="doCheckOut('${session.id}')">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
          Check Out
        </button>
        <button class="btn btn-outline" onclick="exportCSV('${session.id}')">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          CSV
        </button>
        <button class="btn btn-ghost btn-sm" style="color:var(--color-danger)" onclick="confirmClearAttendance('${session.id}')" title="Clear all records">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
        </button>
      </div>
    </div>

    <div id="stats-container"></div>

    <div class="filter-bar">
      <div class="form-group">
        <label for="filter-mode">Filter</label>
        <select id="filter-mode">
          <option value="">All records</option>
          <option value="before">Before</option>
          <option value="after">After</option>
        </select>
      </div>
      <div class="form-group" id="filter-time-group" style="display:none">
        <label for="filter-time">Date/Time</label>
        <input type="datetime-local" id="filter-time">
      </div>
      <button class="btn btn-outline btn-sm" id="filter-btn" style="display:none" onclick="applyFilter('${session.id}')">Apply Filter</button>
    </div>

    <div id="attendance-table"></div>
  `;

  document.getElementById('filter-mode').addEventListener('change', (e) => {
    const show = e.target.value !== '';
    document.getElementById('filter-time-group').style.display = show ? '' : 'none';
    document.getElementById('filter-btn').style.display = show ? '' : 'none';
    if (!show) loadAttendance(session.id);
  });

  await loadAttendance(session.id);
}

function renderStats(records) {
  const total = records ? records.length : 0;
  const checkedIn = records ? records.filter(r => !r.time_out).length : 0;
  const checkedOut = total - checkedIn;

  document.getElementById('stats-container').innerHTML = `
    <div class="stats-bar">
      <div class="stat-card">
        <div class="stat-icon stat-icon-primary">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
        </div>
        <div>
          <div class="stat-value">${total}</div>
          <div class="stat-label">Total Records</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon stat-icon-success">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
        </div>
        <div>
          <div class="stat-value">${checkedIn}</div>
          <div class="stat-label">Currently In</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon stat-icon-warm">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
        </div>
        <div>
          <div class="stat-value">${checkedOut}</div>
          <div class="stat-label">Checked Out</div>
        </div>
      </div>
    </div>
  `;
}

async function loadAttendance(sessionId) {
  try {
    const records = await api.getAttendance(sessionId);
    renderStats(records);
    renderAttendanceTable(records);
  } catch (err) {
    showToast(err.message, 'error');
  }
}

function renderAttendanceTable(records) {
  const container = document.getElementById('attendance-table');

  if (!records || records.length === 0) {
    container.innerHTML = `
      <div class="empty-state">
        <div class="empty-state-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <line x1="17" y1="11" x2="23" y2="11"/>
          </svg>
        </div>
        <h3>No attendance records</h3>
        <p>Click "Check In" to start tracking attendance for this session.</p>
      </div>
    `;
    return;
  }

  container.innerHTML = `
    <div class="table-wrapper">
      <table>
        <thead>
          <tr>
            <th>Participant</th>
            <th>Time In</th>
            <th>Time Out</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          ${records.map(r => {
            const initials = getInitials(r.participant_name);
            return `
            <tr>
              <td>
                <div class="participant-cell">
                  <div class="participant-avatar">${initials}</div>
                  <strong>${escapeHtml(r.participant_name)}</strong>
                </div>
              </td>
              <td>${formatTime(r.time_in)}</td>
              <td>${r.time_out ? formatTime(r.time_out) : '\u2014'}</td>
              <td>${r.time_out
                ? '<span class="badge badge-neutral"><span class="badge-dot"></span>Checked out</span>'
                : '<span class="badge badge-success"><span class="badge-dot"></span>Checked in</span>'
              }</td>
            </tr>`;
          }).join('')}
        </tbody>
      </table>
      <div class="table-footer">${records.length} record${records.length !== 1 ? 's' : ''}</div>
    </div>
  `;
}

async function doCheckIn(sessionId) {
  const btn = document.getElementById('checkin-btn');
  btn.disabled = true;
  try {
    await api.checkIn(sessionId);
    showToast('Checked in!', 'success');
    await loadAttendance(sessionId);
  } catch (err) {
    showToast(err.message, 'error');
  } finally {
    btn.disabled = false;
  }
}

async function doCheckOut(sessionId) {
  const btn = document.getElementById('checkout-btn');
  btn.disabled = true;
  try {
    await api.checkOut(sessionId);
    showToast('Checked out!', 'success');
    await loadAttendance(sessionId);
  } catch (err) {
    showToast(err.message, 'error');
  } finally {
    btn.disabled = false;
  }
}

async function exportCSV(sessionId) {
  try {
    await api.exportCSV(sessionId);
    showToast('CSV downloaded!', 'success');
  } catch (err) {
    showToast(err.message, 'error');
  }
}

async function applyFilter(sessionId) {
  const mode = document.getElementById('filter-mode').value;
  const time = document.getElementById('filter-time').value;
  if (!mode || !time) {
    showToast('Please select a filter mode and time', 'error');
    return;
  }
  try {
    const isoTime = new Date(time).toISOString();
    const records = await api.filterAttendance(sessionId, mode, isoTime);
    renderAttendanceTable(records);
  } catch (err) {
    showToast(err.message, 'error');
  }
}

function confirmClearAttendance(sessionId) {
  const overlay = document.createElement('div');
  overlay.className = 'modal-overlay';
  overlay.onclick = (e) => { if (e.target === overlay) overlay.remove(); };
  overlay.innerHTML = `
    <div class="modal">
      <h3>Clear All Attendance</h3>
      <p>Are you sure? This will permanently delete all attendance records for this session.</p>
      <div class="modal-actions">
        <button class="btn btn-outline" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
        <button class="btn btn-danger" id="confirm-clear-btn">Clear All Records</button>
      </div>
    </div>
  `;
  document.body.appendChild(overlay);

  document.getElementById('confirm-clear-btn').addEventListener('click', async () => {
    try {
      await api.clearAttendance(sessionId);
      overlay.remove();
      showToast('Attendance records cleared', 'success');
      await loadAttendance(sessionId);
    } catch (err) {
      showToast(err.message, 'error');
    }
  });
}
