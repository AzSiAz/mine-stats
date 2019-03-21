package jobs

import (
	"github.com/sirupsen/logrus"
	"mine-stats/models"
	"mine-stats/protocol/minecraft"
	"mine-stats/store"
	"strings"
	"time"
)

type Runnable interface {
	Run()
}

var jobList []*Job

func AddJob(j *Job) {
	jobList = append(jobList, j)
	logrus.WithFields(logrus.Fields{
		"server": j.Server.Name,
		"url": j.Server.Url,
	}).Infoln("Launching new job")
	go j.Loop()
}

func ShutDownJob() {
	for _, job := range jobList {
		close(job.quit)
	}
}

type Job struct {
	ticker *time.Ticker
	quit   chan struct{}
	Server *models.Server
}

func NewJob(server *models.Server) *Job {
	ticker := time.NewTicker(server.Every * time.Second)
	quit := make(chan struct{}, 1)

	return &Job{
		Server: server,
		quit:   quit,
		ticker: ticker,
	}
}

func (j *Job) Loop() {
	for {
		j.Run()

		select {
		case <-j.ticker.C:
		case <-j.quit:
			j.ticker.Stop()
			return
		}
	}
}

func (j *Job) Run() {
	srv := minecraftProtocol.NewMinecraftServer(j.Server.Name, j.Server.Url, j.Server.Port, j.Server.Timeout, j.Server.Every)
	status, err := srv.Query()
	if err != nil {
		logrus.WithError(err).
			WithFields(logrus.Fields{
				"server_name": j.Server.Name,
				"url":         j.Server.Url,
			}).Infoln()
		return
	}
	go j.LogRun(status)
	err = store.GetStore().AddStats(status, j.Server.ID)
	if err != nil {
		logrus.WithError(err).Infoln("Error adding stats to db")
	}
}

func (j *Job) LogRun(status *models.MinecraftStatus) {
	var playerList []string

	for _, player := range status.PlayerInfo.Players {
		playerList = append(playerList, player.Name)
	}

	logrus.WithFields(logrus.Fields{
		"name": j.Server.Name,
		"max": status.PlayerInfo.Max,
		"online": status.PlayerInfo.Current,
		"players": strings.Join(playerList, ","),
	}).Infoln()
}
