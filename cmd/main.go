package main

import (
	"log"

	"github.com/safe-homie/backend/cmd/app"
)

//	@title			SmartHomie API
//	@version		1.0
//	@description	This is API documentation for SmartHomie.
//	@termsOfService	http://swagger.io/terms/

//	@license.name	MIT
//	@license.url	http://opensource.org/licenses/MIT

//	@BasePath	/api/v1

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
