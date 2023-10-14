// Check if the user is logged in
var token = localStorage.getItem("token");
var username = localStorage.getItem("username");

// Get the userCart div
var userCartDiv = document.getElementById("userCart");

if (token && username) {
  // User is logged in, display the userCart
  userCartDiv.style.display = "block";
} else {
  // User is not logged in, hide the userCart
  userCartDiv.style.display = "none";
}
