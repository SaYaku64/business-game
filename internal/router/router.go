package router

import (
	"path/filepath"
	"sync"

	"github.com/SaYaku64/business-game/internal/game"
	"github.com/SaYaku64/business-game/internal/lobby"
	"github.com/SaYaku64/business-game/internal/naming"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
	games  map[string]*GameLobby
	gMux   sync.RWMutex

	lm *lobby.LobbyModule
	gm *game.GameModule
}

func NewRouter(lm *lobby.LobbyModule, gm *game.GameModule) *Router {
	router := gin.Default()

	return &Router{
		engine: router,
		games:  make(map[string]*GameLobby),
		lm:     lm,
		gm:     gm,
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
	r.engine.GET("/game", r.showGamePage)
	r.engine.GET("/ws", r.HandleWebSocket)
	r.engine.GET("/ws/game", r.HandleWSGame)

	apiV1 := r.engine.Group("/api/v1")
	apiV1.GET("/randomName", naming.GetRandName)
	apiV1.GET("/getSessionID", r.GetSessionID)
	apiV1.POST("/createLobby", r.CreateLobbyHandler)
	apiV1.GET("/getLobbiesTable", r.GetLobbiesTable)
	apiV1.GET("/removeLobby", r.RemoveLobby)
	apiV1.POST("/connectLobby", r.ConnectLobby)
	// apiV1.GET("/redirectToLobby", r.RedirectToLobby)
	apiV1.POST("/checkActiveGame", r.CheckActiveGame)
	apiV1.POST("/isLobbyExists", r.IsLobbyExists)

	apiGame := apiV1.Group("/game")
	apiGame.POST("/updatePlates", r.UpdatePlates)
	apiGame.GET("/turn", r.Turn)
	apiGame.GET("/buy", r.Buy)
	apiGame.GET("/payRent", r.PayRent)

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
