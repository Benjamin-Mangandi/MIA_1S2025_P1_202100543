package CommandsUsers

import (
	"Backend/Globals"
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
		flagName := match[1]                      // Nombre del flag
		flagValue := strings.Trim(match[2], `"'`) // Quitar comillas simples y dobles

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

	// Verificar que los parámetros sean válidos
	if *user == "" {
		fmt.Println("Error: el parámetro '-user' es obligatorio")
		return
	}
	if *pass == "" {
		fmt.Println("Error: el parámetro '-pass' es obligatorio")
		return
	}
	if *group == "" {
		fmt.Println("Error: el parámetro '-grp' es obligatorio")
		return
	}

	// Verificar que el nombre del usuario no sea "root"
	if strings.ToLower(*user) == "root" {
		fmt.Println("Error: No se puede crear un usuario con el nombre 'root'")
		return
	}

	UsersManager.Mkusr(*user, *pass, *group)
}
