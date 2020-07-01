package purchase_order

import (
	"context"
	"errors"
	"github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PurchaseOrderRepository struct {
	db *mongo.Database
}

func (r *PurchaseOrderRepository) GetAll() ([]*model.PurchaseOrder, error) {
	collection := r.db.Collection("pos")
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	results := make([]*model.PurchaseOrder, 0)

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem model.PurchaseOrder
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, &elem)
	}

	return results, nil
}

func (r *PurchaseOrderRepository) GetByID(id string) (*model.PurchaseOrder, error) {

	var result model.PurchaseOrder
	collection := r.db.Collection("pos")
	if err := collection.FindOne(context.TODO(), bson.M{"unsignedPurchaseOrder.purchaseOrderId": id}).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *PurchaseOrderRepository) Save(po *model.PurchaseOrder) (string, error) {
	collection := r.db.Collection("pos")
	if len(po.PurchaseOrderId) == 0 {
		po.PurchaseOrderId = uuid.New().String()
	}
	_, err := collection.InsertOne(context.TODO(), po)
	if err != nil {
		return "", err
	}

	return po.PurchaseOrderId, nil
}

func (r *PurchaseOrderRepository) Update(po *model.PurchaseOrder) error {
	collection := r.db.Collection("pos")
	if len(po.PurchaseOrderId) == 0 {
		return errors.New("Purchase order without Id cannot be updated")
	}
	ur, err := collection.ReplaceOne(context.TODO(), bson.M{"unsignedPurchaseOrder.purchaseOrderId": po.PurchaseOrderId}, po)
	if err != nil {
		return err
	}

	if ur.MatchedCount == 0 {
		return errors.New("No such purchase order found")
	}

	return nil
}

func NewPurchaseOrderRepository(db *mongo.Database) *PurchaseOrderRepository {
	return &PurchaseOrderRepository{db: db}
}