package main

import (
	"filmoteka/configs"
	"filmoteka/configs/logger"
	delivery "filmoteka/delivery/http"
	"filmoteka/usecase"
	_ "github.com/swaggo/swag"
)

func main() {
	log := logger.GetLogger()
	err := configs.InitEnv()
	if err != nil {
		log.Errorf("Init env error: %s", err.Error())
		return
	}

	psxCfg, err := configs.GetPsxConfig("configs/db_psx.yaml")
	if err != nil {
		log.Errorf("Create psx config error: %s", err.Error())
		return
	}

	redisCfg, err := configs.GetRedisConfig()
	if err != nil {
		log.Errorf("Create redis config error: %s", err.Error())
		return
	}

	core, err := usecase.GetCore(psxCfg, redisCfg, log)
	if err != nil {
		log.Errorf("Create core error: %s", err.Error())
		return
	}

	api := delivery.GetApi(core, log)

	log.Info("Server running")
	err = api.ListenAndServe("8081")
	if err != nil {
		log.Errorf("Listen and serve error: %s", err.Error())
		return
	}

}
