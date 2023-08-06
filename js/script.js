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
      first_name: encodeURIComponent(formData.get("firstName")),
      last_name: encodeURIComponent(formData.get("lastName")),
      email: encodeURIComponent(formData.get("email")),
      password: encodeURIComponent(formData.get("password")),
      phone: encodeURIComponent(formData.get("phone")),
    };
    

    // console.log("Sending JSON:", JSON.stringify(user));

    console.log(JSON.stringify(user));

    try {
      const response = await fetch("http://localhost:8000/users/signup", {
        method: "POST",
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

      // Redirect to ......... after successful sign-up
      window.location.href = "login.html";
    } catch (error) {
      console.error('JSON error:', error.message);

      // Display the error message in the frontend
      const errorDiv = document.getElementById("errorDiv");
      errorDiv.textContent = error.message;
    }
  });


