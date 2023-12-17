package server

import (
	"net/http"
	"strconv"

	"github.com/bitgemtech/exotic-indexer/exotic"
	badger "github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
)

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

type UTXORangesRequest struct {
	UTXOs               []string `json:"utxos" binding:"required"`
	ExcludeCommonRanges bool     `json:"excludeCommonRanges"`
}

type RangeResponse struct {
	Utxo   string `json:"utxo"`
	Start  int64  `json:"start"`
	Size   int64  `json:"size"`
	End    int64  `json:"end"`
	Offset int64  `json:"offset"`
}

type ExoticRangeResponse struct {
	RangeResponse
	Satributes []exotic.Satribute `json:"satributes"`
}

func (s *Server) getUTXORanges(c *gin.Context) {
	var req UTXORangesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := s.model.GetUTXORanges(req.UTXOs, req.ExcludeCommonRanges)
	if err == badger.ErrKeyNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

type AddressRangeRequest struct {
	Address             string `json:"address" binding:"required"`
	ExcludeCommonRanges bool   `json:"excludeCommonRanges"`
}

func (s *Server) getAddressRanges(c *gin.Context) {
	var req AddressRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := s.model.GetAddressRanges(req.Address, req.ExcludeCommonRanges)
	if err == badger.ErrKeyNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) GetSat(c *gin.Context) {
	satStr := c.Param("sat")
	sat, err := strconv.ParseInt(satStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	res := s.model.GetSat(sat)

	c.JSON(http.StatusOK, res)
}

func (s *Server) GetSatributes(c *gin.Context) {
	c.JSON(http.StatusOK, exotic.SatributeList)
}
