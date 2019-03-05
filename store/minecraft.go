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
		AddedBy:   userID,
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

func (s *Store) GetMinecraftServerList() ([]models.Server, error) {
	var servers []models.Server
	err := s.orm.All(&servers)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (s *Store) GetMinecraftServerByName(name string) (models.Server, error) {
	var server models.Server

	err := s.orm.One("Name", name, &server)
	if err != nil {
		return models.Server{}, err
	}

	return server, nil
}

func (s *Store) GetMinecraftServerByURL(URL string) (models.Server, error) {
	var server models.Server

	err := s.orm.One("Url", URL, &server)
	if err != nil {
		return models.Server{}, err
	}

	return server, nil
}

func (s *Store) GetMinecraftServerByID(id int) (models.Server, error) {
	var server models.Server

	err := s.orm.One("ID", id, &server)
	if err != nil {
		return models.Server{}, err
	}

	return server, nil
}

func (s *Store) DeleteServerByID(id int) error {
	srv, err := s.GetMinecraftServerByID(id)
	if err != nil {
		return err
	}

	err = s.orm.DeleteStruct(&srv)
	if err != nil {
		return err
	}

	return nil
}
