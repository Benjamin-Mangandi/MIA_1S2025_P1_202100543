package CommandsDisk

import (
	"Backend/DiskManager"
	"Backend/Globals"
	"Backend/Responsehandler"
	"flag"
	"os"
	"strings"
)

func MkDisk(params string) {
	// Definir flags con valores predeterminados
	fs := flag.NewFlagSet("mkdisk", flag.ContinueOnError)
	fit := fs.String("fit", "ff", "Tipo de ajuste (bf, ff, wf)")
	unit := fs.String("unit", "m", "Unidad de tamaño (k, m)")
	path := fs.String("path", "", "Ruta donde se creará el disco")
	size := fs.Int("size", 0, "Tamaño del disco")

	fs.Parse(os.Args[1:])
	matches := Globals.Regex.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		flagName := strings.ToLower(match[1])     // Convertir nombre a minúsculas
		flagValue := strings.ToLower(match[2])    // Convertir valor a minúsculas si aplica
		flagValue = strings.Trim(flagValue, "\"") // Eliminar comillas
		fs.Set(flagName, flagValue)
	}
	// Validaciones
	if *size <= 0 {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: El tamaño debe ser mayor a 0"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	if *path == "" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: Se requiere especificar un path"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Validar unidad
	if *unit != "k" && *unit != "m" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: Unidad inválida, debe ser 'k' o 'm'"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Llamar a la función para crear el disco
	DiskManager.Mkdisk(*size, *fit, *unit, *path)

}
