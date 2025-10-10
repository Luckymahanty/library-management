async function loadBooks() {
  const res = await fetch("/books");
  const books = await res.json();

  const list = document.getElementById("bookList");
  list.innerHTML = "";
  books.forEach(b => {
    const li = document.createElement("li");
    li.textContent = `${b.title} by ${b.author}`;
    const btn = document.createElement("button");
    btn.textContent = "Delete";
    btn.onclick = async () => {
      await fetch(`/deletebook?id=${b.id}`);
      loadBooks();
    };
    li.appendChild(btn);
    list.appendChild(li);
  });
}

document.getElementById("addBookForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  const title = document.getElementById("title").value;
  const author = document.getElementById("author").value;

  await fetch("/addbook", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ title, author }),
  });

  document.getElementById("addBookForm").reset();
  loadBooks();
});

loadBooks();
