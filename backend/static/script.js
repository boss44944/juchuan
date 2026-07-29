let deviceId = localStorage.getItem('juchuan_device_id');
if (!deviceId) {
    deviceId = 'device-' + Math.random().toString(36).slice(2);
    localStorage.setItem('juchuan_device_id', deviceId);
}

const state = {
    authed: false,
    ws: null,
    historyPage: 1,
    historySize: 50,
};

function $(id) {
    return document.getElementById(id);
}

function setHint(id, text, isError) {
    const el = $(id);
    if (!el) return;
    el.textContent = text || '';
    el.classList.toggle('error', !!isError);
}

function escapeHTML(text) {
    return String(text || '')
        .replaceAll('&', '&amp;')
        .replaceAll('<', '&lt;')
        .replaceAll('>', '&gt;')
        .replaceAll('"', '&quot;')
        .replaceAll("'", '&#39;');
}

async function apiFetch(url, options) {
    const res = await fetch(url, {
        credentials: 'include',
        ...options,
    });

    if (res.status === 401) {
        state.authed = false;
        toggleView(false);
        throw new Error('未登录或登录已过期');
    }

    return res;
}

function activateTab(name) {
    document.querySelectorAll('.tab').forEach((btn) => {
        btn.classList.toggle('is-active', btn.dataset.tab === name);
    });
    document.querySelectorAll('.tab-panel').forEach((panel) => {
        panel.classList.toggle('is-active', panel.id === 'panel-' + name);
    });
}

function toggleView(isAuthed) {
    $('loginView').hidden = !!isAuthed;
    $('appView').hidden = !isAuthed;
}

function setAddressAndQR() {
    const url = location.origin + '/';

    const address = $('addressLink');
    if (address) {
        address.href = url;
        address.textContent = url;
    }

    const qrURL = '/api/qr?url=' + encodeURIComponent(url) + '&t=' + Date.now();
    const qr = $('qrImage');
    if (qr) {
        qr.src = qrURL;
    }

    const loginQr = $('loginQr');
    if (loginQr) {
        loginQr.src = qrURL;
    }
}

async function refreshAuthStatus() {
    const res = await fetch('/api/auth/status', { credentials: 'include' });
    if (!res.ok) {
        throw new Error('认证状态检查失败');
    }

    const data = await res.json();
    state.authed = !data.requires_password || data.authenticated;
    toggleView(state.authed);
}

async function login(password) {
    const res = await fetch('/api/auth/login', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ password }),
    });

    if (!res.ok) {
        throw new Error('密码错误，请重试');
    }
}

async function logout() {
    await fetch('/api/auth/logout', { method: 'POST', credentials: 'include' });
    state.authed = false;
    if (state.ws) {
        state.ws.close();
        state.ws = null;
    }
    toggleView(false);
}

async function registerDevice() {
    await apiFetch('/api/device/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: deviceId, name: navigator.userAgent }),
    });
}

function connectWS() {
    if (state.ws) {
        state.ws.close();
        state.ws = null;
    }

    const protocol = location.protocol === 'https:' ? 'wss://' : 'ws://';
    const ws = new WebSocket(protocol + location.host + '/ws?device=' + encodeURIComponent(deviceId));
    state.ws = ws;

    ws.onmessage = function (e) {
        let msg;
        try {
            msg = JSON.parse(e.data);
        } catch {
            return;
        }

        if (msg.type === 'file' && msg.url) {
            const confirmText = msg.filename ? ('收到文件: ' + msg.filename + '，是否下载？') : '收到文件，是否下载？';
            if (confirm(confirmText)) {
                location.href = msg.url;
            }
        }
    };

    ws.onclose = function () {
        if (state.authed) {
            setTimeout(connectWS, 2000);
        }
    };
}

async function loadDevices() {
    const res = await apiFetch('/api/devices');
    const body = await res.json();
    const list = Array.isArray(body) ? body : (body.data || []);
    const box = $('devices');

    if (!Array.isArray(list) || list.length === 0) {
        box.innerHTML = '<div class="empty">暂无在线设备</div>';
        return;
    }

    box.innerHTML = list
        .map((d) => '<div class="device-item">' + escapeHTML(d.Name || d.ID || '未知设备') + '</div>')
        .join('');
}

function formatTime(iso) {
    const d = new Date(iso);
    if (Number.isNaN(d.getTime())) return '';
    return d.toLocaleString();
}

async function loadHistory() {
    const res = await apiFetch('/api/history?page=' + state.historyPage + '&size=' + state.historySize);
    const body = await res.json();
    const list = Array.isArray(body) ? body : (body.data || []);
    const box = $('history');

    if (!Array.isArray(list) || list.length === 0) {
        box.innerHTML = '<div class="empty">暂无历史记录</div>';
        return;
    }

    box.innerHTML = list
        .map((h) => {
            const when = formatTime(h.created_at);
            if (h.type === 'file') {
                return '<article class="history-item">'
                    + '<div class="history-main"><strong>' + escapeHTML(h.filename || '未命名文件') + '</strong><span>' + escapeHTML(when) + '</span></div>'
                    + '<div class="history-meta">文件传输</div>'
                    + '<a class="download-link" href="/download/' + h.id + '">下载</a>'
                    + '</article>';
            }

            return '<article class="history-item">'
                + '<div class="history-main"><strong>文字消息</strong><span>' + escapeHTML(when) + '</span></div>'
                + '<div class="history-text">' + escapeHTML(h.content || '') + '</div>'
                + '</article>';
        })
        .join('');
}

async function sendText() {
    const content = $('text').value.trim();
    if (!content) {
        setHint('textStatus', '请输入内容后再发送', true);
        return;
    }

    const res = await apiFetch('/api/text', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content }),
    });

    if (!res.ok) {
        throw new Error('发送失败');
    }

    $('text').value = '';
    setHint('textStatus', '发送成功');
    await loadHistory();
}

async function uploadBlob(blob, name) {
    const form = new FormData();
    form.append('file', blob, name || 'clipboard.png');

    const res = await apiFetch('/upload', { method: 'POST', body: form });
    if (!res.ok) {
        throw new Error('上传失败');
    }

    setHint('uploadStatus', '上传成功');
    await loadHistory();
}

async function uploadFile(file) {
    if (!file) return;
    await uploadBlob(file, file.name);
}

function bindEvents() {
    document.querySelectorAll('.tab').forEach((btn) => {
        btn.addEventListener('click', () => activateTab(btn.dataset.tab));
    });

    $('copyUrlBtn').addEventListener('click', async () => {
        await navigator.clipboard.writeText(location.origin + '/');
        setHint('textStatus', '地址已复制');
    });

    $('sendTextBtn').addEventListener('click', () => {
        sendText().catch((err) => setHint('textStatus', err.message, true));
    });

    $('file').addEventListener('change', async (e) => {
        try {
            await uploadFile(e.target.files && e.target.files[0]);
        } catch (err) {
            setHint('uploadStatus', err.message, true);
        }
    });

    document.addEventListener('dragover', (e) => e.preventDefault());
    document.addEventListener('drop', (e) => {
        if (!state.authed) return;
        e.preventDefault();
        const file = e.dataTransfer && e.dataTransfer.files && e.dataTransfer.files[0];
        if (!file) return;
        uploadBlob(file, file.name).catch((err) => setHint('uploadStatus', err.message, true));
    });

    document.addEventListener('paste', (e) => {
        if (!state.authed) return;
        const items = e.clipboardData ? e.clipboardData.items : [];
        for (const item of items) {
            if (item.type && item.type.startsWith('image/')) {
                uploadBlob(item.getAsFile(), 'clipboard.png').catch((err) => setHint('uploadStatus', err.message, true));
            }
        }
    });

    $('refreshHistoryBtn').addEventListener('click', () => {
        loadHistory().catch(() => {});
    });

    $('loginForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        setHint('loginError', '');
        try {
            await login($('loginPassword').value);
            $('loginPassword').value = '';
            await startProtectedArea();
            toggleView(true);
        } catch (err) {
            setHint('loginError', err.message, true);
        }
    });

    $('logoutBtn').addEventListener('click', () => {
        logout().catch(() => {});
    });
}

async function startProtectedArea() {
    state.authed = true;
    activateTab('send');
    await registerDevice();
    connectWS();
    await Promise.all([loadDevices(), loadHistory()]);
}

async function init() {
    bindEvents();
    setAddressAndQR();

    await refreshAuthStatus();
    if (state.authed) {
        await startProtectedArea();
    }
}

init().catch((err) => {
    console.error(err);
    setHint('loginError', '初始化失败，请刷新页面重试', true);
});
