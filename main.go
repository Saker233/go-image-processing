package main

import (
	"log"

	"github.com/Saker233/go-image-processing/internal/api"
	db "github.com/Saker233/go-image-processing/internal/database"
	"github.com/Saker233/go-image-processing/internal/util"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	godotenv.Load("app.env")
	util.InitS3()
	// key, err := util.UploadS3()
	// log.Println(key)
	// err := util.DownloadS3("images/682f004a-b5ca-4875-a072-3b693450e69e.jpg")
	// if err !=  nil {
	// 	log.Println(err)
	// }
	conn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	store := db.NewStore(conn)
	api.SetupServer(store)

}
