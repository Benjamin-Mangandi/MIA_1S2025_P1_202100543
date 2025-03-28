package CommandsFilesFolders

import (
	"Backend/Globals"
	"flag"
	"fmt"
	"strings"
)

func MkFile(params string) {
	fs := flag.NewFlagSet("mkfile", flag.ContinueOnError)
	path := fs.String("path", "", "Ruta")
	size := fs.String("size", "0", "Tamaño del archivo en bytes")
	cont := fs.String("cont", "", "Path de un archivo en la PC")
	r := fs.String("r", "false", "Creación de carpetas padre")

	// Extraer parámetros usando Regex
	matches := Globals.Regex.FindAllStringSubmatch(params, -1)
	parsedFlags := make(map[string]string)

	rFlagPresent := false // Variable para rastrear si -r está presente sin valor

	for _, match := range matches {
		flagName := match[1]                      // Nombre del flag
		flagValue := strings.Trim(match[2], "\"") // Quitar comillas

		if flagName == "r" && flagValue == "" {
			rFlagPresent = true
			continue // No establecerlo aún, lo haremos después
		}

		// Asigna el flag en la estructura fs
		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Error: No se pudo establecer el flag '%s'\n", flagName)
			return
		}
		parsedFlags[flagName] = flagValue // Guardar para depuración
	}

	// Si el flag `-r` estuvo presente sin valor, establecerlo en "true"
	if rFlagPresent {
		_ = fs.Set("r", "true")
	}

	// Imprimir parámetros detectados para depuración
	fmt.Println("====== Parámetros Escaneados ======")
	for key, value := range parsedFlags {
		fmt.Printf("%s: %s\n", key, value)
	}
	fmt.Println("===================================")

	// Verificar que el parámetro -path sea obligatorio
	if *path == "" {
		fmt.Println("Error: el parámetro '-path' es obligatorio")
		return
	}

	// Verificar si los parámetros se establecieron correctamente
	fmt.Println("Valores:")
	fmt.Println("Path:", *path)
	fmt.Println("Size:", *size)
	fmt.Println("R:", *r)
	fmt.Println("Cont:", *cont)
}
