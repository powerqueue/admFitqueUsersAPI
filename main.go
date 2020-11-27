package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/powerqueue/fitque-users-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var (
	audience string
	domain   string
)

func main() {

	// connect to mongo and migrate schema
	dbConfigs := &models.DBConfigs1{
		Host:     "ln004prd",
		Port:     27017,
		User:     "",
		Password: "",
	}
	models.ConnectAndMigrate(dbConfigs, "fitqueue-db")

	// apiPrefix := "/fitqueue-login-api/v1"

	r := gin.Default()

	r.Use(CORSMiddleware())

	// r.GET("/todo", GetTodoListHandler)
	// r.POST("/todo", AddTodoHandler)
	// r.DELETE("/todo/:id", DeleteTodoHandler)
	// r.PUT("/todo", CompleteTodoHandler)

	v1 := r.Group("/fitqueue-login-api/v1")
	{
		v1.POST("/retrieve-login", RetrieveLoginHandler)
		v1.POST("/create-login", CreateLogin)
		v1.POST("/term-login", TermLogin)
	}

	// authorized := r.Group("/")
	// authorized.Use(authRequired())
	// authorized.GET("/cases/:caseType/:page", GetCaseListHandler)
	// authorized.GET("/case/:caseCode", GetCaseDetailsHandler)
	// authorized.POST("/case", AddCaseHandler)
	// authorized.DELETE("/case/:caseCode", TermCaseHandler)
	// authorized.PUT("/case", UpdateCaseHandler)

	err := r.Run(":3001")
	if err != nil {
		panic(err)
	}
}

//CORSMiddleware Cross-Origin Resource Sharing helper Class
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, OPTIONS, POST, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func terminateWithError(statusCode int, message string, c *gin.Context) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}

func convertHTTPBodyToLoginDefinition(httpBody io.ReadCloser) (models.LoginDefinition, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return models.LoginDefinition{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToLoginDefinition(body)
}

func convertJSONBodyToLoginDefinition(jsonBody []byte) (models.LoginDefinition, int, error) {
	var login models.LoginDefinition
	err := json.Unmarshal(jsonBody, &login)
	if err != nil {
		return models.LoginDefinition{}, http.StatusBadRequest, err
	}
	return login, http.StatusOK, nil
}

//RetrieveLoginHandler - handler definition
func RetrieveLoginHandler(c *gin.Context) {
	fmt.Println("Inside CreateLogin Route Handler")
	login, statusCode, err := convertHTTPBodyToLoginDefinition(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("fitque-db").Collection("login")

	// create a value into which the result can be decoded
	filter := bson.D{{"MemberID", login.MemberID}, {"LocationID", login.LocationID}, {"UserName", login.UserName}, {"EfctvStartDt", login.EfctvStartDt}}

	err = collection.FindOne(context.TODO(), filter).Decode(&login)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", login)

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

	c.JSON(statusCode, login)
}

//CreateLogin - handler method definition
func CreateLogin(c *gin.Context) {
	fmt.Println("Inside CreateLogin Route Handler")
	login, statusCode, err := convertHTTPBodyToLoginDefinition(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}

	fmt.Println("Inside CreateLogin Model Function")
	now := time.Now()
	login.ID = primitive.NewObjectID()
	login.MemberID = strings.ToUpper(login.MemberID)

	if login.UserName == "" {
		login.UserName = "anonymous"
	} else {
		login.UserName = login.UserName
	}

	login.EfctvStartDt = &now

	login.EfctvEndDt = &time.Time{}

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://ln004prd:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error create %s", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error create %s", err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("fitqueue-db").Collection("login")

	//Insert SINGLE here
	// insertResult, err := collection.InsertOne(context.TODO(), ash)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	//Insert MULTIPLE here
	trainers := []interface{}{login}

	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error create %s", err)
	}

	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	//disconnect client
	err = client.Disconnect(context.TODO())

	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error create %s", err)
	}
	fmt.Println("Connection to MongoDB closed.")

	bytes, err := json.Marshal(login)
	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error create %s", err)

	}

	fmt.Println(string(bytes))
	// return login, err

	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "Category: %v\n", loginDef)

	c.JSON(statusCode, login)
}

//TermLogin - handler method definition
func TermLogin(c *gin.Context) {
	fmt.Println("Inside CreateLogin Route Handler")
	login, statusCode, err := convertHTTPBodyToLoginDefinition(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}

	fmt.Println("Inside TermLogin Model Function")
	now := time.Now()
	login.ID = primitive.NewObjectID()
	login.MemberID = strings.ToUpper(login.MemberID)

	if login.UserName == "" {
		login.UserName = "anonymous"
	}

	login.EfctvEndDt = &now

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://ln004prd:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error term %s", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error term %s", err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("fitqueue-db").Collection("login")

	//build filter
	filter := bson.D{{"LoginID", login.LoginID}, {"MemberID", login.MemberID}, {"LocationID", login.LocationID}, {"UserName", login.UserName}}

	// update := bson.D{
	// 	{"$inc", bson.D{
	// 		{"EfctvEndDt", Login.EfctvEndDt},
	// 	}},
	// }

	update := bson.D{
		{"$set", bson.D{
			{"EfctvEndDt", login.EfctvEndDt},
		}},
	}

	// update := bson.M{"EfctvEndDt": primitive.Timestamp{T: uint32(time.Now().Unix())}}

	//update/term based on filter
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error term %s", err)
	} else {
		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	}

	//disconnect client
	err = client.Disconnect(context.TODO())

	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error term %s", err)
	}
	fmt.Println("Connection to MongoDB closed.")

	bytes, err := json.Marshal(login)
	if err != nil {
		c.JSON(statusCode, login)

	}

	fmt.Println(string(bytes))
	c.JSON(statusCode, login)
}
