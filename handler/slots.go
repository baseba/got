package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/baseba/got/view/slotView"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SlotsHandler struct{}
type RoomData struct {
	ID    string
	Money int
	Count int
}

func (h SlotsHandler) HandleSlotsShow(c echo.Context, room string, client *mongo.Database) error {
	if room != "" {

		return render(c, slotView.Show(c.Param("room"), "putas", "gratis"))
	}
	collection := client.Collection("rooms")
	filter := bson.M{"id": room}
	var result RoomData
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
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
	money := strconv.Itoa(result.Money)
	count := strconv.Itoa(result.Count)
	return render(c, slotView.Show(room, money, count))
}
