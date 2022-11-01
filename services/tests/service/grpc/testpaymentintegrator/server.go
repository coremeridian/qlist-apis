// Package testpaymentintegrator maintains API and services relating to grpc operations
package testpaymentintegrator

import (
	context "context"
	"errors"
	"log"

	"coremeridian.xyz/app/qlist/cmd/api/services/tests/models"
	"coremeridian.xyz/app/qlist/cmd/api/services/tests/service"
)

type Server struct {
	UnimplementedTestPaymentIntegratorServer
	service.Endpoints
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func (tp *Server) PublishTest(ctx context.Context, pub *PublishAction) (*Test, error) {
	didUpdate, err := tp.Tests.Update(pub.TestId, struct {
		PriceID     string `bson:"priceid" json:"priceid,omitempty"`
		IsPublished bool   `bson:"is_published" json:"is_published"`
	}{
		PriceID:     pub.PriceId,
		IsPublished: true,
	})
	if didUpdate {
		tp.InfoLog.Printf("Test %s was published\n", pub.TestId)
		return &Test{Id: pub.TestId}, nil
	} else {
		tp.ErrorLog.Printf("Test %s publication failed: %v\n", pub.TestId, err)
		return nil, err
	}
}

func (ts *Server) QualifyUserTestAndSession(ctx context.Context, pub *SessionAction) (*TestSession, error) {
	session, err := ts.TestSessions.Get(&models.KeySet{
		UserID: pub.UserId,
		TestID: pub.TestId,
	})

	if err != nil || session == nil {
		ts.ErrorLog.Println(err)
		return &TestSession{
			IsValid:     false,
			IsPermitted: true,
			Url:         pub.SessionUrl,
		}, nil
	}

	return &TestSession{
		IsValid:     session.IsValid,
		IsPermitted: session.Attempt < session.MaxAttempts,
		Url:         pub.SessionUrl,
	}, nil
}

func (ts *Server) InitiateTestSession(ctx context.Context, pub *SessionAction) (*TestSession, error) {
	test, err := ts.Tests.Get(pub.TestId)
	if err != nil {
		ts.ErrorLog.Println("Test reception failed: ", err)
		return nil, err
	}
	if test != nil {
		_, err = ts.TestSessions.Insert(&models.TestSession{
			UserID:      pub.UserId,
			TestID:      pub.TestId,
			SessionID:   pub.SessionId,
			SessionURL:  pub.SessionUrl,
			MaxAttempts: test.Attempts,
		})
		if err != nil {
			ts.ErrorLog.Println("Test session creation failed", err)
			return nil, err
		}
		return &TestSession{IsValid: true, IsPermitted: true, Url: pub.SessionUrl}, nil
	}
	return nil, errors.New("unable to initiate test session: test does not exist")
}
