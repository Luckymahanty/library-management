async function loadBooks() {
  const res = await fetch("/books");
  const books = await res.json();

  const list = document.getElementById("userBookList");
  list.innerHTML = "";
  books.forEach(b => {
    const li = document.createElement("li");
    li.textContent = `${b.title} by ${b.author}`;
    list.appendChild(li);
  });
}

loadBooks();
