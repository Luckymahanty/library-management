async function fetchBooks() {
  try {
    const res = await fetch('/books');
    if (!res.ok) return;
    const books = await res.json();
    const list = document.getElementById('book-list');
    if (!list) return;
    list.innerHTML = '';
    books.forEach(b => {
      const li = document.createElement('li');
      li.className = 'book';
      li.innerHTML = `<strong>${escapeHtml(b.title)}</strong> — ${escapeHtml(b.author)} <span style="color:gray">(${b.status||'available'})</span>`;
      // if admin page, admin.js will add delete buttons by checking #addBookBtn
      list.appendChild(li);
    });
  } catch (e) {
    console.error('fetchBooks', e);
  }
}

function escapeHtml(s) {
  const d = document.createElement('div'); d.textContent = s; return d.innerHTML;
}

document.addEventListener('DOMContentLoaded', fetchBooks);
