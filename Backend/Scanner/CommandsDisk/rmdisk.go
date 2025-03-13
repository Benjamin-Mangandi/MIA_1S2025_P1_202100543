package CommandsDisk

import (
	"Backend/DiskManager"
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
		if strings.HasPrefix(args[i], "-path=") {
			*path = strings.TrimPrefix(args[i], "-path=")
			*path = strings.Trim(*path, "\"") // Eliminar comillas si las hay
		}
	}

	// Validar si se proporcionó el path
	if *path == "" {
		fmt.Println("Error: El parámetro -path es obligatorio.")
		return
	}

	// Llamar a la función en DiskManager para eliminar el disco
	DiskManager.RmDisk(*path)
}
