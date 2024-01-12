package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/baseba/got/handler"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoomData struct {
	ID    string
	Money int
	Count int
}

var db *sql.DB

func main() {
	// err := godotenv.Load()
	envFile, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal(".env file could not be loaded")
	}
	// username := os.Getenv("USERNAME")
	username := envFile["USERNAME"]
	// password := os.Getenv("PASSWORD")
	password := envFile["PASSWORD"]
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
	app.GET("/slots/:room", func(c echo.Context) error {
		room := c.Param("room")
		database := client.Database("pokeslots")
		return slotsHandler.HandleSlotsShow(c, room, database)

	})
	indexHandler := handler.IndexHandler{}
	app.GET("/", indexHandler.HandleIndexShow)

	app.POST("/go-to", func(c echo.Context) error {
		room := c.FormValue("room")
		url := fmt.Sprintf("/slots/%s", room)
		c.Response().Header().Set("HX-Redirect", url)
		c.Response().WriteHeader(200)
		// c.Redirect(201, url)
		// fmt.Println()
		return nil
	})

	app.POST("/win/:amount", func(c echo.Context) error {
		amount, _ := strconv.Atoi(c.Param("amount"))
		room := c.QueryParam("room")

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
			ID:    room,
			Count: 0,
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
		room := c.QueryParam("room")

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
			ID:    room,
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
		hotbar := ""
		halfconut := count / 3
		for i := 0; i <= 11; i++ {
			if i+1 > halfconut {
				hotbar += "<div class='bg-red-600 p-6'></div>\n"
			} else {
				hotbar += "<div class='bg-blue-600 p-6'></div>\n"
			}
		}
		res := fmt.Sprintf(`<p className="text-sm text-gray-500 mb-1">Perdidas Seguidas</p>
		<p id="count" className="text-sm text-gray-500 mb-1">%d</p>
		<p className="text-sm text-gray-500 mb-1">saldo</p>
		<p id="money" className="text-lg text-black">%d</p>
		<h2>
                    nivel de calor de la maquina
                </h2>
                    <div class="grid grid-cols-12 justify-start">
                        %s
                </div>`, count, money, hotbar)

		return c.String(200, res)
	})

	app.POST("/reset/:room", func(c echo.Context) error {
		room := c.Param("room")
		database := client.Database("pokeslots")
		collection := database.Collection("rooms")
		filter := bson.M{"id": room}

		replacement := RoomData{
			ID:    room,
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
