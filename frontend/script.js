// Global state
let books = [];
let currentUser = null;

// Page Navigation Functions
function showWelcome() {
  document.getElementById('welcome-page').style.display = 'flex';
  document.getElementById('login-page').style.display = 'none';
  document.getElementById('signup-page').style.display = 'none';
  document.getElementById('home-page').style.display = 'none';
}

function showLogin() {
  document.getElementById('welcome-page').style.display = 'none';
  document.getElementById('login-page').style.display = 'flex';
  document.getElementById('signup-page').style.display = 'none';
  document.getElementById('home-page').style.display = 'none';
}

function showSignup() {
  document.getElementById('welcome-page').style.display = 'none';
  document.getElementById('login-page').style.display = 'none';
  document.getElementById('signup-page').style.display = 'flex';
  document.getElementById('home-page').style.display = 'none';
}

function showHome() {
  document.getElementById('welcome-page').style.display = 'none';
  document.getElementById('login-page').style.display = 'none';
  document.getElementById('signup-page').style.display = 'none';
  document.getElementById('home-page').style.display = 'block';
}

// Signup form submission
async function handleSignup(e) {
  e.preventDefault();
  const username = document.getElementById('signup-username').value.trim();
  const password = document.getElementById('signup-password').value;
  const confirmPassword = document.getElementById('signup-confirm-password').value;

  // Validate passwords match
  if (password !== confirmPassword) {
    showSignupError('Passwords do not match!');
    return;
  }

  if (!username || !password) {
    showSignupError('Please fill in all fields');
    return;
  }

  try {
  const res = await fetch('/signup', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ 
      username: username, 
      password: password,
      role: 'user'
    }),
  });
  const text = await res.text();
  
  if (text.includes('successful')) {
    alert('Signup successful! Please login.');
    document.getElementById('signup-form').reset();
    showLogin();
  } else {
    showSignupError(text);
  }
} catch (error) {
  console.error('Signup error:', error);
  showSignupError('Signup failed. Please try again.');
}

// Login form submission
async function handleLogin(e) {
  e.preventDefault();
  const username = document.getElementById('username').value.trim();
  const password = document.getElementById('password').value;

  if (!username || !password) {
    showError('Please enter both username and password');
    return;
  }

  try {
    const res = await fetch('/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: `username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`,
    });
    const text = await res.text();

    if (text.includes('successful')) {
      currentUser = {
        username: username,
        role: 'Administrator'
      };
      // Update UI with user info
      document.getElementById('user-name').textContent = username;
      document.getElementById('user-role').textContent = 'Administrator';
      // Clear form and show home page
      document.getElementById('login-form').reset();
      showHome();
      // Load books from backend
      loadBooks();
    } else {
      showError(text);
    }
  } catch (error) {
    console.error('Login error:', error);
    showError('Login failed. Please try again.');
  }
}

function handleLogout() {
  currentUser = null;
  books = [];
  renderBooks();
  updateStats();
  showWelcome();
}

function showError(message) {
  const errorDiv = document.getElementById('login-error');
  errorDiv.textContent = message;
  errorDiv.style.display = 'block';
  setTimeout(() => {
    errorDiv.style.display = 'none';
  }, 3000);
}

function showSignupError(message) {
  const errorDiv = document.getElementById('signup-error');
  errorDiv.textContent = message;
  errorDiv.style.display = 'block';
  setTimeout(() => {
    errorDiv.style.display = 'none';
  }, 3000);
}

// Book Management Functions
async function loadBooks() {
  try {
    const response = await fetch('/api/books');
    if (response.ok) {
      books = await response.json();
      renderBooks();
      updateStats();
    }
  } catch (error) {
    console.error('Error loading books:', error);
    // Use empty array if API fails
    books = [];
    renderBooks();
    updateStats();
  }
}

async function handleAddBook(e) {
  e.preventDefault();
  const title = document.getElementById('book-title').value.trim();
  const author = document.getElementById('book-author').value.trim();

  if (!title || !author) {
    alert('Please fill in all fields');
    return;
  }

  try {
    const response = await fetch('/api/books', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title, author })
    });

    if (response.ok) {
      const newBook = await response.json();
      books.push(newBook);
      document.getElementById('add-book-form').reset();
      renderBooks();
      updateStats();
      alert('Book added successfully!');
    } else {
      const error = await response.text();
      alert('Failed to add book: ' + error);
    }
  } catch (error) {
    console.error('Error adding book:', error);
    // Fallback: Add book locally if API fails
    const newBook = {
      id: Date.now(),
      title: title,
      author: author,
      addedDate: new Date().toISOString()
    };
    books.push(newBook);
    document.getElementById('add-book-form').reset();
    renderBooks();
    updateStats();
    alert('Book added locally (API unavailable)');
  }
}

async function deleteBook(id) {
  if (!confirm('Are you sure you want to delete this book?')) {
    return;
  }

  try {
    const response = await fetch(`/api/books/${id}`, {
      method: 'DELETE'
    });

    if (response.ok) {
      books = books.filter(book => book.id !== id);
      renderBooks();
      updateStats();
    } else {
      alert('Failed to delete book');
    }
  } catch (error) {
    console.error('Error deleting book:', error);
    // Fallback: Delete locally if API fails
    books = books.filter(book => book.id !== id);
    renderBooks();
    updateStats();
  }
}

function editBook(id) {
  alert('Edit feature coming soon! This will connect to your backend API.');
}

// Render Functions
function renderBooks(filteredBooks = books) {
  const bookList = document.getElementById('book-list');
  
  if (filteredBooks.length === 0) {
    bookList.innerHTML = '<div class="no-books">No books found. Add one above!</div>';
    return;
  }

  bookList.innerHTML = filteredBooks.map(book => `
    <div class="book-card" data-book-id="${book.id}">
      <div class="book-title">${escapeHtml(book.title)}</div>
      <div class="book-author">by ${escapeHtml(book.author)}</div>
      <div class="book-actions">
        <button class="btn-small btn-edit" onclick="editBook(${book.id})">Edit</button>
        <button class="btn-small btn-delete" onclick="deleteBook(${book.id})">Delete</button>
      </div>
    </div>
  `).join('');
}

function handleSearch() {
  const searchTerm = document.getElementById('search-books').value.toLowerCase().trim();
  
  if (!searchTerm) {
    renderBooks();
    return;
  }

  const filtered = books.filter(book =>
    book.title.toLowerCase().includes(searchTerm) ||
    book.author.toLowerCase().includes(searchTerm)
  );
  renderBooks(filtered);
}

function updateStats() {
  document.getElementById('total-books').textContent = books.length;
}

// Utility Functions
function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}

// Event Listeners
document.addEventListener('DOMContentLoaded', function() {
  // Initialize the application
  showWelcome();

  // Setup login form event listener
  const loginForm = document.getElementById('login-form');
  if (loginForm) {
    loginForm.addEventListener('submit', handleLogin);
  }

  // Setup signup form event listener
  const signupForm = document.getElementById('signup-form');
  if (signupForm) {
    signupForm.addEventListener('submit', handleSignup);
  }

  // Setup add book form event listener
  const addBookForm = document.getElementById('add-book-form');
  if (addBookForm) {
    addBookForm.addEventListener('submit', handleAddBook);
  }

  // Setup search with debouncing
  const searchInput = document.getElementById('search-books');
  if (searchInput) {
    let searchTimeout;
    searchInput.addEventListener('input', function() {
      clearTimeout(searchTimeout);
      searchTimeout = setTimeout(handleSearch, 300);
    });
  }
});

// Export functions for use in HTML onclick attributes
window.showWelcome = showWelcome;
window.showLogin = showLogin;
window.showSignup = showSignup;
window.showHome = showHome;
window.handleLogin = handleLogin;
window.handleSignup = handleSignup;
window.handleLogout = handleLogout;
window.handleAddBook = handleAddBook;
window.deleteBook = deleteBook;
window.editBook = editBook;
window.handleSearch = handleSearch;
