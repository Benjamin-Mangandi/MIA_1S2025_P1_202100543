package CommandsUsers

import (
	"Backend/Globals"
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
	if *user == "" || *group == "" {
		fmt.Println("Error: los parametros '-user' '-grp' son obligatorios")
		return
	}

	// Verificar que el nombre del usuario no sea "root"
	if strings.ToLower(*user) == "root" {
		fmt.Println("Error: No se puede mover el grupo del usuario 'root'")
		return
	}
}
