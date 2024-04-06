package main

import (
	"anncouncement/configs"
	"anncouncement/configs/logger"
	delivery_auth_grpc "anncouncement/services/authorization/delivery/grpc"
	delivery "anncouncement/services/authorization/delivery/http"
	"anncouncement/services/authorization/usecase"
	"os"
	"sync"
)

func main() {
	log := logger.GetLogger()
	err := configs.InitEnv()
	if err != nil {
		log.Errorf("Init env error: %s", err.Error())
		return
	}

	grpcCfg, err := configs.GetGrpcConfig()
	if err != nil {
		log.Errorf("failed to parse grpc configs file: %s", err.Error())
		return
	}

	psxCfg, err := configs.GetAuthPsxConfig()
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
	grpcServ, err := delivery_auth_grpc.NewServer(log)
	if err != nil {
		log.Errorf("Create grpc new server error: %s", err.Error())
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		port := os.Getenv("AUTH_APP_PORT")
		log.Infof("Server auth running on %s", port)

		err = api.ListenAndServe(port)
		if err != nil {
			log.Errorf("Listen and serve error: %s", err.Error())
			return
		}
	}()

	go func() {
		defer wg.Done()
		log.Infof("Server auth-grpc running on %s", grpcCfg.Port)

		err := grpcServ.ListenAndServeGrpc(grpcCfg)
		if err != nil {
			log.Errorf("Listen and serve grpc error: %s", err.Error())
			return
		}
	}()

	wg.Wait()
}
