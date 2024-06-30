package main

import (
	"log"

	"github.com/altsaqif/go-rest/cmd/delivery"
)

// @Golang API
// @version 1.0
// @description This is a API use Golang with GIN Framework.
// @termsOfService https://example.com/terms/
// @contact.name API Support
// @contact.email altsaqifnugraha19@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /api/v1

func main() {
	log.Println("Starting REST API with GIN")
	srv := delivery.NewServer()
	srv.Run()
}
