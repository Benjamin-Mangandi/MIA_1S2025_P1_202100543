package CommandsDisk

import (
	"Backend/DiskManager"
	"Backend/Responsehandler"
	"flag"
	"fmt"
	"strings"
)

func RmDisk(input string) {
	// Definir flags
	fs := flag.NewFlagSet("rmdisk", flag.ExitOnError)
	path := fs.String("path", "", "Ruta del disco a eliminar")

	// Buscar y extraer los flags del input
	args := strings.Fields(input)

	for i := 0; i < len(args); i++ {
		argName := strings.ToLower(args[i])
		if strings.HasPrefix(argName, "-path=") {
			*path = strings.TrimPrefix(argName, "-path=")
			*path = strings.Trim(*path, "\"") // Eliminar comillas si las hay
		}

	}

	// Validar si se proporcionó el path
	if *path == "" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: El parámetro -path es obligatorio."
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		fmt.Println()
		return
	}

	// Llamar a la función en DiskManager para eliminar el disco
	DiskManager.RmDisk(*path)
}
