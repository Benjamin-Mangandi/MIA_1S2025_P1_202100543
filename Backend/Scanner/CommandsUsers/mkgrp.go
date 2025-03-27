package CommandsUsers

import (
	"Backend/Globals"
	"Backend/Responsehandler"
	"Backend/UsersManager"
	"flag"
	"fmt"
	"strings"
)

func Mkgrp(params string) {
	fs := flag.NewFlagSet("mkgrp", flag.ExitOnError)
	name := fs.String("name", "", "Nombre del grupo")

	// Extraer parámetros usando Regex
	matches := Globals.Regex.FindAllStringSubmatch(params, -1)
	parsedFlags := make(map[string]string)

	for _, match := range matches {
		flagName := strings.ToLower(match[1])     // Nombre del flag
		flagValue := strings.Trim(match[2], "\"") // Quitar comillas

		// Asigna el flag en la estructura fs
		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Error: No se pudo establecer el flag '%s'\n", flagName)
			return
		}
		parsedFlags[flagName] = flagValue // Guardar para depuración
	}

	// Verificar que el parámetro -name se haya ingresado
	if *name == "" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: el parámetro '-name' es obligatorio"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Verificar que el nombre del grupo no sea "root"
	if strings.ToLower(*name) == "root" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No se puede crear un grupo con el nombre 'root'"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	UsersManager.Mkgrp(*name)
}
