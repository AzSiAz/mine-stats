package store

import (
	"github.com/asdine/storm"
	"github.com/sirupsen/logrus"
	"mine-stats/models"
	"os"
)

type Store struct {
	AdminAdded bool
	FirstAdmin bool
	Orm        *storm.DB
}

var store *Store
var initDone = false

func GetStore() *Store {
	if initDone {
		return store
	} else {
		panic("Trying to access store without init first")
	}
}

func NewStore(path string, firstAdmin bool) (*Store, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, err
	}

	store = &Store{Orm: db, FirstAdmin: firstAdmin}

	err = store.initStorm(db)
	if err != nil {
		return nil, err
	}

	store.checkAlreadyAddedAdmin()

	initDone = true

	return store, nil
}

func (s *Store) initStorm(db *storm.DB) (err error) {
	//err = db.Init(&models.ServerTypes{})
	err = s.Orm.Init(&models.Server{})
	err = s.Orm.Init(&models.Stats{})
	err = s.Orm.Init(&models.User{})

	return
}

func (s *Store) checkAlreadyAddedAdmin() {
	var user models.User
	err := s.Orm.One("Role", models.AdminRole, &user)
	if err != nil {
		if err == storm.ErrNotFound {
			s.AdminAdded = false
			return
		}
		logrus.Fatalln("Can't check if there is already an admin, try again or remove first_admin flag")
		os.Exit(1)
	}
	s.AdminAdded = true
}

func (s *Store) Close() (err error) {
	err = s.Orm.Close()

	return
}
