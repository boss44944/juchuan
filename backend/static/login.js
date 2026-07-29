function $(id) {
    return document.getElementById(id);
}

function setError(text) {
    const el = $('loginError');
    if (!el) return;
    el.textContent = text || '';
}

function setQR() {
    const url = location.origin + '/';
    const qrURL = '/api/qr?url=' + encodeURIComponent(url) + '&t=' + Date.now();
    $('loginQr').src = qrURL;
}

async function checkAuth() {
    const res = await fetch('/api/auth/status', { credentials: 'include' });
    if (!res.ok) return;
    const data = await res.json();
    if (!data.requires_password || data.authenticated) {
        location.replace('/app');
    }
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

function bindEvents() {
    $('loginForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        setError('');
        try {
            await login($('loginPassword').value);
            $('loginPassword').value = '';
            location.assign('/app');
        } catch (err) {
            setError(err.message || '登录失败');
        }
    });
}

async function init() {
    setQR();
    bindEvents();
    await checkAuth();
}

init().catch(() => setError('初始化失败，请刷新重试'));
