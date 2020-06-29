package rfp

import (
	"context"
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

	results := make([]*model.RFP, 0)

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

	return &result, nil
}
func (r *RFPRepository) Save(rfp *model.RFP) (string, error) {
	collection := r.db.Collection("rfps")
	if len(rfp.RFPId) == 0 {
		rfp.RFPId = uuid.New().String()
	}
	_, err := collection.InsertOne(context.TODO(), rfp)
	if err != nil {
		return "", err
	}

	return rfp.RFPId, nil
}

func NewRFPRepository(db *mongo.Database) *RFPRepository {
	return &RFPRepository{db: db}
}
