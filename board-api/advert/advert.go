package advert

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/lib/pq"
	"errors"
	"time"
)

var InvalidFieldErr = errors.New("invalid input field(s) length")

type Advert struct {
	ID          int        `json:"-"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Images      []string   `json:"images"`
	Price       float64    `json:"price"`
	CreatedAt   *time.Time `json:"-"`
	UpdatedAt   *time.Time `json:"-"`
	DeletedAt   *time.Time `json:"-"`
}

type AdStorage struct {
	DB *sql.DB
}

func StorageConnect(driver, source string) (*AdStorage, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}

	initSQLFile, err := os.Open("init.sql")
	if err != nil {
		return nil, err
	}
	defer initSQLFile.Close()

	initQuery, err := ioutil.ReadAll(initSQLFile)
	if err != nil {
		return nil, err
	}

	_, err = db.Query(string(initQuery))
	if err != nil {
		return nil, err
	}

	return &AdStorage{DB: db}, nil
}

func (as AdStorage) CreateAdvert(jsonData []byte) (id *int, err error) {
	var a Advert
	err = json.Unmarshal(jsonData, &a)
	if err != nil {
		return nil, err
	}

	badFields := len(a.Images) > 3 || len(a.Description) > 1000 ||
		len(a.Name) > 200 || len(a.Images) < 1 ||
		(a.Description == "") || (a.Name == "") || len(a.Images) == 0 ||
		(a.Price) == 0.0
	if badFields {
		return nil, InvalidFieldErr
	}

	adRow := as.DB.QueryRow("INSERT INTO adverts"+
		"(id,name,description,images,price) "+
		"VALUES (DEFAULT ,$1,$2,$3,$4) RETURNING id",
		a.Name, a.Description, pq.Array(a.Images), a.Price)

	var ad Advert
	err = adRow.Scan(&ad.ID)
	if err != nil {
		return nil, err
	}

	return &ad.ID, nil

}

func (as AdStorage) GetSingleAdvert(adID int, fields []string) (*Advert, error) {
	ad := new(Advert)
	ad.Images = make([]string, 1, 3)

	err := as.DB.QueryRow("SELECT adverts.name,adverts.images[1],adverts.price "+
		"FROM adverts WHERE id=$1", adID).
		Scan(&ad.Name, &ad.Images[0], &ad.Price)
	if err != nil {
		return nil, err
	}
	for _, v := range fields {
		switch {
		case v == "images":
			err := as.DB.QueryRow("SELECT adverts.images"+
				" FROM adverts WHERE id=$1", adID).
				Scan(pq.Array(&ad.Images))
			if err != nil {
				return nil, err
			}

		case v == "description":
			err := as.DB.QueryRow("SELECT adverts.description"+
				" FROM adverts WHERE id=$1", adID).
				Scan(&ad.Description)
			if err != nil {
				return nil, err
			}
		}
	}
	return ad, nil
}

func (as AdStorage) GetAdverts(offset int, sortBy, orderBy string) ([]*Advert, error) {
	var err error
	var adsRecords *sql.Rows
	ads := make([]*Advert, 0)
	//по-умолчанию сортировка id asc;offset=0
	sortByField := "id"
	if sortBy == "created_at" || sortBy == "price" {
		sortByField = sortBy
	}
	if orderBy != "desc" {
		orderBy = "asc"
	}

	switch orderBy {
	case "desc":
		adsRecords, err = as.DB.Query("SELECT adverts.name,adverts.images[1],adverts.price"+
			" FROM adverts ORDER BY $1 DESC LIMIT 10 OFFSET $2;", sortByField, offset)
	default:
		adsRecords, err = as.DB.Query("SELECT adverts.name,adverts.images[1],adverts.price"+
			" FROM adverts ORDER BY $1 ASC LIMIT 10 OFFSET $2;", sortByField, offset)
	}

	if err != nil {
		return nil, err
	}
	defer adsRecords.Close()

	for adsRecords.Next() {
		singleAd := new(Advert)
		singleAd.Images = make([]string, 1)
		err := adsRecords.Scan(&singleAd.Name, &singleAd.Images[0], &singleAd.Price)
		if err != nil {
			return nil, err
		}
		ads = append(ads, singleAd)

	}

	return ads, nil
}
