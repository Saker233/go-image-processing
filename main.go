package main

import (
	"log"

	"github.com/Saker233/go-image-processing/internal/api"
	db "github.com/Saker233/go-image-processing/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	godotenv.Load("app.env")

	conn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	store := db.NewStore(conn)
	api.SetupServer(store)

}
