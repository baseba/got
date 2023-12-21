package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/baseba/got/handler"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoomData struct {
	ID    int64
	Money int
	Count int
}

var db *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file could not be loaded")
	}
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	fmt.Println(username)
	URI := fmt.Sprintf("mongodb+srv://%s:%s@pokeslots.ibl4gcl.mongodb.net/?retryWrites=true&w=majority", username, password)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(URI).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	app := echo.New()

	userHandler := handler.UserHandler{}
	app.GET("/user", userHandler.HandleUserShow)
	slotsHandler := handler.SlotsHandler{}
	app.GET("/slots/:room", slotsHandler.HandleSlotsShow)

	app.POST("/win/:amount", func(c echo.Context) error {
		amount, _ := strconv.Atoi(c.Param("amount"))
		room, _ := strconv.Atoi(c.QueryParam("room"))

		// add data to the db
		database := client.Database("pokeslots")
		collection := database.Collection("rooms")
		filter := bson.M{"id": room}

		var result RoomData
		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				fmt.Println("No matching document found.")
			} else {
				log.Fatal(err)
			}
		} else {
			// Print the result
			fmt.Printf("encontramos" + strconv.Itoa(result.Money))
		}
		money := result.Money + amount
		count := 0
		replacement := RoomData{
			ID:    int64(room),
			Count: count,
			Money: money,
		}
		// Example: Perform the update or insert operation with upsert set to true
		updateResult, err := collection.ReplaceOne(context.TODO(), filter, replacement, options.Replace().SetUpsert(true))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Matched %v document and modified %v document.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

		fmt.Println(c.Path())
		res := fmt.Sprintf(`<p className="text-sm text-gray-500 mb-1">Perdidas Seguidas</p>
		<p id="count" className="text-sm text-gray-500 mb-1">%d</p>
		<p className="text-sm text-gray-500 mb-1">saldo</p>
		<p id="money" className="text-lg text-black">%d</p>`, count, money)

		return c.String(200, res)
	})

	app.POST("/lose/:amount", func(c echo.Context) error {
		amount, _ := strconv.Atoi(c.Param("amount"))
		room, _ := strconv.Atoi(c.QueryParam("room"))

		// add data to the db
		database := client.Database("pokeslots")
		collection := database.Collection("rooms")
		filter := bson.M{"id": room}

		var result RoomData
		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				fmt.Println("No matching document found.")
			} else {
				log.Fatal(err)
			}
		} else {
			// Print the result
			fmt.Printf("encontramos" + strconv.Itoa(result.Money))
		}
		money := result.Money - amount
		count := result.Count + 1
		replacement := RoomData{
			ID:    int64(room),
			Count: count,
			Money: money,
		}
		// Example: Perform the update or insert operation with upsert set to true
		updateResult, err := collection.ReplaceOne(context.TODO(), filter, replacement, options.Replace().SetUpsert(true))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Matched %v document and modified %v document.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

		fmt.Println(c.Path())
		res := fmt.Sprintf(`<p className="text-sm text-gray-500 mb-1">Perdidas Seguidas</p>
		<p id="count" className="text-sm text-gray-500 mb-1">%d</p>
		<p className="text-sm text-gray-500 mb-1">saldo</p>
		<p id="money" className="text-lg text-black">%d</p>`, count, money)

		return c.String(200, res)
	})

	app.POST("/reset/:room", func(c echo.Context) error {
		room, _ := strconv.Atoi(c.QueryParam("room"))
		database := client.Database("pokeslots")
		collection := database.Collection("rooms")
		filter := bson.M{"id": room}

		replacement := RoomData{
			ID:    int64(room),
			Count: 0,
			Money: 0,
		}
		// Example: Perform the update or insert operation with upsert set to true
		updateResult, err := collection.ReplaceOne(context.TODO(), filter, replacement, options.Replace().SetUpsert(true))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Matched %v document and modified %v document.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

		res := fmt.Sprintf(`<p className="text-sm text-gray-500 mb-1">Perdidas Seguidas</p>
		<p id="count" className="text-sm text-gray-500 mb-1">%d</p>
		<p className="text-sm text-gray-500 mb-1">saldo</p>
		<p id="money" className="text-lg text-black">%d</p>`, 0, 0)
		return c.String(200, res)
	})

	app.Start(":3000") //envs?
	fmt.Println("im alive!")

}
