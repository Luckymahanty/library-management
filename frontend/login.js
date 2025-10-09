document.getElementById("loginForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value.trim();

  try {
    const res = await fetch("/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });

    const text = await res.text();
    alert(text);

    if (res.ok) {
      window.location.href = "index.html";
    }
  } catch (err) {
    alert("⚠️ Error connecting to server: " + err.message);
  }
});
