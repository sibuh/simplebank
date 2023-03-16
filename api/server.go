package api

import (
	db "assignment_01/simplebank/db/sqlc"
	"assignment_01/simplebank/util"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server := &Server{store: store}
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	router.POST("/tranfers", server.createTransfer)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts/", server.listAccount)
	server.router = router
	return server
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
