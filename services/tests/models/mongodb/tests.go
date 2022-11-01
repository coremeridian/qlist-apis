package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"coremeridian.xyz/app/qlist/cmd/api/services/tests/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestModel struct {
	DB *mongo.Database
}

const TestTable = "assessments"

func (m *TestModel) Counter(ctx context.Context) *mongo.SingleResult {
	counter := m.DB.Collection("counters").FindOneAndUpdate(
		ctx,
		bson.M{"_id": "testid"},
		bson.M{"$inc": bson.M{"sequence_value": 1}},
		nil,
	)

	return counter
}

func (m *TestModel) Insert(test *models.Test) (id interface{}, err error) {
	if m.DB == nil {
		fmt.Println("No Database connection")
		return nil, ErrNoConnection
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	counter := m.Counter(ctx)
	if counter.Err() != nil {
		return nil, counter.Err()
	}
	var count struct {
		N int `bson:"sequence_value" json:"sequence_value"`
	}
	err = counter.Decode(&count)
	if err != nil {
		return nil, err
	}

	test.TestID = count.N
	test.CreatedAt = time.Now()
	collection := m.DB.Collection(TestTable)
	res, err := collection.InsertOne(ctx, test)
	if err != nil {
		fmt.Printf("Error inserting test: %v", test)
		return nil, ErrNoInsertion
	}

	id = res.InsertedID
	fmt.Println("Inserted test into collection", id)
	return
}

func (m *TestModel) Update(idString string, t interface{}) (didUpdate bool, err error) {
	if m.DB == nil {
		fmt.Println("No Database connection")
		return false, ErrNoConnection
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update, err := toDoc(t)
	if err != nil {
		return false, err
	}

	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return false, err
	}

	_, err = m.DB.Collection(TestTable).UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
		nil,
	)
	return err == nil, err
}

func (m *TestModel) Get(idString string) (test *models.Test, err error) {
	if m.DB == nil {
		fmt.Println("No Database connection")
		return nil, ErrNoConnection
	}

	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return nil, err
	}

	collection := m.DB.Collection(TestTable)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	test = &models.Test{}
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(test)
	return
}

func (m *TestModel) Latest(n int, options models.TestOptions) (tests []*models.Test, err error) {
	if m.DB == nil {
		fmt.Println("No Database connection")
		return nil, ErrNoConnection
	}

	collection := m.DB.Collection(TestTable)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	optionsDoc := &bson.M{}
	if options != nil {
		optionsDoc, err = toDoc(options)
		if err != nil {
			return
		}
	}
	cursor, err := collection.Find(ctx, optionsDoc)
	if err != nil {
		log.Fatal(err)
	}

	tests = make([]*models.Test, n)
	if err = cursor.All(ctx, &tests); err != nil {
		log.Fatal(err)
	}

	return
}
