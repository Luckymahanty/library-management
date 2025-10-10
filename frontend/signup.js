document.getElementById("signupForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value.trim();
  const role = document.getElementById("role").value;

  if (!username || !password) {
    alert("All fields are required!");
    return;
  }

  const res = await fetch("/signup", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password, role }),
  });

  if (res.ok) {
    alert("✅ Signup successful! You can now login.");
    window.location.href = "login.html";
  } else {
    const err = await res.text();
    alert("❌ " + err);
  }
});
