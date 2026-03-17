// API client for Presence backend
const API_BASE = '/api/v1';

class Api {
  constructor() {
    this.token = localStorage.getItem('token');
  }

  setToken(token) {
    this.token = token;
    if (token) {
      localStorage.setItem('token', token);
    } else {
      localStorage.removeItem('token');
    }
  }

  getToken() {
    return this.token;
  }

  isAuthenticated() {
    return !!this.token;
  }

  async request(method, path, body) {
    const headers = { 'Content-Type': 'application/json' };
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    const opts = { method, headers };
    if (body && method !== 'GET') {
      opts.body = JSON.stringify(body);
    }

    const res = await fetch(`${API_BASE}${path}`, opts);

    if (res.status === 401) {
      this.setToken(null);
      window.location.hash = '#/login';
      throw new Error('Session expired. Please log in again.');
    }

    // Handle CSV downloads
    if (res.headers.get('content-type')?.includes('text/csv')) {
      return res.blob();
    }

    const json = await res.json();

    if (!res.ok) {
      throw new Error(json.error || 'Something went wrong');
    }

    return json.data;
  }

  // Auth
  register(data) { return this.request('POST', '/auth/register', data); }
  login(data) { return this.request('POST', '/auth/login', data); }

  // User
  getProfile() { return this.request('GET', '/users/me'); }
  updateProfile(data) { return this.request('PUT', '/users/me', data); }
  deleteProfile() { return this.request('DELETE', '/users/me'); }

  // Sessions
  createSession(data) { return this.request('POST', '/sessions', data); }
  listSessions() { return this.request('GET', '/sessions'); }
  getSession(id) { return this.request('GET', `/sessions/${id}`); }
  updateSession(id, data) { return this.request('PUT', `/sessions/${id}`, data); }
  deleteSession(id) { return this.request('DELETE', `/sessions/${id}`); }

  // Attendance
  checkIn(sessionId) { return this.request('POST', `/sessions/${sessionId}/checkin`); }
  checkOut(sessionId) { return this.request('POST', `/sessions/${sessionId}/checkout`); }
  getAttendance(sessionId) { return this.request('GET', `/sessions/${sessionId}/attendance`); }
  clearAttendance(sessionId) { return this.request('DELETE', `/sessions/${sessionId}/attendance`); }
  filterAttendance(sessionId, mode, time) {
    return this.request('GET', `/sessions/${sessionId}/attendance/filter?mode=${mode}&time=${encodeURIComponent(time)}`);
  }

  async exportCSV(sessionId) {
    const headers = {};
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }
    const res = await fetch(`${API_BASE}/sessions/${sessionId}/export/csv`, { headers });
    if (!res.ok) {
      const json = await res.json().catch(() => ({}));
      throw new Error(json.error || 'Export failed');
    }
    const blob = await res.blob();
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `attendance_${sessionId}.csv`;
    a.click();
    URL.revokeObjectURL(url);
  }
}

window.api = new Api();
