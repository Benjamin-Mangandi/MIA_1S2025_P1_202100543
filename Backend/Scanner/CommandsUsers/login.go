package CommandsUsers

import (
	"Backend/FileSystem"
	"Backend/Globals"
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

	// Imprimir parámetros detectados para depuración
	fmt.Println("====== Parámetros Escaneados ======")
	for key, value := range parsedFlags {
		fmt.Printf("%s: %s\n", key, value)
	}
	fmt.Println("===================================")

	// Validación de campos obligatorios
	if *user == "" || *pass == "" || *id == "" {
		fmt.Println("Error: Los parámetros '-user', '-pass' y '-id' son obligatorios")
		return
	}
	FileSystem.Login(*user, *pass, *id)
}
