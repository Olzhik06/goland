package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Wallet struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PatientID primitive.ObjectID `bson:"patient_id" json:"patient_id"`
	Balance   float64            `bson:"balance" json:"balance"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func NewWallet(patientID primitive.ObjectID) Wallet {
	return Wallet{
		ID:        primitive.NewObjectID(),
		PatientID: patientID,
		Balance:   0.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
