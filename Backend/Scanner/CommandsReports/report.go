package CommandsReports

import (
	"Backend/Globals"
	"Backend/Reports"
	"flag"
	"fmt"
	"strings"
)

func Report(params string) {
	fs := flag.NewFlagSet("rep", flag.ExitOnError)
	name := fs.String("name", "", "Tipo de reporte")
	path := fs.String("path", "", "Ubicacion donde se creara el reporte")
	id := fs.String("id", "", "ID de la particion")
	path_file_ls := fs.String("path_file_ls", "", "ID de la particion")

	matches := Globals.Regex.FindAllStringSubmatch(params, -1)

	// Mapa para almacenar los valores ingresados por el usuario
	parsedFlags := make(map[string]string)

	for _, match := range matches {
		flagName := match[1]                      // Flag tal cual fue escrito
		flagValue := strings.Trim(match[2], "\"") // Quita comillas si las tiene

		// Asigna el flag en la estructura fs
		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Error: No se pudo establecer el flag '%s'\n", flagName)
			return
		}
		parsedFlags[flagName] = flagValue // Guardar para depuración
	}

	// Validación de campos obligatorios
	if *name == "" || *path == "" || *id == "" {
		fmt.Println("Error: Los parámetros '-name', '-path' y '-id' son obligatorios")
		return
	}
	if *name == "is" || *name == "file" && *path_file_ls == "" {
		fmt.Println("Error: Los parámetros '-path_file_ls' es obligatorio")
		return
	}

	if *name == "mbr" {
		Reports.CreateMBR_Report(*path, *id)
	}
	if *name == "disk" {
		Reports.CreateDiskReport(*path, *id)
	}
	if *name == "inode" {
		Reports.CreateInode_Report(*path, *id)
	}
	if *name == "block" {
		fmt.Println("LLamando a la funcion para block")
	}
	if *name == "bm_inode" {
		Reports.CreateBmInodeReport(*path, *id)
	}
	if *name == "bm_block" {
		Reports.CreateBmBlockReport(*path, *id)
	}
	if *name == "tree" {
		fmt.Println("LLamando a la funcion para block")
	}
	if *name == "sb" {
		Reports.CreateSbReport(*path, *id)
	}
	if *name == "file" {
		Reports.CreateFileReport(*path, *id, *path_file_ls)
	}
	if *name == "is" {
		fmt.Println("LLamando a la funcion para bm_block")
	}
}
