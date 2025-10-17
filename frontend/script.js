// Shared utility functions

// Check if user is logged in
function checkAuth() {
    const username = sessionStorage.getItem('username');
    const role = sessionStorage.getItem('role');
    
    if (!username || !role) {
        window.location.href = 'login.html';
        return null;
    }
    
    return { username, role };
}

// Fetch all books from backend
async function fetchBooks() {
    try {
        const response = await fetch('http://localhost:8080/books');
        
        if (!response.ok) {
            throw new Error('Failed to fetch books');
        }
        
        const books = await response.json();
        return books;
    } catch (error) {
        console.error('Error fetching books:', error);
        return [];
    }
}

// Display books in grid
function displayBooks(books, container, isAdmin = false) {
    if (!books || books.length === 0) {
        container.innerHTML = `
            <div class="loading">
                <i class="fas fa-book-open"></i>
                <p>No books available at the moment.</p>
            </div>
        `;
        return;
    }
    
    container.innerHTML = '';
    
    books.forEach((book, index) => {
        const bookCard = createBookCard(book, isAdmin, index);
        container.appendChild(bookCard);
    });
}

// Create book card element
function createBookCard(book, isAdmin, index) {
    const card = document.createElement('div');
    card.className = 'book-card';
    card.style.animationDelay = `${index * 0.1}s`;
    
    const available = book.quantity > 0;
    
    card.innerHTML = `
        <div class="book-icon">
            <i class="fas fa-book"></i>
        </div>
        <h3>${book.title}</h3>
        <div class="book-info">
            <div class="book-info-item">
                <i class="fas fa-user"></i>
                <span>${book.author}</span>
            </div>
            <div class="book-info-item">
                <i class="fas fa-tag"></i>
                <span>${book.genre}</span>
            </div>
            <div class="book-info-item">
                <i class="fas fa-boxes"></i>
                <span>${book.quantity} available</span>
            </div>
        </div>
        <div class="book-actions">
            ${isAdmin ? `
                <button class="btn btn-primary" onclick="editBook('${book.id}', '${book.title}', '${book.author}', '${book.genre}', ${book.quantity})">
                    <i class="fas fa-edit"></i> Edit
                </button>
                <button class="btn btn-danger" onclick="deleteBook('${book.id}')">
                    <i class="fas fa-trash"></i> Delete
                </button>
            ` : `
                <button class="btn ${available ? 'btn-success' : 'btn-secondary'}" 
                        onclick="borrowBook('${book.id}')" 
                        ${!available ? 'disabled' : ''}>
                    <i class="fas ${available ? 'fa-hand-holding' : 'fa-ban'}"></i> 
                    ${available ? 'Borrow' : 'Unavailable'}
                </button>
            `}
        </div>
    `;
    
    return card;
}

// Search functionality
function setupSearch(books, container, isAdmin) {
    const searchInput = document.getElementById('searchInput');
    
    if (!searchInput) return;
    
    searchInput.addEventListener('input', (e) => {
        const searchTerm = e.target.value.toLowerCase();
        
        const filteredBooks = books.filter(book => 
            book.title.toLowerCase().includes(searchTerm) ||
            book.author.toLowerCase().includes(searchTerm) ||
            book.genre.toLowerCase().includes(searchTerm)
        );
        
        displayBooks(filteredBooks, container, isAdmin);
    });
}

// Logout functionality
function setupLogout() {
    const logoutBtn = document.getElementById('logoutBtn');
    
    if (logoutBtn) {
        logoutBtn.addEventListener('click', () => {
            sessionStorage.clear();
            window.location.href = 'index.html';
        });
    }
}

// Show notification
function showNotification(message, type = 'success') {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.style.cssText = `
        position: fixed;
        top: 80px;
        right: 20px;
        background: ${type === 'success' ? '#10b981' : '#ef4444'};
        color: white;
        padding: 1rem 1.5rem;
        border-radius: 10px;
        box-shadow: 0 5px 20px rgba(0,0,0,0.2);
        z-index: 3000;
        animation: slideIn 0.3s ease;
    `;
    notification.textContent = message;
    
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease';
        setTimeout(() => notification.remove(), 300);
    }, 3000);
}

// Add animation styles
const style = document.createElement('style');
style.textContent = `
    @keyframes slideIn {
        from {
            transform: translateX(400px);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }
    
    @keyframes slideOut {
        from {
            transform: translateX(0);
            opacity: 1;
        }
        to {
            transform: translateX(400px);
            opacity: 0;
        }
    }
`;
document.head.appendChild(style);
