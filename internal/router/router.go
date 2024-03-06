package router

import (
	"path/filepath"

	"github.com/SaYaku64/monopoly/internal/lobby"
	"github.com/SaYaku64/monopoly/internal/naming"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine

	lm *lobby.LobbyModule
}

func NewRouter(lm *lobby.LobbyModule) *Router {
	router := gin.Default()

	return &Router{
		engine: router,
		lm:     lm,
	}
}

func (r *Router) Load() {
	absPath, _ := filepath.Abs("../front")

	r.engine.LoadHTMLGlob(absPath + "/templates/*")

	r.engine.Use(static.Serve("/", static.LocalFile(absPath+"/static", false)))

	r.initializeRoutes()

}

func (r *Router) RunRouter() error {
	return r.engine.Run()
}

func (r *Router) initializeRoutes() {
	r.engine.Use()

	r.engine.GET("/", r.showIndexPage)

	apiV1 := r.engine.Group("/api/v1")
	apiV1.GET("/randomName", naming.GetRandName)
	apiV1.POST("/createLobby", r.CreateLobbyHandler)
	apiV1.GET("/getLobbiesTable", r.GetLobbiesTable)
	apiV1.GET("/removeLobby", r.RemoveLobby)

	// router.POST("/login", performLogin)
	// router.POST("/register", register)
	// router.POST("/change", changeTerritory)
	// router.GET("/logout", logout)
	// router.GET("/survey", ensureLoggedIn(), showSurveyPage)
	// router.POST("/survey", ensureLoggedIn(), surveyComplete)
	// router.GET("/infographics", ensureLoggedIn(), showInfographicsPage)
	// router.POST("/infographics", ensureLoggedIn(), infographicsShow)
	// router.GET("/clearInfo", ensureLoggedIn(), clearInfo)

	// router.GET("/object/:name", showObjectPage)
}
