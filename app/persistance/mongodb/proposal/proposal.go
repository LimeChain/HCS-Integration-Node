package proposal

import (
	"context"
	"fmt"
	"github.com/Limechain/HCS-Integration-Node/app/domain/proposal/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProposalRepository struct {
	db *mongo.Database
}

func (r *ProposalRepository) GetAll() ([]*model.Proposal, error) {
	collection := r.db.Collection("proposals")
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	results := make([]*model.Proposal, 0)

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem model.Proposal
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, &elem)
	}

	return results, nil
}

func (r *ProposalRepository) GetByID(id string) (*model.Proposal, error) {

	var result model.Proposal
	collection := r.db.Collection("proposals")
	if err := collection.FindOne(context.TODO(), bson.M{"proposalId": id}).Decode(&result); err != nil {
		return nil, err
	}

	fmt.Println(result)
	return &result, nil
}
func (r *ProposalRepository) Save(rfp *model.Proposal) (string, error) {
	collection := r.db.Collection("proposals")
	if len(rfp.ProposalId) == 0 {
		rfp.ProposalId = uuid.New().String()
	}
	_, err := collection.InsertOne(context.TODO(), rfp)
	if err != nil {
		return "", err
	}

	return rfp.ProposalId, nil
}

func NewProposalRepository(db *mongo.Database) *ProposalRepository {
	return &ProposalRepository{db: db}
}
