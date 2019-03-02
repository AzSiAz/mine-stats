package store

import (
	"mine-stats/models"
	minecraftProtocol "mine-stats/protocol/minecraft"
	"time"
)

func (s *Store) AddServer(rawServer *minecraftProtocol.MinecraftServer, userID int) (server *models.Server, err error) {
	server = &models.Server{
		Name:      rawServer.Name,
		Port:      rawServer.Port,
		Url:       rawServer.Address,
		Timeout:   rawServer.Timeout,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		AddedBy: userID,

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

func (s *Store) GetMinecraftServer() ([]models.Server, error) {
	var servers []models.Server
	err := s.orm.All(&servers)
	if err != nil {
		return nil, err
	}

	return servers, nil
}