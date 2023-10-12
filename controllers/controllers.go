package controllers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"go-api/database"
	"go-api/models"
	generate "go-api/tokens"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	emailRegex                          = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	passwordRegex                       = regexp.MustCompile(`^[a-zA-Z0-9]*[A-Z]+[a-zA-Z0-9]*\d+[a-zA-Z0-9]*$`)
	UserCollection    *mongo.Collection = database.UserData(database.Client, "Users")
	ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userpassword string, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Password Incorrect"
		valid = false
	}

	return valid, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if *user.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email cannot be empty"})
			return
		}

		// Validate email format
		if !emailRegex.MatchString(*user.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			return
		}

		// Validate password format (at least 8 characters, 1 capital letter, and 1 number)
		if !passwordRegex.MatchString(*user.Password) {
			c.JSON(http.StatusBadRequest,
				gin.H{"error": "Password must be at least 8 characters long and contain at least 1 capital letter and 1 number"})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone is already in use"})
			return
		}

		passwordHash := HashPassword(*user.Password)
		user.Password = &passwordHash

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}

		defer cancel()

		c.JSON(http.StatusCreated, "Successfully Signed Up!!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var founduser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Login Incorrect"})
			return
		}

		PasswordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)
		defer cancel()

		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID)
		defer cancel()

		generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)
		c.JSON(http.StatusFound, founduser)

	}
}

func AddProductAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var product models.Product

		if err := c.ShouldBind(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product.Product_ID = primitive.NewObjectID()

		productName := c.PostForm("product_name")

		priceStr := c.PostForm("price")
		price, err := strconv.ParseUint(priceStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price value"})
			return
		}

		ratingStr := c.PostForm("rating")
		rating, err := strconv.ParseUint(ratingStr, 5, 8)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating value"})
			return
		}
		ratingUint8 := uint8(rating)

		image, imageHeader, err := c.Request.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload error"})
			return
		}

		categories := c.PostFormArray("category")

		imagePath := "assets/" + imageHeader.Filename
		imageFile, err := os.Create(imagePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Image save error"})
			return
		}
		defer imageFile.Close()

		_, err = io.Copy(imageFile, image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Image save error"})
			return
		}

		product.Image = &imagePath
		product.Price = &price
		product.Rating = &ratingUint8
		product.Category = categories
		product.Product_Name = &productName

		_, anyerr := ProductCollection.InsertOne(ctx, product)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, "Successfully added our Product Admin!!")
	}
}

func DeleteProductAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is invalid")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		ProductID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Product ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = ProductCollection.DeleteOne(ctx, bson.M{"_id": ProductID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, "Successfully Deleted the Product")
	}
}

func EditProductAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")

		if productQueryId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid"})
			c.Abort()
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "Invalid Product ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var updatedProduct models.Product

		if err := c.ShouldBindJSON(&updatedProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		update := bson.M{
			"$set": bson.M{
				"product_name": updatedProduct.Product_Name,
				"price":        updatedProduct.Price,
				"rating":       updatedProduct.Rating,
				"category":     updatedProduct.Category,
			},
		}

		_, err = ProductCollection.UpdateOne(ctx, bson.M{"_id": productId}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to Update Product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Message": "Product updates successfully"})
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productlist []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something Went Wrong Please Try After Some Time")
			return
		}

		err = cursor.All(ctx, &productlist)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		defer cancel()

		c.IndentedJSON(200, productlist)

	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchproducts []models.Product
		queryParam := c.Query("name")
		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}

		err = searchquerydb.All(ctx, &searchproducts)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		defer searchquerydb.Close(ctx)

		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}

		defer cancel()
		c.IndentedJSON(200, searchproducts)
	}
}

func SearchProductByCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchproductsbycategory []models.Product

		queryParam := c.Query("name")
		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"category": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, "Something went wrong fetching the dbquery")
		}

		err = searchquerydb.All(ctx, &searchproductsbycategory)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, "Invalid Request")
			return
		}

		defer searchquerydb.Close(ctx)

		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid Request")
			return
		}

		defer cancel()
		c.IndentedJSON(200, searchproductsbycategory)
	}
}
