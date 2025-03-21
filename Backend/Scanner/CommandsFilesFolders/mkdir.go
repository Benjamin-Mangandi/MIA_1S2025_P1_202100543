package CommandsFilesFolders

import (
	"Backend/Globals"
	"flag"
	"fmt"
	"strings"
)

func Mkdir(params string) {
	fs := flag.NewFlagSet("mkdir", flag.ContinueOnError)
	path := fs.String("path", "", "Ruta")
	p := fs.String("p", "false", "Creacion de carpetas padre") // Cambié "-r" a "-p" para que coincida con el parámetro correcto.

	// Extraer parámetros usando Regex
	matches := Globals.Regex.FindAllStringSubmatch(params, -1)
	parsedFlags := make(map[string]string)

	pFlagPresent := false // Variable para rastrear si -p está presente sin valor

	for _, match := range matches {
		flagName := match[1]                      // Nombre del flag
		flagValue := strings.Trim(match[2], "\"") // Quitar comillas

		if flagName == "p" && flagValue == "" {
			pFlagPresent = true
			continue // No establecerlo aún, lo haremos después
		}

		// Asigna el flag en la estructura fs
		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Error: No se pudo establecer el flag '%s'\n", flagName)
			return
		}
		parsedFlags[flagName] = flagValue // Guardar para depuración
	}

	// Si el flag `-p` estuvo presente sin valor, establecerlo en "true"
	if pFlagPresent {
		_ = fs.Set("p", "true")
	}

	// Imprimir parámetros detectados para depuración
	fmt.Println("====== Parámetros Escaneados ======")
	for key, value := range parsedFlags {
		fmt.Printf("%s: %s\n", key, value)
	}
	fmt.Println("===================================")

	// Verificar que el parámetro -path sea obligatorio
	if *path == "" {
		fmt.Println("Error: el parametro '-path' es obligatorio")
		return
	}

	// Verificar si -p se estableció correctamente
	fmt.Println("Valor de -p:", *p)
}
