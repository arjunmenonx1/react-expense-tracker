package main

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Name of the database.
	DBName             = "expenses"
	ExpenseCollections = "expense"
	//URI                 = "mongodb://<user>:<password>@localhost:27017"
	URI = "mongodb://localhost:27017"
)

var (
	/* Used to create a singleton object of MongoDB client.
	Initialized and exposed through  GetMongoClient().*/
	clientInstance *mongo.Client
	//Used during creation of singleton client object in GetMongoClient().
	clientInstanceError error
	//Used to execute client creation procedure only once.
	mongoOnce sync.Once
)

type Expense struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	ExpenseID string             `bson:"expenseID" json:"expenseID"`
	Title     string             `bson:"title" json:"title"`
	Amount    float64            `bson:"amount" json:"amount,omitempty"`
	Date      time.Time          `bson:"date" json:"date,omitempty"`
}

//GetMongoClient - Return mongodb connection to work with
func GetMongoClient() (*mongo.Client, error) {
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(URI)
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
	})
	return clientInstance, clientInstanceError
}

//CreateExpense - Insert a new document in the collection.
func CreateExpense(expense Expense) error {
	//Get MongoDB connection using
	client, err := GetMongoClient()
	if err != nil {
		return err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(DBName).Collection(ExpenseCollections)
	//Perform InsertOne operation & validate against the error.
	_, err = collection.InsertOne(context.TODO(), expense)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

//CreateManyExpenses - Insert multiple documents at once in the collection.
func CreateManyExpenses(list []Expense) error {
	//Map struct slice to interface slice as InsertMany accepts interface slice as parameter
	insertableList := make([]interface{}, len(list))
	for i, v := range list {
		insertableList[i] = v
	}
	//Get MongoDB connection using
	client, err := GetMongoClient()
	if err != nil {
		return err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(DBName).Collection(ExpenseCollections)
	//Perform InsertMany operation & validate against the error.
	_, err = collection.InsertMany(context.TODO(), insertableList)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

//GetIssuesByCode - Get All issues for collection
func GetExpensesByID(id string) (Expense, error) {
	result := Expense{}
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "expenseID", Value: id}}
	//Get MongoDB connection using connectionhelper.
	client, err := GetMongoClient()
	if err != nil {
		return result, err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(DBName).Collection(ExpenseCollections)
	//Perform FindOne operation & validate against the error.
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	//Return result without any error.
	return result, nil
}

//GetAllIssues - Get All issues for collection
func GetAllIssues() ([]Expense, error) {
	//Define filter query for fetching specific document from collection
	filter := bson.D{{}} //bson.D{{}} specifies 'all documents'
	issues := []Expense{}
	//Get MongoDB connection using connectionhelper.
	client, err := GetMongoClient()
	if err != nil {
		return issues, err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(DBName).Collection(ExpenseCollections)
	//Perform Find operation & validate against the error.
	cur, findError := collection.Find(context.TODO(), filter)
	if findError != nil {
		return issues, findError
	}
	//Map result to slice
	for cur.Next(context.TODO()) {
		t := Expense{}
		err := cur.Decode(&t)
		if err != nil {
			return issues, err
		}
		issues = append(issues, t)
	}
	// once exhausted, close the cursor
	cur.Close(context.TODO())
	if len(issues) == 0 {
		return issues, mongo.ErrNoDocuments
	}
	return issues, nil
}

//DeleteExpense - delete one expense for collection
func DeleteExpense(code string) error {
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "expenseID", Value: code}}
	//Get MongoDB connection using connectionhelper.
	client, err := GetMongoClient()
	if err != nil {
		return err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(DBName).Collection(ExpenseCollections)
	//Perform DeleteOne operation & validate against the error.
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

//DeleteAll - Delete all expenses for collection
func DeleteAll() error {
	//Define filter query for fetching specific document from collection
	selector := bson.D{{}} // bson.D{{}} specifies 'all documents'
	//Get MongoDB connection using connectionhelper.
	client, err := GetMongoClient()
	if err != nil {
		return err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(DBName).Collection(ExpenseCollections)
	//Perform DeleteMany operation & validate against the error.
	_, err = collection.DeleteMany(context.TODO(), selector)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

func main() {
	// expense := Expense{
	// 	ID:        primitive.NewObjectID(),
	// 	ExpenseID: "1d",
	// 	Title:     "First expense",
	// 	Amount:    3.50,
	// 	Date:      time.Now(),
	// }
	//CreateExpense(expense)
	DeleteExpense("1d")
}
