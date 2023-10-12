// Function to display cart items
function displayCart() {
    const cartContainer = document.getElementById("cart-container");
    cartContainer.innerHTML = "";
  
    if (cart.length === 0) {
      // The cart is empty; display a message or hide the cart content
      cartContainer.innerHTML = "<p>Your cart is empty.</p>";
    } else {
      // The cart has items; dynamically create cart items based on the real product data
      let total = 0;
  
      cart.forEach((item) => {
        const cartItem = document.createElement("div");
        cartItem.className = "cart-item";
  
        const cartItemImage = document.createElement("img");
        cartItemImage.src = decodeURIComponent(item.image);
  
        const cartItemInfo = document.createElement("div");
        cartItemInfo.className = "cart-item-info";
  
        const itemName = document.createElement("h5");
        itemName.textContent = item.name;
  
        const itemPrice = document.createElement("p");
        itemPrice.textContent = item.price;
  
        cartItemInfo.appendChild(itemName);
        cartItemInfo.appendChild(itemPrice);
  
        cartItem.appendChild(cartItemImage);
        cartItem.appendChild(cartItemInfo);
        cartContainer.appendChild(cartItem);
  
        total += parseFloat(item.price.replace("$", ""));
      });
  
      const totalPriceElement = document.getElementById("total-price");
      totalPriceElement.textContent = "Total Price: $" + total.toFixed(2);
    }
  }
  