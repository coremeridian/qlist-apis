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

const TestSessionTable = "sessions"

type TestSessionModel struct {
	DB *mongo.Database
}

func (m *TestSessionModel) Insert(ts *models.TestSession) (id interface{}, err error) {
	if m.DB == nil {
		fmt.Println("No Database connection")
		return nil, ErrNoConnection
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ts.MaxAttempts = 1
	ts.CreatedAt = time.Now()
	collection := m.DB.Collection(TestSessionTable)
	res, err := collection.InsertOne(ctx, ts)
	if err != nil {
		fmt.Printf("Error inserting test: %v", ts)
		return nil, ErrNoInsertion
	}

	id = res.InsertedID.(primitive.ObjectID)
	fmt.Println("Inserted test session into collection", id)

	return
}

func (m *TestSessionModel) Update(idString string, ts interface{}) (didUpdate bool, err error) {
	if m.DB == nil {
		fmt.Println("No Database connection")
		return false, ErrNoConnection
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update, err := toDoc(ts)
	if err != nil {
		return false, err
	}

	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return false, err
	}

	_, err = m.DB.Collection(TestSessionTable).UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
		nil,
	)
	return err == nil, err
}

func (m *TestSessionModel) Get(keys *models.KeySet) (testSession *models.TestSession, err error) {
	if m.DB == nil {
		fmt.Println("No Database connection")
		return nil, ErrNoConnection
	}

	collection := m.DB.Collection(TestSessionTable)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	keySet, err := toDoc(keys)
	if err != nil {
		return nil, err
	}

	testSession = &models.TestSession{}
	err = collection.FindOne(ctx, keySet).Decode(testSession)
	return
}

func (m *TestSessionModel) Latest(n int) (testSessions []*models.TestSession, err error) {
	if m.DB == nil {
		fmt.Println("No Database connection")
		return nil, ErrNoConnection
	}

	collection := m.DB.Collection(TestSessionTable)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	testSessions = make([]*models.TestSession, n)
	if err = cursor.All(ctx, &testSessions); err != nil {
		log.Fatal(err)
	}
	return
}
