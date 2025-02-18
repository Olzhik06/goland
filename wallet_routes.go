package routes

import (
	"github.com/AskatNa/OnlineClinic/api/controllers"
	"github.com/AskatNa/OnlineClinic/api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func WalletRoutes(router *gin.Engine) {
	routes := router.Group("/wallet")
	{
		routes.GET("/", WalletPageHandler)
		routes.POST("/topup", TopUpRequestHandler)
		routes.POST("/confirm", ConfirmTopUpHandler)
		routes.POST("/pay", ChargePatientHandler)
		routes.POST("/register", RegisterPatientHandler)
	}
}

func RegisterPatientHandler(c *gin.Context) {
	var patient models.User
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(400, gin.H{"error": "Неверные данные"})
		return
	}

	res, err := controllers.RegisterPatient(patient)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Пациент зарегистрирован", "id": res.InsertedID})
}

func TopUpRequestHandler(c *gin.Context) {
	var req controllers.TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Неверные данные"})
		return
	}

	res, err := controllers.RequestTopUp(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Запрос на пополнение создан", "transaction_id": res.InsertedID})
}

func ConfirmTopUpHandler(c *gin.Context) {
	var req struct {
		TransactionID string `json:"transaction_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Неверные данные"})
		return
	}

	transactionID, err := primitive.ObjectIDFromHex(req.TransactionID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Неверный ID транзакции"})
		return
	}

	err = controllers.ConfirmTopUp(transactionID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Пополнение подтверждено"})
}

func ChargePatientHandler(c *gin.Context) {
	var req controllers.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Неверные данные"})
		return
	}

	err := controllers.ChargePatient(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Оплата успешна"})
}
func WalletPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "wallets.html", nil)
}
