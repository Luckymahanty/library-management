// Global state
let books = [];
let currentUser = null;

// Page Navigation
function showWelcome() {
  document.getElementById('welcome-page').style.display = 'flex';
  document.getElementById('login-page').style.display = 'none';
  document.getElementById('signup-page').style.display = 'none';
}

function showLogin() {
  document.getElementById('welcome-page').style.display = 'none';
  document.getElementById('login-page').style.display = 'flex';
  document.getElementById('signup-page').style.display = 'none';
}

function showSignup() {
  document.getElementById('welcome-page').style.display = 'none';
  document.getElementById('login-page').style.display = 'none';
  document.getElementById('signup-page').style.display = 'flex';
}

// ====================== SIGNUP =========================
async function handleSignup(e) {
  e.preventDefault();

  const username = document.getElementById('signup-username').value.trim();
  const password = document.getElementById('signup-password').value;
  const confirmPassword = document.getElementById('signup-confirm-password').value;
  const role = document.getElementById('signup-role') ? document.getElementById('signup-role').value : 'user';

  if (!username || !password) {
    showSignupError('Please fill all fields');
    return;
  }

  if (password !== confirmPassword) {
    showSignupError('Passwords do not match!');
    return;
  }

  try {
    const res = await fetch('/signup', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: `username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}&role=${encodeURIComponent(role)}`
    });

    const text = await res.text();

    if (text.toLowerCase().includes('successful')) {
      alert('Signup successful! Please login.');
      document.getElementById('signup-form').reset();
      showLogin();
    } else {
      showSignupError(text);
    }
  } catch (err) {
    console.error('Signup error:', err);
    showSignupError('Signup failed, please try again.');
  }
}

// ====================== LOGIN =========================
async function handleLogin(e) {
  e.preventDefault();

  const username = document.getElementById('username').value.trim();
  const password = document.getElementById('password').value;

  if (!username || !password) {
    showError('Please enter username and password');
    return;
  }

  try {
    const res = await fetch('/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: `username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`
    });

    const text = await res.text();

    if (text.toLowerCase().includes('successful')) {
      alert('Login successful!');

      let role = 'user';
      if (text.toLowerCase().includes('admin')) {
        role = 'admin';
      }

      // redirect to proper dashboard
      if (role === 'admin') {
        window.location.href = 'admin.html';
      } else {
        window.location.href = 'user.html';
      }

    } else {
      showError(text);
    }
  } catch (err) {
    console.error('Login error:', err);
    showError('Login failed, please try again.');
  }
}

// ====================== LOGOUT =========================
function handleLogout() {
  alert('Logged out successfully!');
  window.location.href = 'index.html';
}

// ====================== ERROR DISPLAY =========================
function showError(msg) {
  const err = document.getElementById('login-error');
  if (err) {
    err.textContent = msg;
    err.style.display = 'block';
    setTimeout(() => (err.style.display = 'none'), 3000);
  } else {
    alert(msg);
  }
}

function showSignupError(msg) {
  const err = document.getElementById('signup-error');
  if (err) {
    err.textContent = msg;
    err.style.display = 'block';
    setTimeout(() => (err.style.display = 'none'), 3000);
  } else {
    alert(msg);
  }
}

// ====================== BOOK MANAGEMENT (for dashboard) =========================
async function loadBooks() {
  try {
    const res = await fetch('/api/books');
    if (res.ok) {
      books = await res.json();
      renderBooks();
    }
  } catch (err) {
    console.error('Error loading books:', err);
  }
}

async function handleAddBook(e) {
  e.preventDefault();
  const title = document.getElementById('book-title').value.trim();
  const author = document.getElementById('book-author').value.trim();

  if (!title || !author) {
    alert('Please fill all fields');
    return;
  }

  try {
    const res = await fetch('/api/books', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title, author })
    });

    if (res.ok) {
      const newBook = await res.json();
      books.push(newBook);
      renderBooks();
      alert('Book added successfully!');
    }
  } catch (err) {
    console.error('Error adding book:', err);
  }
}

async function deleteBook(id) {
  if (!confirm('Are you sure you want to delete this book?')) return;

  try {
    const res = await fetch(`/api/books/${id}`, { method: 'DELETE' });
    if (res.ok) {
      books = books.filter(b => b.id !== id);
      renderBooks();
    }
  } catch (err) {
    console.error('Error deleting book:', err);
  }
}

function renderBooks() {
  const list = document.getElementById('book-list');
  if (!list) return;

  if (books.length === 0) {
    list.innerHTML = '<p>No books available.</p>';
    return;
  }

  list.innerHTML = books.map(
    b => `
      <div class="book-card">
        <h3>${b.title}</h3>
        <p>by ${b.author}</p>
        <button onclick="deleteBook(${b.id})">Delete</button>
      </div>
    `
  ).join('');
}

// ====================== INITIALIZATION =========================
document.addEventListener('DOMContentLoaded', () => {
  const loginForm = document.getElementById('login-form');
  if (loginForm) loginForm.addEventListener('submit', handleLogin);

  const signupForm = document.getElementById('signup-form');
  if (signupForm) signupForm.addEventListener('submit', handleSignup);

  const addBookForm = document.getElementById('add-book-form');
  if (addBookForm) addBookForm.addEventListener('submit', handleAddBook);
});

// Export for HTML buttons
window.showWelcome = showWelcome;
window.showLogin = showLogin;
window.showSignup = showSignup;
window.handleLogin = handleLogin;
window.handleSignup = handleSignup;
window.handleLogout = handleLogout;
window.handleAddBook = handleAddBook;
window.deleteBook = deleteBook;

