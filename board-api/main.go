package main

import (
	"log"
	"github.com/DiGregory/testAvito/board-api/advert"
	_ "github.com/lib/pq"
)

type apiApp struct {
	Addr          string
	AdvertStorage *advert.AdStorage
}

func createApiApp(addr, advertStorageDSN string) (*apiApp, error) {
	advertStorage, err := advert.StorageConnect("postgres", advertStorageDSN)
	if err != nil {
		return nil, err
	}

	return &apiApp{Addr: addr, AdvertStorage: advertStorage}, nil
}

func main() {
	DSN := "host=192.168.99.100 user=postgres-dev password=1234 dbname=dev port=5432 sslmode=disable"

	myApp, err := createApiApp(":5000", DSN)
	if err != nil {
		log.Fatal("can`t run API application: ",err)
	}

	err = myApp.createApiHandlers()
	if err != nil {
		log.Fatal("can`t run API application: ",err)
	}
}
