const pwShowHide = document.querySelectorAll(".eye-icon");

pwShowHide.forEach(eyeIcon => {
    eyeIcon.addEventListener("click", () => {
        let pwField = eyeIcon.previousElementSibling;
        
        if (pwField.type === "password") {
            pwField.type = "text";
            eyeIcon.classList.replace("bx-hide", "bx-show");
        } else {
            pwField.type = "password";
            eyeIcon.classList.replace("bx-show", "bx-hide");
        }
    });
}); 

document
.getElementById("signupForm")
.addEventListener("submit", async function (event) {
  event.preventDefault();

  const formData = new FormData(event.target);
  const user = {
    first_name: formData.get("firstName"),
    last_name: formData.get("lastName"),
    email: formData.get("email"),
    password: formData.get("password"),
    phone: formData.get("phone"),
  };

  console.log("Sending JSON:", JSON.stringify(user));

  try {
    const response = await fetch("http://localhost:8000/users/signup", {
      method: "POST",
      mode: "no-cors",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(user),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || "Sign up failed");
    }

    const successMessage = await response.text();
    console.log(successMessage); // "Successfully Signed Up!!"

    // Redirect to Project1/landing.html after successful sign-up
    window.location.href = "login.html";
  } catch (error) {
    console.error(error.message); // Handle errors
  }
});

const loginForm = document.getElementById("loginForm");

loginForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const formData = new FormData(loginForm);
  const email = formData.get("email");
  const password = formData.get("password");

  try {
    const response = await axios.post("http://localhost:8000/users/login", {
      email: email,
      password: password,
    });

    // Assuming the backend returns the user data upon successful login
    const user = response.data;

    // Redirect to index.html on successful login
    window.location.href = "index.html";

  } catch (error) {
    // Handle login errors, display error message, etc.
    console.error("Login error:", error);
  }
});