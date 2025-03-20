package main

import (
	"Backend/Responsehandler"
	"Backend/Scanner"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/scannear", func(ctx *gin.Context) {
		Responsehandler.Clear()
		var requestData struct {
			Input string `json:"input"`
		}

		if err := ctx.BindJSON(&requestData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato JSON inválido"})
			return
		}
		Scanner.Scan(requestData.Input) // Pasar el dato a la función Scanner.Scan
		ctx.JSON(200, gin.H{"consola": Responsehandler.GlobalResponse.Content})
	})
	// Iniciar el servidor en el puerto 7777
	router.Run(":7777")
}
