package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gedion-01/Go-Crud-Challenge/api"
	"github.com/Gedion-01/Go-Crud-Challenge/db"
	"github.com/Gedion-01/Go-Crud-Challenge/server"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	listenAddr := flag.String("listen-addr", ":5000", "the listen address of the api server")

	memoryDB := db.NewMemoryDB()
	srv := server.NewServer(memoryDB)
	personHandler := api.NewPersonHandler(memoryDB)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	apiv1 := app.Group("/api/v1")
	apiv1.Put("/person/:id", api.CheckUserExists(memoryDB), personHandler.HandlePutPerson)
	apiv1.Post("/person", personHandler.HandlePostPerson)
	apiv1.Get("/person/:id", personHandler.HandleGetPerson)
	apiv1.Get("/person", personHandler.HandleGetAllPersons)
	apiv1.Delete("/person/:id", api.CheckUserExists(memoryDB), personHandler.HandleDeletePerson)

	go func() {
		if err := app.Listen(*listenAddr); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()
	<-stop
	srv.Stop()
}
