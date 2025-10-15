document.getElementById('signupForm').addEventListener('submit', async (e)=> {
  e.preventDefault();
  const username = document.getElementById('su-username').value.trim();
  const password = document.getElementById('su-password').value.trim();
  const role = document.getElementById('su-role').value;
  const res = await fetch('/signup', {
    method:'POST',
    headers:{'Content-Type':'application/json'},
    body: JSON.stringify({username,password,role})
  });
  if (res.ok) {
    alert('Signup successful. Please login.');
    location.href='login.html';
  } else {
    alert('Signup failed');
  }
});
