// Check authentication
const auth = checkAuth();
if (auth && auth.role !== 'user') {
    window.location.href = 'admin.html';
}

// Set user greeting
const userGreeting = document.getElementById('userGreeting');
if (userGreeting && auth) {
    userGreeting.textContent = `Welcome, ${auth.username}`;
}

// Setup logout
setupLogout();

// Get books container
const booksGrid = document.getElementById('booksGrid');
let allBooks = [];

// Load books on page load
async function loadBooks() {
    allBooks = await fetchBooks();
    displayBooks(allBooks, booksGrid, false);
    setupSearch(allBooks, booksGrid, false);
}

// Borrow book function
async function borrowBook(bookId) {
    const username = sessionStorage.getItem('username');
    
    if (!username) {
        showNotification('Please login to borrow books', 'error');
        return;
    }
    
    try {
        const response = await fetch('http://localhost:8080/borrow', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                bookId: bookId,
                username: username 
            })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            showNotification('Book borrowed successfully!', 'success');
            
            // Reload books to update availability
            setTimeout(() => {
                loadBooks();
            }, 1000);
        } else {
            showNotification(data.message || 'Failed to borrow book', 'error');
        }
    } catch (error) {
        console.error('Borrow error:', error);
        showNotification('Unable to process request. Please try again.', 'error');
    }
}

// Initialize
loadBooks();
