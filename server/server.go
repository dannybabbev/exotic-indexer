package server

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	model *ServerModel
}

func NewServer(model *ServerModel) *Server {
	return &Server{
		model: model,
	}
}

func (s *Server) Start() {
	r := gin.Default()

	r.GET("/", s.healthHandler)
	r.HEAD("/", s.healthHandler)

	r.GET("/health", s.healthHandler)
	r.POST("/utxo-ranges", s.getUTXORanges)
	r.GET("/sat/:sat", s.GetSat)
	r.POST("/address-ranges", s.getAddressRanges)

	// Info section
	r.GET("/info/satributes", s.GetSatributes)

	r.Run()
}
