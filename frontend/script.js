// Global state
let books = [];
let currentUser = null;

// Page Navigation Functions
function showWelcome() {
  document.getElementById('welcome-page').style.display = 'flex';
  document.getElementById('login-page').style.display = 'none';
  document.getElementById('home-page').style.display = 'none';
}

function showLogin() {
  document.getElementById('welcome-page').style.display = 'none';
  document.getElementById('login-page').style.display = 'flex';
  document.getElementById('home-page').style.display = 'none';
}

function showHome() {
  document.getElementById('welcome-page').style.display = 'none';
  document.getElementById('login-page').style.display = 'none';
  document.getElementById('home-page').style.display = 'block';
}

// Authentication Functions
function handleLogin(e) {
  e.preventDefault();
  
  const username = document.getElementById('username').value.trim();
  const password = document.getElementById('password').value;

  // Basic validation (replace with actual API call to backend)
  if (!username || !password) {
    showError('Please enter both username and password');
    return;
  }

  // TODO: Replace with actual authentication API call
  // Example: fetch('/api/auth/login', { method: 'POST', body: JSON.stringify({ username, password }) })
  
  if (username && password) {
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
    showError('Invalid credentials');
  }
}

function handleLogout() {
  // TODO: Call backend logout API
  // Example: fetch('/api/auth/logout', { method: 'POST' })
  
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

// Book Management Functions
async function loadBooks() {
  // TODO: Replace with actual API call to fetch books from backend
  // Example:
  // try {
  //   const response = await fetch('/api/books');
  //   books = await response.json();
  //   renderBooks();
  //   updateStats();
  // } catch (error) {
  //   console.error('Error loading books:', error);
  // }
  
  // For now, use mock data or empty array
  books = [];
  renderBooks();
  updateStats();
}

function handleAddBook(e) {
  e.preventDefault();
  
  const title = document.getElementById('book-title').value.trim();
  const author = document.getElementById('book-author').value.trim();

  if (!title || !author) {
    alert('Please fill in all fields');
    return;
  }

  // TODO: Replace with actual API call to add book to backend
  // Example:
  // try {
  //   const response = await fetch('/api/books', {
  //     method: 'POST',
  //     headers: { 'Content-Type': 'application/json' },
  //     body: JSON.stringify({ title, author })
  //   });
  //   const newBook = await response.json();
  //   books.push(newBook);
  // } catch (error) {
  //   console.error('Error adding book:', error);
  // }

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
  
  // Show success feedback
  alert('Book added successfully!');
}

function deleteBook(id) {
  if (!confirm('Are you sure you want to delete this book?')) {
    return;
  }

  // TODO: Replace with actual API call to delete book from backend
  // Example:
  // try {
  //   await fetch(`/api/books/${id}`, { method: 'DELETE' });
  //   books = books.filter(book => book.id !== id);
  // } catch (error) {
  //   console.error('Error deleting book:', error);
  // }

  books = books.filter(book => book.id !== id);
  renderBooks();
  updateStats();
}

function editBook(id) {
  // TODO: Implement edit functionality
  // This could open a modal or navigate to an edit page
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
  
  // TODO: Update other stats from backend API
  // Example:
  // fetch('/api/stats').then(res => res.json()).then(stats => {
  //   document.getElementById('active-members').textContent = stats.activeMembers;
  //   document.getElementById('books-issued').textContent = stats.booksIssued;
  //   document.getElementById('due-today').textContent = stats.dueToday;
  // });
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
  
  // Setup form event listeners
  const loginForm = document.getElementById('login-form');
  if (loginForm) {
    loginForm.addEventListener('submit', handleLogin);
  }
  
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
window.showHome = showHome;
window.handleLogin = handleLogin;
window.handleLogout = handleLogout;
window.handleAddBook = handleAddBook;
window.deleteBook = deleteBook;
window.editBook = editBook;
window.handleSearch = handleSearch;
