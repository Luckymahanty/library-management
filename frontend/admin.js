document.addEventListener('DOMContentLoaded', () => {
  const titleInput = document.getElementById('title');
  const authorInput = document.getElementById('author');
  const addBookBtn = document.getElementById('addBookBtn');
  const list = document.getElementById('book-list');

  // Fetch all books and show with delete buttons
  async function loadBooks() {
    try {
      const res = await fetch('/books');
      const books = await res.json();
      list.innerHTML = '';

      books.forEach(b => {
        const li = document.createElement('li');
        li.innerHTML = `
          <strong>${b.title}</strong> — ${b.author} 
          <button class="delete" data-id="${b.id}">Delete</button>
        `;
        list.appendChild(li);
      });
    } catch (err) {
      console.error('Failed to load books', err);
    }
  }

  // Add new book
  addBookBtn.addEventListener('click', async () => {
    const title = titleInput.value.trim();
    const author = authorInput.value.trim();
    if (!title || !author) return alert('Enter title and author');

    await fetch('/books', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title, author })
    });

    titleInput.value = '';
    authorInput.value = '';
    await loadBooks(); // refresh list
  });

  // Handle delete clicks
  list.addEventListener('click', async (e) => {
    if (e.target.matches('.delete')) {
      const id = e.target.dataset.id;
      const ok = confirm('Are you sure you want to delete this book?');
      if (ok) {
        await fetch('/books/' + id, { method: 'DELETE' });
        e.target.closest('li').remove();
      }
    }
  });

  // Initial load
  loadBooks();
});
