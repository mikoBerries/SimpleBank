package api

import (
	"net/http"
	"time"

	cv "github.com/MikoBerries/SimpleBank/api/costumValidator"
	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type server struct {
	store  db.Store
	router *gin.Engine
}

//NewServer Create new Server with gin router
func NewServer(store db.Store) *server {
	s := &server{}
	s.store = store

	// Default router With the Logger and Recovery middleware already attached
	router := gin.Default()

	//set logger mode
	// gin.SetMode(gin.TestMode)
	// gin.SetMode(gin.ReleaseMode)

	//Costume Debuger
	// gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	// 	log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	//   }
	// router.Group("/").Use(AuthRequired())

	//similar with handle func
	router.POST("/createAccount", s.createAccount)
	router.GET("/account/:id", s.getAccountByID)
	router.GET("/account", s.getListAccount)

	router.POST("/transfer", s.createTransfer)

	//register our costum validator to default gin (validator/v10)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//adding own costum validator ("tag-name", func(validator.FieldLevel)bool )
		v.RegisterValidation("bookabledate", cv.BookableDate)
		v.RegisterValidation("IsCurrency", cv.IsCurrency)
	}

	s.router = router
	return s
}

//StartServerAddress start server with given adress and some configuration
func (server *server) StartServerAddress(addrs string) error {
	s := &http.Server{
		Addr:           addrs,
		Handler:        server.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return s.ListenAndServe()
}
