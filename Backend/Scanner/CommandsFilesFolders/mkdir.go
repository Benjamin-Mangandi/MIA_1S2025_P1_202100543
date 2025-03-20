package CommandsFilesFolders

import (
	"Backend/Globals"
	"flag"
	"fmt"
	"strings"
)

func Mkdir(params string) {
	fs := flag.NewFlagSet("mkdir", flag.ContinueOnError)
	path := fs.String("path", "", "Ruta")
	p := fs.String("r", "false", "Creacion de carpetas padre")
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
	if *path == "" {
		fmt.Println("Error: el parametro '-path' es  obligatorio")
		return
	}
	fmt.Println(*p)
}
