package controllers

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/AskatNa/OnlineClinic/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var walletCollection *mongo.Collection
var transactionCollection *mongo.Collection

func SetCollections(userColl, walletColl, transColl *mongo.Collection) {
	userCollection = userColl
	walletCollection = walletColl
	transactionCollection = transColl
}

func RegisterPatient(patient models.User) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	patient.ID = primitive.NewObjectID()
	res, err := userCollection.InsertOne(ctx, patient)
	if err != nil {
		return nil, err
	}

	wallet := models.NewWallet(patient.ID)
	_, err = walletCollection.InsertOne(ctx, wallet)
	if err != nil {
		log.Println("Ошибка при создании кошелька:", err)
	}

	return res, nil
}

type TopUpRequest struct {
	PatientID primitive.ObjectID `json:"patient_id"`
	Amount    float64            `json:"amount"`
}

func RequestTopUp(req TopUpRequest) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wallet models.Wallet
	err := walletCollection.FindOne(ctx, bson.M{"patient_id": req.PatientID}).Decode(&wallet)
	if err != nil {
		return nil, err
	}

	transaction := bson.M{
		"patient_id": req.PatientID,
		"amount":     req.Amount,
		"status":     "pending",
		"created_at": time.Now(),
	}

	res, err := transactionCollection.InsertOne(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func ConfirmTopUp(transactionID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var transaction bson.M
	err := transactionCollection.FindOne(ctx, bson.M{"_id": transactionID}).Decode(&transaction)
	if err != nil {
		return err
	}

	if transaction["status"] != "pending" {
		return errors.New("транзакция уже обработана")
	}

	_, err = walletCollection.UpdateOne(
		ctx,
		bson.M{"patient_id": transaction["patient_id"]},
		bson.M{"$inc": bson.M{"balance": transaction["amount"]}},
	)
	if err != nil {
		return err
	}

	_, err = transactionCollection.UpdateOne(
		ctx,
		bson.M{"_id": transactionID},
		bson.M{"$set": bson.M{"status": "confirmed"}},
	)

	return err
}

type PaymentRequest struct {
	PatientID primitive.ObjectID `json:"patient_id"`
	Amount    float64            `json:"amount"`
	Service   string             `json:"service"`
}

func ChargePatient(req PaymentRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wallet models.Wallet
	err := walletCollection.FindOne(ctx, bson.M{"patient_id": req.PatientID}).Decode(&wallet)
	if err != nil {
		return errors.New("кошелек не найден")
	}

	if wallet.Balance < req.Amount {
		return errors.New("недостаточно средств")
	}

	_, err = walletCollection.UpdateOne(
		ctx,
		bson.M{"patient_id": req.PatientID},
		bson.M{"$inc": bson.M{"balance": -req.Amount}},
	)
	if err != nil {
		return err
	}

	transaction := bson.M{
		"patient_id": req.PatientID,
		"amount":     req.Amount,
		"service":    req.Service,
		"status":     "charged",
		"created_at": time.Now(),
	}

	_, err = transactionCollection.InsertOne(ctx, transaction)
	return err
}
