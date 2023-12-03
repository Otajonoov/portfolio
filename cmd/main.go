package main

import (
	"log"
	"net/http"
	_ "portfolio/internal/delivery/rest/docs"
	"portfolio/internal/delivery/rest/router"
	"portfolio/internal/pkg/config"
	"portfolio/internal/pkg/db"
	"portfolio/internal/service"
)

func main() {
	cfg := config.Load(".")
	connDB := db.ConnectToDatabase()

	s := service.NewUserService(connDB)
	apiServer := router.New(router.RouterOption{Service: &s})

	server := &http.Server{
		Addr:    ":" + cfg.HttpPort,
		Handler: apiServer,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	select {}
}


// if err := apiServer.Run(fmt.Sprintf(":%s", cfg.HttpPort)); err != nil {
// 	log.Fatalf("failed to run server: %v", err)
// }