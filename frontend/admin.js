async function fetchBooks() {
  const res = await fetch("/books");
  const books = await res.json();

  const tableBody = document.querySelector("#booksTable tbody");
  tableBody.innerHTML = "";

  books.forEach(book => {
    const row = `
      <tr>
        <td>${book.id}</td>
        <td>${book.title}</td>
        <td>${book.author}</td>
        <td><button onclick="deleteBook(${book.id})">Delete</button></td>
      </tr>
    `;
    tableBody.innerHTML += row;
  });
}

document.getElementById("addBookForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  const title = document.getElementById("title").value;
  const author = document.getElementById("author").value;

  const res = await fetch("/addbook", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ title, author }),
  });

  if (res.ok) {
    alert("Book added!");
    fetchBooks();
  } else {
    alert("Error adding book");
  }
});

async function deleteBook(id) {
  const res = await fetch(`/deletebook?id=${id}`, { method: "DELETE" });
  if (res.ok) {
    alert("Book deleted!");
    fetchBooks();
  } else {
    alert("Failed to delete book");
  }
}

function logout() {
  window.location.href = "index.html";
}

fetchBooks();
