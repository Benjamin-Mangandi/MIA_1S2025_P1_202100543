package CommandsFileSystem

import (
	"Backend/FileSystem"
	"Backend/Globals"
	"fmt"
	"strings"
)

func Cat(params string) {
	matches := Globals.Regex.FindAllStringSubmatch(params, -1)
	files := make(map[string]string) // Almacena los archivos detectados

	// Iterar sobre las coincidencias encontradas
	for _, match := range matches {
		flagName := match[1]                      // fileN (ej. file1, file2, file3)
		flagValue := strings.Trim(match[2], "\"") // Ruta del archivo sin comillas

		files[flagName] = flagValue
	}

	// Si no se encontraron archivos, mostrar error
	if len(files) == 0 {
		fmt.Println("Error: No se especificaron archivos.")
		return
	}

	// Mostrar archivos detectados
	fmt.Println("Archivos detectados:")
	for key, value := range files {
		fmt.Printf("%s -> %s\n", key, value)
	}
	FileSystem.Cat(files)
}
