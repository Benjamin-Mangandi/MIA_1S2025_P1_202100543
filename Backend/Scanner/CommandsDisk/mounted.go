package CommandsDisk

import (
	"Backend/DiskManager"
	"Backend/Responsehandler"
	"strings"
)

func Mounted(params string) {
	// Verificar que no se pasen parámetros
	if strings.TrimSpace(params) != "" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: El comando 'mounted' no acepta parámetros"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	// Llamar a la función para imprimir las particiones montadas
	DiskManager.PrintMountedPartitions()
}
