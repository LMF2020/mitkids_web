package job

import (
	"github.com/robfig/cron/v3"
	"mitkid_web/conf"
	"mitkid_web/service"
	"mitkid_web/utils/log"
	"time"
)

var s *service.Service

func Init(conf *conf.Config, service *service.Service) {
	s = service
	cron := cron.New()
	cron.Start()
	defer cron.Stop()

	cron.AddFunc(conf.Job.EndClassOccurrClassOccurrencesCron, endClassOccurrClassOccurrencesJob)
}

func endClassOccurrClassOccurrencesJob() {
	log.Logger.Info("job run endClassOccurrClassOccurrencesJob")
	time := time.Now()
	s.EndClassOccurrClassOccurrencesByDateTimeSql(&time)

}
