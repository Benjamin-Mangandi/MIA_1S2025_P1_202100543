package CommandsFileSystem

import (
	"Backend/FileSystem"
	"Backend/Globals"
	"flag"
	"fmt"
	"strings"
)

func MkFile(params string) {
	fs := flag.NewFlagSet("mkfile", flag.ContinueOnError)
	path := fs.String("path", "", "Ruta")
	size := fs.String("size", "0", "Tamaño del archivo en bytes")
	cont := fs.String("cont", "", "Path de un archivo en la PC")
	r := fs.String("r", "false", "Creación de carpetas padre")

	// Extraer parámetros usando Regex
	matches := Globals.Regex.FindAllStringSubmatch(params, -1)
	parsedFlags := make(map[string]string)

	rFlagPresent := false // Variable para rastrear si -r está presente sin valor

	for _, match := range matches {
		flagName := strings.ToLower(match[1])     // fileN (ej. file1, file2, file3)
		flagValue := strings.ToLower(match[2])    // Convertir valor a minúsculas si aplica
		flagValue = strings.Trim(flagValue, "\"") // Eliminar comillas

		if flagName == "r" && flagValue == "" {
			rFlagPresent = true
			continue // No establecerlo aún, lo haremos después
		}

		// Asigna el flag en la estructura fs
		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Error: No se pudo establecer el flag '%s'\n", flagName)
			return
		}
		parsedFlags[flagName] = flagValue // Guardar para depuración
	}

	// Si el flag `-r` estuvo presente sin valor, establecerlo en "true"
	if rFlagPresent {
		_ = fs.Set("r", "true")
	}

	// Verificar que el parámetro -path sea obligatorio
	if *path == "" {
		fmt.Println("Error: el parámetro '-path' es obligatorio")
		return
	}
	FileSystem.Mkfile(*path, *r, *size, *cont)
}
