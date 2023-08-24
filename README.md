```bash
# You can start the project with below commands
docker-compose up -d
go run main.go
```

- **SIGNUP FUNCTION API CALL (POST REQUEST)**

<http://localhost:8000/users/signup>

```json
{
  "first_name": "Sanchar",
  "last_name": "Limbu",
  "email": "test@gmail.com",
  "password": "12345678",
  "phone": "9841111111"
}
```

Response: "Successfully Signed Up!!"

- **LOGIN FUNCTION API CALL (POST REQUEST)**

  <http://localhost:8000/users/login>

```json
{
  "email": "test@gmail.com",
  "password": "12345678"
}
```

Response:

```json
{
    "_id": "64c283e8593e19b180b38f20",
    "first_name": "Sanchar",
    "last_name": "Limbu",
    "password": "$2a$14$XfMCoZ9Tdk4F8HeRDaP1FO4VVQWySHgnwQxvRjo.w/0AeQQ6FD5/y",
    "email": "test@gmail.com",
    "phone": "9841111111",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InRlc3RAZ21haWwuY29tIiwiRmlyc3RfTmFtZSI6IlNhbmNoYXIiLCJMYXN0X05hbWUiOiJMaW1idSIsIlVpZCI6IjY0YzI4M2U4NTkzZTE5YjE4MGIzOGYyMCIsImV4cCI6MTY5MDU1NTc1Nn0.TgS9RlTcOmMT-LyPFOqdFswu4fbqkWS3CngVNFmW6sk",
    "Refresh_Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6IiIsIkZpcnN0X05hbWUiOiIiLCJMYXN0X05hbWUiOiIiLCJVaWQiOiIiLCJleHAiOjE2OTEwNzQxNTZ9.JTJKeD1jNdHBmkeNy778aBqZB8IUzpjoVm0FC2CgSyI",
    "created_at": "2023-07-27T14:49:12Z",
    "updtaed_at": "2023-07-27T14:49:12Z",
    "user_id": "64c283e8593e19b180b38f20",
    "usercart": [],
    "address": [],
    "orders": []
}
```

- **Admin add Product Function (POST REQUEST)**

  <http://localhost:8000/admin/addproduct>

```json
{
  "product_name": "Alienware x15",
  "price": 2500,
  "rating": 10,
  "image": "alienware.jpg",
  "category": ["Laptop", "Gaming"]
}
```

Response : "Successfully added our Product Admin!!"

- **View all the Products in db GET REQUEST**

  <http://localhost:8000/users/productview>

Response

```json
[
    {
        "Product_ID": "64c215ae975cd1706d1492c5",
        "product_name": "Alienware x15",
        "price": 2500,
        "rating": 10,
        "image": "alienware.jpg",
        "category": [
            "Laptop"
        ]
    },
    {
        "Product_ID": "64c281df593e19b180b38f1f",
        "product_name": "Alienware x15",
        "price": 2500,
        "rating": 10,
        "image": "alienware.jpg",
        "category": [
            "Laptop",
            "Gaming"
        ]
    }
]
```

- **Search Product by regex function (GET REQUEST)**

<http://localhost:8000/users/search?name=Al>

Response:

```json
[
    {
        "Product_ID": "64c215ae975cd1706d1492c5",
        "product_name": "Alienware x15",
        "price": 2500,
        "rating": 10,
        "image": "alienware.jpg",
        "category": [
            "Laptop"
        ]
    },
    {
        "Product_ID": "64c281df593e19b180b38f1f",
        "product_name": "Alienware x15",
        "price": 2500,
        "rating": 10,
        "image": "alienware.jpg",
        "category": [
            "Laptop",
            "Gaming"
        ]
    }
]
```

- **Search Product by Category (GET REQUEST)**

  <http://localhost:8000/users/search-category?name=Gaming>

Response:

```json
[
    {
        "Product_ID": "64c281df593e19b180b38f1f",
        "product_name": "Alienware x15",
        "price": 2500,
        "rating": 10,
        "image": "alienware.jpg",
        "category": [
            "Laptop",
            "Gaming"
        ]
    }
]
```

- **Adding the Products to the Cart (GET REQUEST)**

  <http://localhost:8000/addtocart?id=64a3b75a227c961f850ef099&userID=64a3b752227c961f850ef098>

Response: "Successfully Added to the cart"

- **List the Products in the Cart (GET REQUEST)**

  <http://localhost:8000/listcart?id=64c283e8593e19b180b38f20>

Response:

```json
2500[
    {
        "Product_ID": "64c215ae975cd1706d1492c5",
        "product_name": "Alienware x15",
        "price": 2500,
        "rating": 10,
        "image": "alienware.jpg"
    }
]
```

- **Removing Item From the Cart (GET REQUEST)**

  <http://localhost:8000/removeitem?id=64c215ae975cd1706d1492c5&userID=64c283e8593e19b180b38f20>

Response: "Successfully removed from cart"

- **Addding the Address (POST REQUEST)**

  <http://localhost:8000/addaddress?id=64c283e8593e19b180b38f20>

  The Address array is limited to a single value more addresses is not acceptable

```json
{
  "house_name": "house",
  "street_name": "street",
  "city_name": "bhaktapur",
  "pin_code": "69420"
}
```

Response: "Successfully Added the Address"

- **Editing the Address(PUT REQUEST)**

  <http://localhost:8000/edithomeaddress?id=64c283e8593e19b180b38f20>

```json
{
    "house_name": "Lab",
    "street_name": "One",
    "city_name": "Crocas",
    "pin_code": "69420"
}
```

Response: "Successfully Updated the address"

- **Delete Addresses(GET REQUEST)**

  <http://localhost:8000/deleteaddresses?id=64c283e8593e19b180b38f20>

Response: "Successfully Deleted!"

- **Cart Checkout Function and placing the order(GET REQUEST)**

  <http://localhost:8000/cartcheckout?id=64c283e8593e19b180b38f20>

Response: "Successfully Placed the order"

- **Instantly Buying the Products(GET REQUEST)**
  <http://localhost:8000/instantbuy?userid=64c283e8593e19b180b38f20&pid=64c215ae975cd1706d1492c5>

Response: "Successully placed the order"
