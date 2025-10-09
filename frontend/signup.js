document.getElementById("signupForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value.trim();
  const role = document.getElementById("role").value;

  try {
    const res = await fetch("/signup", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password, role }),
    });

    const text = await res.text();
    alert(text);

    if (res.ok) {
      window.location.href = "login.html";
    }
  } catch (err) {
    alert("⚠️ Error connecting to server: " + err.message);
  }
});
