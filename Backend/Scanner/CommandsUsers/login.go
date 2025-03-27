package CommandsUsers

import (
	"Backend/Globals"
	"Backend/Responsehandler"
	"Backend/UsersManager"
	"flag"
	"fmt"
	"strings"
)

func Login(params string) {
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	user := fs.String("user", "", "Usuario")
	pass := fs.String("pass", "", "Contraseña")
	id := fs.String("id", "", "ID Particion")

	matches := Globals.Regex.FindAllStringSubmatch(params, -1)

	// Mapa para almacenar los valores ingresados por el usuario
	parsedFlags := make(map[string]string)

	for _, match := range matches {
		flagName := match[1]                      // Flag tal cual fue escrito
		flagValue := strings.Trim(match[2], "\"") // Quita comillas si las tiene

		// Asigna el flag en la estructura fs
		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Error: No se pudo establecer el flag '%s'\n", flagName)
			return
		}
		parsedFlags[flagName] = flagValue // Guardar para depuración
	}

	// Validación de campos obligatorios
	if *user == "" || *pass == "" || *id == "" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: Los parámetros '-user', '-pass' y '-id' son obligatorios"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	UsersManager.Login(*user, *pass, *id)
}
