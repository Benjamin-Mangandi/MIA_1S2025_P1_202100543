package CommandsFileSystem

import (
	"Backend/FileSystem"
	"Backend/Globals"
	"flag"
	"fmt"
	"strings"
)

func Mkfs(params string) {
	fs := flag.NewFlagSet("mkfs", flag.ExitOnError)
	id := fs.String("id", "", "Id")
	type_ := fs.String("type", "full", "Tipo")
	fs_ := fs.String("fs", "2fs", "Fs")

	// Convertir parámetros a minúsculas para hacer la búsqueda insensible a mayúsculas
	params = strings.ToLower(params)

	// Expresión regular para extraer parámetros
	matches := Globals.Regex.FindAllStringSubmatch(params, -1)
	validFlags := map[string]*string{"id": id, "type": type_, "fs": fs_}

	for _, match := range matches {
		flagName, flagValue := strings.ToLower(match[1]), strings.Trim(match[2], "\"")

		if ptr, exists := validFlags[flagName]; exists {
			*ptr = flagValue
		} else {
			fmt.Printf("Error: Flag '%s' no encontrada\n", flagName)
			return
		}
	}

	// Validar que 'id' no esté vacío
	if *id == "" {
		fmt.Println("Error: El parámetro 'id' es obligatorio.")
		return
	}

	FileSystem.Mkfs(strings.ToLower(*id), *type_, *fs_)
}

//mkfs -id=431a
//mkfs -type=full -id=431a
