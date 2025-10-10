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
      </tr>
    `;
    tableBody.innerHTML += row;
  });
}

function logout() {
  window.location.href = "index.html";
}

fetchBooks();
