// Package models is a compendium of higher order components of data
package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateSignature = errors.New("models: duplicate model signature")
)

type KeySet struct {
	UserID    string `bson:"userid,omitempty"`
	TestID    string `bson:"testid,omitempty"`
	SessionID string `bson:"sessionid,omitempty"`
}

type Test struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TestID      int                `bson:"testid" json:"testid"`
	Title       string             `bson:"title" json:"title,omitempty"`
	Author      string             `bson:"author" json:"author,omitempty"`
	Description string             `bson:"description" json:"description,omitempty"`
	Content     string             `bson:"content" json:"content,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at,omitempty"`
	IsPublished bool               `bson:"is_published" json:"is_published"`
	Price       float32            `bson:"price,omitempty" json:"price,omitempty"`
	PriceID     string             `bson:"priceid" json:"priceid,omitempty"`
	Attempts    int                `bson:"attempts" json:"attempts"`
}

type TestSession struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"userid" json:"userid"`
	TestID      string             `bson:"testid" json:"testid"`
	SessionID   string             `bson:"sessionid" json:"sessionid"`
	SessionURL  string             `bson:"sessionurl" json:"sessionurl"`
	IsValid     bool               `bson:"is_valid" json:"is_valid"`
	Attempt     int                `bson:"attempt" json:"attempt,omitempty"`
	MaxAttempts int                `bson:"max_attempts" json:"max_attempts,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at,omitempty"`
}

type TestOptions *testOptions
type testOptions struct {
	IsPublished *bool `bson:"is_published,omitempty" json:"is_published,omitempty"`
	Attempts    *int  `bson:"attempts,omitempty" json:"attempts,omitempty"`
}

func NewTestOptionsFromQuery(query url.Values) (TestOptions, error) {
	queryParams := make(map[string]interface{})
	wunp := "with_unpublished"
	if !query.Has(wunp) || (query.Has(wunp) && query.Get(wunp) != "true") {
		queryParams["is_published"] = true
	}
	if query.Has("attempts") {
		if intVar, err := strconv.Atoi(query.Get("attempts")); err == nil {
			queryParams["attempts"] = intVar
		}
	}
	jsonString, err := json.Marshal(queryParams)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(bytes.NewBuffer(jsonString))
	decoder.DisallowUnknownFields()
	to := NewTestOptions(false)
	err = decoder.Decode(to)
	if err != nil {
		return nil, err
	}
	return to, nil
}

func NewTestOptions(companyOptions bool) TestOptions {
	to := &testOptions{}
	if companyOptions {
		to.Published(true)
	}
	return to
}

func (t *testOptions) Published(isPublished bool) *testOptions {
	t.IsPublished = &isPublished
	return t
}

func (t *testOptions) HasAttempts(attempts int) *testOptions {
	t.Attempts = &attempts
	return t
}
