package store

import (
	"github.com/asdine/storm/q"
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
		Every:     rawServer.Every,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		AddedBy:   userID,
	}

	err = s.Orm.Save(server)

	return
}

func (s *Store) AddStats(data *models.MinecraftStatus, serverID int) error {
	stats := &models.Stats{
		Time:          time.Now(),
		CurrentPlayer: data.PlayerInfo.Current,
		MaxPlayer:     data.PlayerInfo.Max,
		ServerID:      serverID,
	}

	err := s.Orm.Save(stats)

	return err
}

func (s *Store) GetMinecraftServerList() ([]models.Server, error) {
	var servers []models.Server
	err := s.Orm.All(&servers)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (s *Store) GetMinecraftServerListByUser(userID int) ([]models.Server, error) {
	var servers []models.Server
	err := s.Orm.Find("AddedBy", userID, &servers)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (s *Store) GetMinecraftServerForUserByID(userID, serverID int) ([]models.Server, error) {
	var server []models.Server
	err := s.Orm.Select(q.And(
		q.Eq("ID", serverID),
		q.Eq("AddedBy", userID),
	)).Limit(1).Find(&server)

	if err != nil {
		return nil, err
	}

	return server, nil
}

func (s *Store) GetMinecraftServerByName(name string) (models.Server, error) {
	var server models.Server

	err := s.Orm.One("Name", name, &server)
	if err != nil {
		return models.Server{}, err
	}

	return server, nil
}

func (s *Store) GetMinecraftServerByURL(URL string) (models.Server, error) {
	var server models.Server

	err := s.Orm.One("Url", URL, &server)
	if err != nil {
		return models.Server{}, err
	}

	return server, nil
}

func (s *Store) GetMinecraftServerByID(id int) (models.Server, error) {
	var server models.Server

	err := s.Orm.One("ID", id, &server)
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

	err = s.Orm.DeleteStruct(&srv)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteServerForUserByID(userID, serverID int) error {
	srv, err := s.GetMinecraftServerForUserByID(userID, serverID)
	if err != nil {
		return err
	}

	err = s.Orm.DeleteStruct(&srv)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateServer(userID, serverID int, newData models.Server) (models.Server, error) {
	srvs, err := s.GetMinecraftServerForUserByID(userID, serverID)
	if err != nil {
		return models.Server{}, err
	}

	srv := srvs[0]

	srv.Name = newData.Name
	srv.Timeout = newData.Timeout
	srv.Port = newData.Port
	srv.Every = newData.Every
	srv.Url = newData.Url

	err = s.Orm.Update(&srv)
	if err != nil {
		return models.Server{}, err
	}

	return srv, nil
}
