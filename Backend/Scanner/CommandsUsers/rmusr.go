package CommandsUsers

import (
	"Backend/Globals"
	"Backend/UsersManager"
	"flag"
	"fmt"
	"strings"
)

func Rmusr(params string) {
	fs := flag.NewFlagSet("rmusr", flag.ExitOnError)
	name := fs.String("user", "", "Nombre del usuario")

	// Extraer parámetros usando Regex
	matches := Globals.Regex.FindAllStringSubmatch(params, -1)
	parsedFlags := make(map[string]string)

	for _, match := range matches {
		flagName := match[1]                      // Nombre del flag
		flagValue := strings.Trim(match[2], "\"") // Quitar comillas

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

	// Verificar que el parámetro -name se haya ingresado
	if *name == "" {
		fmt.Println("Error: el parámetro '-name' es obligatorio")
		return
	}

	// Verificar que el nombre del grupo no sea "root"
	if strings.ToLower(*name) == "root" {
		fmt.Println("Error: No se puede eliminar el grupo con el nombre 'root'")
		return
	}

	UsersManager.Rmusr(*name)
}
