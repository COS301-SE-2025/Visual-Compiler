package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Lexing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Lexing works"})
}
