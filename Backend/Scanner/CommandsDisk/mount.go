package CommandsDisk

import (
	"Backend/DiskManager"
	"Backend/Globals"
	"flag"
	"fmt"
	"strings"
)

func Mount(params string) {
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	path := fs.String("path", "", "Ruta del disco")
	name := fs.String("name", "", "Nombre de la partición")

	// Extraer los parámetros usando regex
	matches := Globals.Regex.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := strings.ToLower(match[1])     // Convertir nombre del flag a minúsculas
		flagValue := strings.Trim(match[2], "\"") // Eliminar comillas

		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Error: No se pudo establecer el flag '%s'\n", flagName)
			return
		}
	}

	// Validación de campos obligatorios
	if *path == "" || *name == "" {
		fmt.Println("Error: Los parámetros '-path' y '-name' son obligatorios")
		return
	}

	// Llamar a la función Mount con el nombre en minúsculas
	DiskManager.Mount(*path, strings.ToLower(*name))

}

//mount -path="/home/benjamin/discos/disco1.mia" -name=particion1
//mount -path="/home/benjamin/discos/disco1.mia" -name=particion2
//mount -path="/home/benjamin/discos/disco2.mia" -name=part3
//mount -path="/home/benjamin/discos/disco1.mia" -name=part3
