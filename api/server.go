package api

import (
	db "assignment_01/simplebank/db/sqlc"
	"assignment_01/simplebank/token"
	"assignment_01/simplebank/util"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
)

type Server struct {
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()

	return server, nil
}
func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/accounts", server.createAccount)
	router.POST("/tranfers", server.createTransfer)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts/", server.listAccount)
	server.router = router
}
func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
func (server *Server) Start(addres string) error {
	config, err := util.LoadConfig("../.")
	if err != nil {
		log.Fatal("can not load configration:", err)
	}
	return server.router.Run(config.Addres)
}
