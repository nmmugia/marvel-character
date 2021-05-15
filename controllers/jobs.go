package controllers

import (
	"log"
	"sync"
	"time"

	"github.com/nmmugia/marvel/service"
	"github.com/robfig/cron"
)

func ListJobs(characterUsecase service.CharacterUsecase) {
	getAllDataByCacheHourly("@every 1h", characterUsecase)
}

func getAllDataByCacheHourly(schedule string, characterUsecase service.CharacterUsecase) {
	c := cron.New()
	var mutex = &sync.Mutex{}
	c.AddFunc(schedule, func() {
		log.Printf("[jobs][setCacheHourlyCharacterData] Job started at %+v", time.Now())
		mutex.Lock()
		defer mutex.Unlock()
		errx := characterUsecase.GetAllDataByCache()
		if errx.Err != nil {
			log.Printf("[jobs][setCacheHourlyCharacterData] error while GetAllDataByCache, cause: %v. Message:%s", errx.Err, errx.Message)
		}
		log.Printf("[jobs][setCacheHourlyCharacterData] Job ended at %+v", time.Now())
	})
	c.Start()
}
