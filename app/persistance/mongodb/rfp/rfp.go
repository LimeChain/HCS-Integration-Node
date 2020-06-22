package rfp

import (
	"context"
	"fmt"
	"github.com/Limechain/HCS-Integration-Node/app/domain/rfp/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RFPRepository struct {
	db *mongo.Database
}

func (r *RFPRepository) GetAll() ([]*model.RFP, error) {
	collection := r.db.Collection("rfps")
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	var results []*model.RFP

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem model.RFP
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, &elem)
	}

	return results, nil
}

func (r *RFPRepository) GetByID(id string) (*model.RFP, error) {

	var result model.RFP
	collection := r.db.Collection("rfps")
	if err := collection.FindOne(context.TODO(), bson.M{"rfpId": id}).Decode(&result); err != nil {
		return nil, err
	}

	fmt.Println(result)
	return &result, nil
}
func (r *RFPRepository) Save(rfp *model.RFP) (string, error) {
	collection := r.db.Collection("rfps")
	if len(rfp.RFPID) == 0 {
		rfp.RFPID = uuid.New().String()
	}
	_, err := collection.InsertOne(context.TODO(), rfp)
	if err != nil {
		return "", err
	}

	return rfp.RFPID, nil
}

func NewRFPRepository(db *mongo.Database) *RFPRepository {
	return &RFPRepository{db: db}
}
