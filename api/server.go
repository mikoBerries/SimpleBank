package api

import (
	"fmt"
	"net/http"
	"time"

	cv "github.com/MikoBerries/SimpleBank/api/costumValidator"
	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/token"
	"github.com/MikoBerries/SimpleBank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type server struct {
	store  db.Store
	token  token.Maker
	config util.Config
	router *gin.Engine
}

//NewServer Create new Server with gin router
func NewServer(config util.Config, store db.Store) (*server, error) {
	server := &server{}
	//sign db transaction logic (sqlc)
	server.store = store
	//sign config file connection (viper)
	server.config = config
	//crete tokeMaker and sign it to server
	tokenMaker, err := token.NewPasetoMaker(server.config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	//sign token maker (Paseto / JWT)
	server.token = tokenMaker

	//Register our costum validator to default gin (validator/v10)
	err = setCostumeBindingValidator()
	if err != nil {
		return nil, err
	}
	//sign router (Gin)
	server.setGinRouter()
	return server, nil
}

func setCostumeBindingValidator() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//adding own costum validator ("tag-name", func(validator.FieldLevel)bool )
		v.RegisterValidation("bookabledate", cv.BookableDate)
		v.RegisterValidation("IsCurrency", cv.IsCurrency)
		return nil
	}
	return fmt.Errorf("error seting binnding gin")
}

func (server *server) setGinRouter() {
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
	router.POST("/createUser", server.createUser)

	router.POST("/user/login", server.userLogin)

	//goruping to use middleWare
	authRouter := router.Group("/", authMiddleWare(server.token))

	authRouter.POST("/createAccount", server.createAccount)
	authRouter.GET("/account/:id", server.getAccountByID)
	authRouter.GET("/account", server.getListAccount)
	authRouter.POST("/transfer", server.createTransfer)

	server.router = router
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
