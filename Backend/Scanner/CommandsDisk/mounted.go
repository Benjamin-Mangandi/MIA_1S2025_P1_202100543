package CommandsDisk

import (
	"Backend/DiskManager"
	"fmt"
	"strings"
)

func Mounted(params string) {
	// Verificar que no se pasen parámetros
	if strings.TrimSpace(params) != "" {
		fmt.Println("Error: El comando 'mounted' no acepta parámetros")
		return
	}
	// Llamar a la función para imprimir las particiones montadas
	DiskManager.PrintMountedPartitions()
}
