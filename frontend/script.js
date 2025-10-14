function fetchBooks() {
  fetch("/books")
    .then(res => res.json())
    .then(data => {
      const list = document.getElementById("book-list");
      list.innerHTML = "";
      data.forEach(b => {
        list.innerHTML += `
          <li>${b.title} by ${b.author}
            <button onclick="deleteBook(${b.id})">Delete</button>
          </li>`;
      });
    });
}

function addBook() {
  const title = document.getElementById("bookTitle").value;
  const author = document.getElementById("bookAuthor").value;

  fetch("/addbook", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ title, author }),
  })
  .then(res => res.json())
  .then(() => fetchBooks());
}

function deleteBook(id) {
  fetch(`/deletebook?id=${id}`, { method: "DELETE" })
    .then(res => res.json())
    .then(() => fetchBooks());
}
