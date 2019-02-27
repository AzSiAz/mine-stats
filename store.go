package main

import (
	"github.com/asdine/storm"
	"mine-stats/models"
	"mine-stats/protocol/minecraft"
	"time"
)

type Store struct {
	orm *storm.DB
}

func initStorm(db *storm.DB) (err error) {
	//err = db.Init(&models.ServerTypes{})
	err = db.Init(&models.Server{})
	err = db.Init(&models.Stats{})

	return
}

func NewStore(path string) (store *Store, err error) {
	db, err := storm.Open(path)
	err = initStorm(db)

	store = &Store{orm: db}

	return
}

func (s *Store) Close() (err error) {
	err = s.orm.Close()

	return
}

func (s *Store) AddServer(rawServer *minecraftProtocol.MinecraftServer) (server *models.Server, err error) {
	server = &models.Server{
		Name:      rawServer.Name,
		Port:      rawServer.Port,
		Url:       rawServer.Address,
		Timeout:   rawServer.Timeout,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.orm.Save(server)

	return
}

func (s *Store) AddStats(data *models.MinecraftStatus, serverID uint) (stats *models.Stats, err error) {
	stats = &models.Stats{
		Time:          time.Now(),
		CurrentPlayer: data.PlayerInfo.Current,
		MaxPlayer:     data.PlayerInfo.Max,
		ServerID:      serverID,
	}

	err = s.orm.Save(stats)

	return
}
