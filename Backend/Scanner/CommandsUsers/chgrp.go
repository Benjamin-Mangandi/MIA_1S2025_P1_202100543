package CommandsUsers

import (
	"Backend/Globals"
	"Backend/Responsehandler"
	"Backend/UsersManager"
	"flag"
	"fmt"
	"strings"
)

func Chgrp(params string) {
	fs := flag.NewFlagSet("Chgrp", flag.ExitOnError)
	user := fs.String("user", "", "Nombre del usuario")
	group := fs.String("grp", "", "Nombre del grupo")

	// Extraer parámetros usando Regex
	matches := Globals.Regex.FindAllStringSubmatch(params, -1)
	parsedFlags := make(map[string]string)

	for _, match := range matches {
		flagName := strings.ToLower(match[1])     // Nombre del flag
		flagValue := strings.Trim(match[2], "\"") // Quitar comillas

		// Asigna el flag en la estructura fs
		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Error: No se pudo establecer el flag '%s'\n", flagName)
			return
		}
		parsedFlags[flagName] = flagValue // Guardar para depuración
	}

	// Verificar que el parámetro -name se haya ingresado
	if *user == "" || *group == "" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: los parametros '-user' '-grp' son obligatorios"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Verificar que el nombre del usuario no sea "root"
	if strings.ToLower(*user) == "root" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No se puede mover el grupo del usuario 'root'"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	UsersManager.Chgrp(*user, *group)
}
