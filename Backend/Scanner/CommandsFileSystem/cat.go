package CommandsFileSystem

import (
	"Backend/FileSystem"
	"Backend/Globals"
	"Backend/Responsehandler"
	"strings"
)

func Cat(params string) {
	matches := Globals.Regex.FindAllStringSubmatch(params, -1)
	files := make(map[string]string) // Almacena los archivos detectados

	// Iterar sobre las coincidencias encontradas
	for _, match := range matches {
		flagName := strings.ToLower(match[1])     // fileN (ej. file1, file2, file3)
		flagValue := strings.ToLower(match[2])    // Convertir valor a minúsculas si aplica
		flagValue = strings.Trim(flagValue, "\"") // Eliminar comillas

		files[flagName] = flagValue
	}

	// Si no se encontraron archivos, mostrar error
	if len(files) == 0 {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No se especificaron archivos."
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	FileSystem.Cat(files)
}
