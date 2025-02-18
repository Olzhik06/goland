package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AskatNa/OnlineClinic/api/controllers"
	"github.com/AskatNa/OnlineClinic/api/routes"
	"github.com/AskatNa/OnlineClinic/config/configs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error

	// Подключение к MongoDB
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Couldn't connect to MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB")

	// Инициализация коллекций
	userCollection := client.Database("online_clinic").Collection("users")
	walletCollection := client.Database("online_clinic").Collection("wallets")
	transactionCollection := client.Database("online_clinic").Collection("transactions")

	// Передаём коллекции в контроллер
	controllers.SetCollections(userCollection, walletCollection, transactionCollection)

	// Настройка роутера
	router := gin.Default()
	configs.ConnectDB()
	router.LoadHTMLGlob("ui/html/*.html")

	router.GET("/wallets", func(c *gin.Context) {
		c.HTML(http.StatusOK, "wallets.html", nil)
	})

	// Регистрация маршрутов
	routes.UnauthRoutes(router)
	routes.UserRoutes(router)
	routes.DoctorRoutes(router)
	routes.PatientRoutes(router)
	routes.AppointmentRoutes(router)
	routes.WalletRoutes(router) // Подключаем маршруты кошелька

	fmt.Println("The server is running on port :9000")
	log.Fatal(http.ListenAndServe(":9000", router))
}
