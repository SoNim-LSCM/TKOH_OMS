package main

import (
	"github.com/SoNim-LSCM/TKOH_OMS/app"
	_ "github.com/SoNim-LSCM/TKOH_OMS/docs"
)

// @title			TKOH OMS
// @version		0.1
// @description	TKOH OMS Backend Server
// @contact.name	So Nim Wai
// @license.name	LSCM
// @host			ams.lscm.hk
// @BasePath		/oms/
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// setup and run app
	app.SetupAndRunApp()
}
