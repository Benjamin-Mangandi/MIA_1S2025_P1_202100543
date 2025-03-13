package CommandsDisk

import (
	"Backend/DiskManager"
	"Backend/Globals"
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
		fmt.Println("Error: El tamaño debe ser mayor a 0")
		return
	}
	if *path == "" {
		fmt.Println("Error: Se requiere la ruta del archivo")
		return
	}
	if !validUnits[*unit] {
		fmt.Println("Error: Unidad inválida, debe ser 'b', 'k' o 'm'")
		return
	}
	if !validTypes[*type_] {
		fmt.Println("Error: Tipo inválido, debe ser 'p' (primaria), 'e' (extendida) o 'l' (logica)")
		return
	}
	if !validFits[*fit] {
		fmt.Println("Errora: Ajuste inválido, debe ser 'bf', 'ff' o 'wf'")
		return
	}

	// Llamar a la función Fdisk con los valores procesados
	DiskManager.Fdisk(*size, *path, *name, *unit, *type_, *fit)

}

//fdisk -size=600 -path="/home/benjamin/discos/Disco1.mia" -name=Particion1
//fdisk -Size=2000 -path=/home/benjamin/discos/Disco1.mia -name=Particion2
//fdisk -type=E -path=/home/benjamin/discos/Disco1.mia -Unit=K -name=Particion3 -size=300
//fdisk -size=1 -type=L -unit=M -fit=BF -path=/home/benjamin/discos/Disco1.mia -name="Particion4"
//fdisk -path=/home/benjamin/discos/Disco2.mia -name=Part3 -Unit=K -size=200
