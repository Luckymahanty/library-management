document.getElementById("signupForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value.trim();
  const role = document.getElementById("role").value;

  const res = await fetch("/signup", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password, role }),
  });

  if (res.ok) {
    alert("✅ Signup successful! Please login now.");
    window.location.href = "login.html";
  } else {
    const msg = await res.text();
    alert("❌ " + msg);
  }
});
