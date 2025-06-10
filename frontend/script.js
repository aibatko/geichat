let token = null;
let currentGroupId = null;
let ws = null;
let signInMode = false;

document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('authForm').addEventListener('submit', handleAuth);
    document.getElementById('toggleAuth').addEventListener('click', toggleMode);
    document.getElementById('signOutBtn').addEventListener('click', signOut);
    document.getElementById('newGroupBtn').addEventListener('click', newGroup);
    document.getElementById('sendBtn').addEventListener('click', sendMessage);
    document.getElementById('inviteBtn').addEventListener('click', inviteUser);
    checkStoredAuth();
});

function toggleMode() {
    signInMode = !signInMode;
    document.getElementById('authTitle').textContent = signInMode ? 'Sign In' : 'Create Account';
    document.getElementById('toggleAuth').textContent = signInMode ? 'Create account' : 'Sign In';
}

function handleAuth(e) {
    e.preventDefault();
    const username = document.getElementById('usernameInput').value;
    const password = document.getElementById('passwordInput').value;
    const endpoint = signInMode ? 'signin' : 'signup';
    fetch(`/api/auth/${endpoint}`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({username, password})
    }).then(r => {
        if (!r.ok) throw new Error('failed');
        return r.json();
    }).then(d => {
        token = d.token;
        localStorage.setItem('token', token);
        localStorage.setItem('username', username);
        showApp();
    }).catch(() => alert('Auth failed'));
}

function checkStoredAuth() {
    token = localStorage.getItem('token');
    if (token) {
        showApp();
    }
}

function showApp() {
    document.getElementById('authSection').style.display = 'none';
    document.getElementById('app').style.display = 'block';
    document.getElementById('welcome').textContent = 'Hello ' + localStorage.getItem('username');
    loadGroups();
}

function signOut() {
    token = null;
    localStorage.clear();
    document.getElementById('app').style.display = 'none';
    document.getElementById('authSection').style.display = 'block';
}

function loadGroups() {
    fetch('/api/groups', {headers: {Authorization: 'Bearer ' + token}})
        .then(r => r.json())
        .then(groups => {
            const list = document.getElementById('groupList');
            list.innerHTML = '';
            groups.forEach(g => {
                const li = document.createElement('li');
                li.className = 'list-group-item list-group-item-action';
                li.textContent = g.name;
                li.onclick = () => openGroup(g.id, g.name);
                list.appendChild(li);
            });
        });
}

function newGroup() {
    const name = prompt('Group name?');
    if (!name) return;
    fetch('/api/groups', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', Authorization: 'Bearer ' + token },
        body: JSON.stringify({name})
    }).then(r => r.json()).then(loadGroups);
}

function openGroup(id, name) {
    currentGroupId = id;
    document.getElementById('currentGroup').textContent = name;
    document.getElementById('chatSection').style.display = 'block';
    document.getElementById('groupSection').style.display = 'none';
    connectWS();
    loadMessages();
}

function connectWS() {
    if (ws) ws.close();
    const proto = location.protocol === 'https:' ? 'wss' : 'ws';
    ws = new WebSocket(`${proto}://${location.host}/api/ws?token=${token}&channel=${currentGroupId}`);
    ws.onmessage = e => {
        const msg = JSON.parse(e.data);
        if (msg.type === 'message') addMessage(msg.username, msg.content);
    };
}

function sendMessage() {
    const input = document.getElementById('messageInput');
    const text = input.value.trim();
    if (!text || !ws) return;
    ws.send(text);
    input.value = '';
}

function addMessage(user, text) {
    const box = document.getElementById('messages');
    const div = document.createElement('div');
    div.textContent = user + ': ' + text;
    box.appendChild(div);
    box.scrollTop = box.scrollHeight;
}

function loadMessages(before) {
    let url = `/api/messages?channel=${currentGroupId}`;
    if (before) url += `&before=${before}`;
    fetch(url)
        .then(r => r.json())
        .then(msgs => {
            const box = document.getElementById('messages');
            box.innerHTML = '';
            msgs.forEach(m => addMessage(m.username, m.content));
        });
}

function inviteUser() {
    const user = prompt('Username to invite');
    if (!user) return;
    fetch(`/api/groups/${currentGroupId}/invite`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', Authorization: 'Bearer ' + token },
        body: JSON.stringify({username: user})
    }).then(r => {
        if (!r.ok) alert('Failed to invite');
    });
}
