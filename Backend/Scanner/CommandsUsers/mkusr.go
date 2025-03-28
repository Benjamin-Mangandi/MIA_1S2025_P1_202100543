package CommandsUsers

import (
	"Backend/Globals"
	"Backend/Responsehandler"
	"Backend/UsersManager"
	"flag"
	"fmt"
	"strings"
)

func Mkusr(params string) {
	fs := flag.NewFlagSet("mkusr", flag.ExitOnError) // Corrección del nombre
	user := fs.String("user", "", "Nombre del usuario")
	pass := fs.String("pass", "", "Contraseña del usuario")
	group := fs.String("grp", "", "Grupo del usuario")

	// Extraer parámetros usando Regex
	matches := Globals.Regex.FindAllStringSubmatch(params, -1)
	parsedFlags := make(map[string]string)

	for _, match := range matches {
		flagName := strings.ToLower(match[1])     // Nombre del flag
		flagValue := strings.Trim(match[2], `"'`) // Quitar comillas simples y dobles

		// Asigna el flag en la estructura fs
		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Error: No se pudo establecer el flag '%s'\n", flagName)
			return
		}
		parsedFlags[flagName] = flagValue // Guardar para depuración
	}

	// Verificar que los parámetros sean válidos
	if *user == "" || *pass == "" || *group == "" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: Los parámetros '-user', '-pass' y '-group' son obligatorios"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Validar que ningún parámetro tenga más de 10 caracteres
	if len(*user) > 10 || len(*pass) > 10 || len(*group) > 10 {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: Los parámetros '-user', '-pass' y '-group' no pueden tener más de 10 caracteres."
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Verificar que el nombre del usuario no sea "root"
	if strings.ToLower(*user) == "root" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No se puede crear un usuario con el nombre 'root'"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	UsersManager.Mkusr(*user, *pass, *group)
}
