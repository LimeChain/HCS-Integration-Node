package contract

import (
	"context"
	"errors"
	"github.com/Limechain/HCS-Integration-Node/app/domain/contract/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContractRepository struct {
	db *mongo.Database
}

func (r *ContractRepository) GetAll() ([]*model.Contract, error) {
	collection := r.db.Collection("contracts")
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	results := make([]*model.Contract, 0)

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem model.Contract
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, &elem)
	}

	return results, nil
}

func (r *ContractRepository) GetByID(id string) (*model.Contract, error) {

	var result model.Contract
	collection := r.db.Collection("contracts")
	if err := collection.FindOne(context.TODO(), bson.M{"proposalId": id}).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *ContractRepository) Save(contract *model.Contract) (string, error) {
	collection := r.db.Collection("contracts")
	if len(contract.ContractId) == 0 {
		contract.ContractId = uuid.New().String()
	}
	_, err := collection.InsertOne(context.TODO(), contract)
	if err != nil {
		return "", err
	}

	return contract.ContractId, nil
}

func (r *ContractRepository) Update(contract *model.Contract) error {
	collection := r.db.Collection("contracts")
	if len(contract.ContractId) == 0 {
		return errors.New("Contract without Id cannot be updated")
	}
	ur, err := collection.ReplaceOne(context.TODO(), bson.M{"unsignedContract.contractId": contract.ContractId}, contract)
	if err != nil {
		return err
	}

	if ur.MatchedCount == 0 {
		return errors.New("No such contract found")
	}

	return nil
}

func NewContractRepositiry(db *mongo.Database) *ContractRepository {
	return &ContractRepository{db: db}
}
