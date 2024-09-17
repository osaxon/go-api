package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"go-api/internal/server"
	"log"
	"os"
	"strconv"
)

func main() {
	//config.GenDal()
	run()
}

func run() {
	s := server.New()
	s.RegisterFiberRoutes()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := s.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
}
