package CommandsDisk

import (
	"Backend/DiskManager"
	"Backend/Globals"
	"Backend/Responsehandler"
	"flag"
	"fmt"
	"strings"
)

func FDisk(input string) {
	// Definir flags con valores predeterminados
	fs := flag.NewFlagSet("fdisk", flag.ContinueOnError)
	size := fs.Int("size", 0, "Tamaño")
	path := fs.String("path", "", "Ruta")
	name := fs.String("name", "", "Nombre")
	unit := fs.String("unit", "k", "Unidad (b/k/m)")
	type_ := fs.String("type", "p", "Tipo (p/e/l)")
	fit := fs.String("fit", "wf", "Ajuste (bf/ff/wf)")

	// Expresión regular para capturar los flags en el input
	matches := Globals.Regex.FindAllStringSubmatch(input, -1)

	// Mapas para validaciones rápidas
	validUnits := map[string]bool{"b": true, "k": true, "m": true}
	validTypes := map[string]bool{"p": true, "e": true, "l": true}
	validFits := map[string]bool{"bf": true, "ff": true, "wf": true}

	// Procesar los flags extraídos
	for _, match := range matches {
		flagName := strings.ToLower(match[1])     // Convertir nombre a minúsculas
		flagValue := strings.ToLower(match[2])    // Convertir valor a minúsculas si aplica
		flagValue = strings.Trim(flagValue, "\"") // Eliminar comillas

		// Intentar asignar el valor al flag correspondiente
		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Advertencia: Flag desconocido '%s'\n", flagName)
		}
	}

	// Validaciones de parámetros
	if *size <= 0 {
		response := "---------------------\n" +
			"Error: El tamaño de la particion debe ser mayor a 0"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	if *path == "" {
		response := "---------------------\n" +
			"Error: Se requiere la ruta del archivo"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	if !validUnits[*unit] {
		response := "---------------------\n" +
			"Error: Unidad inválida, debe ser 'b', 'k' o 'm'"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	if !validTypes[*type_] {
		response := "---------------------\n" +
			"Error: Tipo inválido, debe ser 'p' (primaria), 'e' (extendida) o 'l' (logica)"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	if !validFits[*fit] {
		response := "---------------------\n" +
			"Error: Ajuste inválido, debe ser 'bf', 'ff' o 'wf'"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Llamar a la función Fdisk con los valores procesados
	DiskManager.Fdisk(*size, *path, *name, *unit, *type_, *fit)

}
