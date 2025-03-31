package CommandsFileSystem

import (
	"Backend/FileSystem"
	"Backend/Globals"
	"flag"
	"fmt"
	"strings"
)

func Mkdir(params string) {
	fs := flag.NewFlagSet("mkdir", flag.ContinueOnError)
	path := fs.String("path", "", "Ruta")
	p := fs.String("p", "false", "Creacion de carpetas padre")

	// Extraer parámetros usando Regex
	matches := Globals.Regex.FindAllStringSubmatch(params, -1)
	parsedFlags := make(map[string]string)

	pFlagPresent := false // Variable para rastrear si -p está presente sin valor

	for _, match := range matches {
		flagName := strings.ToLower(match[1])     // fileN (ej. file1, file2, file3)
		flagValue := strings.ToLower(match[2])    // Convertir valor a minúsculas si aplica
		flagValue = strings.Trim(flagValue, "\"") // Eliminar comillas

		if flagName == "p" && flagValue == "" {
			pFlagPresent = true
			continue // No establecerlo aún, lo haremos después
		}

		// Asigna el flag en la estructura fs
		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Error: No se pudo establecer el flag '%s'\n", flagName)
			return
		}
		parsedFlags[flagName] = flagValue // Guardar para depuración
	}

	// Si el flag `-p` estuvo presente sin valor, establecerlo en "true"
	if pFlagPresent {
		_ = fs.Set("p", "true")
	}

	// Verificar que el parámetro -path sea obligatorio
	if *path == "" {
		fmt.Println("Error: el parametro '-path' es obligatorio")
		return
	}
	FileSystem.Mkdir(*path, *p)
}
