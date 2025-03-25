package CommandsDisk

import (
	"Backend/DiskManager"
	"Backend/Globals"
	"Backend/Responsehandler"
	"flag"
	"fmt"
	"strings"
)

func Mount(params string) {
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	path := fs.String("path", "", "Ruta del disco")
	name := fs.String("name", "", "Nombre de la partición")

	// Extraer los parámetros usando regex
	matches := Globals.Regex.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := strings.ToLower(match[1])     // Convertir nombre del flag a minúsculas
		flagValue := strings.ToLower(match[2])    // Convertir valor a minúsculas si aplica
		flagValue = strings.Trim(flagValue, "\"") // Eliminar comillas
		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Error: No se pudo establecer el flag '%s'\n", flagName)
			return
		}
	}

	// Validación de campos obligatorios
	if *path == "" || *name == "" {
		response := "---------------------\n" +
			"Error: Los parámetros '-path' y '-name' son obligatorios"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Llamar a la función Mount
	DiskManager.Mount(*path, *name)

}
