const API_URL = "http://localhost:8080";

// Fetch and display books
async function getBooks() {
  const res = await fetch(`${API_URL}/books`);
  const books = await res.json();

  const list = document.getElementById("book-list");
  list.innerHTML = "";
  books.forEach(book => {
    let li = document.createElement("li");
    li.textContent = `${book.id}. ${book.title} by ${book.author} [${book.status}]`;
    list.appendChild(li);
  });
}

// Add a new book
async function addBook() {
  const title = document.getElementById("title").value;
  const author = document.getElementById("author").value;

  if (!title || !author) {
    alert("Please enter book title and author");
    return;
  }

  await fetch(`${API_URL}/add`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ title, author })
  });

  document.getElementById("title").value = "";
  document.getElementById("author").value = "";
  getBooks(); // refresh list
}

// Load books on page load
getBooks();
