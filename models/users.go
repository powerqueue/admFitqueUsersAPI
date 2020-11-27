package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//LoginCollectionName - constant definition
const LoginCollectionName = "login"

//ILoginRepository - interface definition
type ILoginRepository interface {
	GetLogin(location string, username string, memberID string, loginDate *time.Time) ([]*LoginDefinition, error)
	CreateLogin(AddressHierarchy *LoginDefinition) (*LoginDefinition, error)
	TermLogin(AddressHierarchy *LoginDefinition) (*LoginDefinition, error)
}

//LoginRepository - struct definition
type LoginRepository struct {
	MongoClient IMongoClient
	// Validator   core.IValidator
}

//NewLoginRepository - func definition
func NewLoginRepository(client IMongoClient) *LoginRepository {
	return &LoginRepository{
		MongoClient: client,
	}
}

//LoginDefinition - struc definition
type LoginDefinition struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id" validate:"required"`
	LoginID      string             `bson:"LoginID" json:"LoginID" validate:"required"`
	MemberID     string             `bson:"MemberID" json:"MemberID" validate:"required"`
	LocationID   string             `bson:"LocationID" json:"LocationID" validate:"required"`
	UserName     string             `bson:"UserName" json:"UserName" validate:"required"`
	EfctvStartDt *time.Time         `bson:"EfctvStartDt" json:"EfctvStartDt" validate:"required"`
	EfctvEndDt   *time.Time         `bson:"EfctvEndDt" json:"EfctvEndDt"`
}

//GetLogin - func definition
func (loginRepo *LoginRepository) GetLogin(location string, username string, memberID string, loginDate *time.Time) ([]*LoginDefinition, error) {
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
	var login []*LoginDefinition
	// var result Trainer

	filter := bson.D{{"MemberID", memberID}, {"LocationID", location}, {"UserName", username}, {"EfctvStartDt", loginDate}}

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

	if len(login) == 0 {
		return login, err
	}

	return login, err
}

//CreateLogin - create Login document
func (loginRepo *LoginRepository) CreateLogin(Login *LoginDefinition) (*LoginDefinition, error) {
	fmt.Println("Inside CreateLogin Model Function")
	now := time.Now()
	Login.ID = primitive.NewObjectID()
	Login.MemberID = strings.ToUpper(Login.MemberID)

	if Login.UserName == "" {
		Login.UserName = "anonymous"
	} else {
		Login.UserName = Login.UserName
	}

	Login.EfctvStartDt = &now

	Login.EfctvEndDt = &time.Time{}

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://ln004prd:27017")

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

	collection := client.Database("fitqueue-db").Collection("login")

	//Insert SINGLE here
	// insertResult, err := collection.InsertOne(context.TODO(), ash)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	//Insert MULTIPLE here
	trainers := []interface{}{Login}

	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error create %s", err)
	}

	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	//disconnect client
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

	bytes, err := json.Marshal(Login)
	if err != nil {
		return Login, err

	}

	fmt.Println(string(bytes))
	return Login, err
}

//TermLogin - term Login document
func (loginRepo *LoginRepository) TermLogin(Login *LoginDefinition) (*LoginDefinition, error) {
	fmt.Println("Inside TermLogin Model Function")
	now := time.Now()
	Login.ID = primitive.NewObjectID()
	Login.MemberID = strings.ToUpper(Login.MemberID)

	if Login.UserName == "" {
		Login.UserName = "anonymous"
	}

	Login.EfctvEndDt = &now

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://ln004prd:27017")

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

	collection := client.Database("fitqueue-db").Collection("login")

	//build filter
	filter := bson.D{{"LoginID", Login.LoginID}, {"MemberID", Login.MemberID}, {"LocationID", Login.LocationID}, {"UserName", Login.UserName}}

	// update := bson.D{
	// 	{"$inc", bson.D{
	// 		{"EfctvEndDt", Login.EfctvEndDt},
	// 	}},
	// }

	update := bson.D{
		{"$set", bson.D{
			{"EfctvEndDt", Login.EfctvEndDt},
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

	bytes, err := json.Marshal(Login)
	if err != nil {
		return Login, err

	}

	fmt.Println(string(bytes))
	return Login, err
}
