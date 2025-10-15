document.getElementById('loginForm').addEventListener('submit', async (e)=> {
  e.preventDefault();
  const username = document.getElementById('username').value.trim();
  const password = document.getElementById('password').value.trim();
  const res = await fetch('/login', {
    method:'POST',
    headers:{'Content-Type':'application/json'},
    body: JSON.stringify({username,password})
  });
  if (!res.ok) {
    alert('Login failed');
    return;
  }
  const user = await res.json();
  if (user.role === 'admin') location.href = 'admin.html';
  else location.href = 'user.html';
});
