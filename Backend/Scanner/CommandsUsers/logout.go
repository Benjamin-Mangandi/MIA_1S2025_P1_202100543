package CommandsUsers

import (
	"Backend/FileSystem"
	"fmt"
	"strings"
)

func LogOut(params string) {
	// Verificar que no se pasen parámetros
	if strings.TrimSpace(params) != "" {
		fmt.Println("Error: El comando 'mounted' no acepta parámetros")
		return
	}
	// Llamar a la función cerrar sesión
	FileSystem.Logout()
}
