package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"coremeridian.xyz/app/qlist/cmd/api/services/tests/models/mongodb"
	"coremeridian.xyz/app/qlist/cmd/api/services/tests/service"
	"coremeridian.xyz/app/qlist/cmd/api/services/tests/service/grpc/testpaymentintegrator"
	"github.com/costal/go-misc-tools/httpapp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	var databaseName string
	addrn := flag.String("addr", "", "address accessing server")
	port := flag.Int("port", 1900, "address port")
	mode := flag.String("mode", "development", "deployment mode")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *addrn, *port)
	isProduction := strings.ToLower(*mode) == "production"
	isStaging := strings.ToLower(*mode) == "staging"
	notDevelopment := isStaging || isProduction

	if isProduction {
		databaseName = os.Getenv("TESTSDB_PROD")
	} else {
		databaseName = os.Getenv("TESTSDB_DEV")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	database, err := mongodb.GetDB(databaseName)
	if err != nil {
		panic(err)
	}

	client := database.Client()
	defer client.Disconnect(ctx)

	endpoints := service.Endpoints{
		Tests:        &mongodb.TestModel{DB: database},
		TestSessions: &mongodb.TestSessionModel{DB: database},
	}

	app := &service.HTTPApplication{
		HTTP:      httpapp.DefaultApp("Psychometrics (Tests) API", os.Getenv("MAINAPP_URL")),
		Endpoints: endpoints,
	}

	app.HTTP.Auth = httpapp.Auth{
		Realm:        os.Getenv("REALM"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURI:  os.Getenv("REDIRECT_URI"),
		GrantType:    "authorization_code",
	}

	app.Init()

	server := http.Server{
		Addr:         addr,
		ErrorLog:     app.HTTP.ErrorLog(),
		Handler:      service.Router(app),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	var wg sync.WaitGroup
	listener, closeListener := createListener(addr)

	wg.Add(2)
	mainServerChan := func() <-chan struct{} {
		mainServerChan := make(chan struct{})
		go func() {
			serverType := map[bool]string{true: "HTTPS", false: "HTTP"}[notDevelopment]
			/*
			   (declaration of a non-const map, calculation of a parameter, map lookup of the parameter, assignment) and the insane number of things you are doing implicitly:
			   - inline construction of a map object, - inline population of the map object (hashing true and false, allocating two nodes for them, assigning that to the map),
			   - production of bool key, - hashing of bool key, - binary search of map, - release of 2 map nodes, - release of map, assignment
			*/
			defer wg.Done()
			defer closeListener()
			defer func() {
				app.HTTP.InfoLog().Printf("Closing %s \"%s\" server on port %d", serverType, app.HTTP.Name(), listener.Addr().(*net.TCPAddr).Port)
			}()

			app.HTTP.InfoLog().Printf("Starting %s \"%s\" server on port %d", serverType, app.HTTP.Name(), listener.Addr().(*net.TCPAddr).Port)

			var serverError error
			go func() {
				if notDevelopment {
					serverError = server.ServeTLS(listener, os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE"))
				} else {
					serverError = server.Serve(listener)
				}
			}()
			for serverError == nil {
				mainServerChan <- struct{}{}
			}
			app.HTTP.ErrorLog().Fatalf("(%s) %v\n", serverType, serverError)
		}()
		return mainServerChan
	}()

	<-mainServerChan

	go func() {
		port := listener.Addr().(*net.TCPAddr).Port + 10
		listener, closeListener := createListener(fmt.Sprintf("%s:%d", *addrn, port))
		name := "Test Publisher (Payment-Sub) API"
		app.HTTP.InfoLog().Printf("Starting %s \"%s\" server on port %d", "gRPC", name, port)
		defer wg.Done()
		defer closeListener()
		defer func() {
			app.HTTP.InfoLog().Printf("Closing %s \"%s\" server on port %d", "gRPC", name, port)
		}()

		cred, err := credentials.NewServerTLSFromFile(
			os.Getenv("CERT_FILE"),
			os.Getenv("KEY_FILE"),
		)
		if err != nil {
			app.HTTP.ErrorLog().Panic(err)
		}
		opts := []grpc.ServerOption{grpc.Creds(cred)}
		grpcServer := grpc.NewServer(opts...)
		testpaymentintegrator.RegisterTestPaymentIntegratorServer(
			grpcServer,
			&testpaymentintegrator.Server{
				Endpoints: endpoints,
				ErrorLog:  app.HTTP.ErrorLog(),
				InfoLog:   app.HTTP.InfoLog(),
			})
		grpcServer.Serve(listener)
	}()
	wg.Wait()
}

func createListener(addr string) (l net.Listener, close func()) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	return l, func() { _ = l.Close() }
}
