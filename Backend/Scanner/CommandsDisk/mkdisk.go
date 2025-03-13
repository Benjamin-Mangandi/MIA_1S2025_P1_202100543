package CommandsDisk

import (
	"Backend/DiskManager"
	"Backend/Globals"
	"flag"
	"fmt"
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
		fmt.Println("Error: El tamaño debe ser mayor a 0")
		return
	}

	if *path == "" {
		fmt.Println("Error: Se requiere especificar un path")
		return
	}

	// Validar unidad
	if *unit != "k" && *unit != "m" {
		fmt.Println("Error: Unidad inválida, debe ser 'k' o 'm'")
		return
	}

	// Llamar a la función para crear el disco
	DiskManager.Mkdisk(*size, *fit, *unit, *path)

}

//mkdisk -Size=3000 -unit=K -path=/home/benjamin/discos/Disco1.mia
//mkdisk -path=/home/benjamin/discos/Disco2.mia -Unit=K -size=3000
//mkdisk -size=5 -unit=M -path="/home/benjamin/discos/Disco3.mia"​
//mkdisk -size=10 -Path="/home/benjamin/discos/Disco4.mia"
//mkdisk -size=1024 -unit=M -fit=BF -path="/home/benjamin/discos/DiscoDefinitivo.mia"
