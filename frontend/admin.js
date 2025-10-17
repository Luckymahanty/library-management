// Check authentication
const auth = checkAuth();
if (auth && auth.role !== 'admin') {
    window.location.href = 'user.html';
}

// Set admin greeting
const adminGreeting = document.getElementById('adminGreeting');
if (adminGreeting && auth) {
    adminGreeting.textContent = `Admin: ${auth.username}`;
}

// Setup logout
setupLogout();

// Get elements
const booksGrid = document.getElementById('booksGrid');
const addBookBtn = document.getElementById('addBookBtn');
const addBookModal = document.getElementById('addBookModal');
const editBookModal = document.getElementById('editBookModal');
const addBookForm = document.getElementById('addBookForm');
const editBookForm = document.getElementById('editBookForm');

let allBooks = [];

// Load books on page load
async function loadBooks() {
    allBooks = await fetchBooks();
    displayBooks(allBooks, booksGrid, true);
    setupSearch(allBooks, booksGrid, true);
}

// Modal controls
addBookBtn.addEventListener('click', () => {
    addBookModal.classList.add('show');
});

// Close modals
document.querySelectorAll('.close').forEach(closeBtn => {
    closeBtn.addEventListener('click', function() {
        this.closest('.modal').classList.remove('show');
    });
});

// Close modal on outside click
window.addEventListener('click', (e) => {
    if (e.target.classList.contains('modal')) {
        e.target.classList.remove('show');
    }
});

// Add book
addBookForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const bookData = {
        title: document.getElementById('bookTitle').value,
        author: document.getElementById('bookAuthor').value,
        genre: document.getElementById('bookGenre').value,
        quantity: parseInt(document.getElementById('bookQuantity').value)
    };
    
    try {
        const response = await fetch('http://localhost:8080/add', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(bookData)
        });
        
        const data = await response.json();
        
        if (response.ok) {
            showNotification('Book added successfully!', 'success');
            addBookModal.classList.remove('show');
            addBookForm.reset();
            
            // Reload books
            setTimeout(() => {
                loadBooks();
            }, 500);
        } else {
            showNotification(data.message || 'Failed to add book', 'error');
        }
    } catch (error) {
        console.error('Add book error:', error);
        showNotification('Unable to add book. Please try again.', 'error');
    }
});

// Edit book function
function editBook(id, title, author, genre, quantity) {
    document.getElementById('editBookId').value = id;
    document.getElementById('editBookTitle').value = title;
    document.getElementById('editBookAuthor').value = author;
    document.getElementById('editBookGenre').value = genre;
    document.getElementById('editBookQuantity').value = quantity;
    
    editBookModal.classList.add('show');
}

// Update book
editBookForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const bookData = {
        id: document.getElementById('editBookId').value,
        title: document.getElementById('editBookTitle').value,
        author: document.getElementById('editBookAuthor').value,
        genre: document.getElementById('editBookGenre').value,
        quantity: parseInt(document.getElementById('editBookQuantity').value)
    };
    
    try {
        const response = await fetch('http://localhost:8080/update', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(bookData)
        });
        
        const data = await response.json();
        
        if (response.ok) {
            showNotification('Book updated successfully!', 'success');
            editBookModal.classList.remove('show');
            
            // Reload books
            setTimeout(() => {
                loadBooks();
            }, 500);
        } else {
            showNotification(data.message || 'Failed to update book', 'error');
        }
    } catch (error) {
        console.error('Update book error:', error);
        showNotification('Unable to update book. Please try again.', 'error');
    }
});

// Delete book function
async function deleteBook(bookId) {
    if (!confirm('Are you sure you want to delete this book?')) {
        return;
    }
    
    try {
        const response = await fetch('http://localhost:8080/delete', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ id: bookId })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            showNotification('Book deleted successfully!', 'success');
            
            // Reload books
            setTimeout(() => {
                loadBooks();
            }, 500);
        } else {
            showNotification(data.message || 'Failed to delete book', 'error');
        }
    } catch (error) {
        console.error('Delete book error:', error);
        showNotification('Unable to delete book. Please try again.', 'error');
    }
}

// Initialize
loadBooks();
