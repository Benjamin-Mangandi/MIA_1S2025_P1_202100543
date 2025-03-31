package CommandsReports

import (
	"Backend/Globals"
	"Backend/Reports"
	"Backend/Responsehandler"
	"flag"
	"fmt"
	"strings"
)

func Report(params string) {
	fs := flag.NewFlagSet("rep", flag.ExitOnError)
	name := fs.String("name", "", "Tipo de reporte")
	path := fs.String("path", "", "Ubicacion donde se creara el reporte")
	id := fs.String("id", "", "ID de la particion")
	path_file_ls := fs.String("path_file_ls", "", "Path de archivo en el disco")

	matches := Globals.Regex.FindAllStringSubmatch(params, -1)

	// Mapa para almacenar los valores ingresados por el usuario
	parsedFlags := make(map[string]string)

	for _, match := range matches {
		flagName := strings.ToLower(match[1])     // Convertir nombre a minúsculas
		flagValue := strings.ToLower(match[2])    // Convertir valor a minúsculas si aplica
		flagValue = strings.Trim(flagValue, "\"") // Eliminar comillas
		// Asigna el flag en la estructura fs
		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Error: No se pudo establecer el flag '%s'\n", flagName)
			return
		}
		parsedFlags[flagName] = flagValue // Guardar para depuración
	}

	// Validación de campos obligatorios
	if *name == "" || *path == "" || *id == "" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: Los parámetros '-name', '-path' y '-id' son obligatorios" + "\n"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	if *name == "is" || *name == "file" && *path_file_ls == "" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: Los parámetros '-path_file_ls' es obligatorio" + "\n"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	switch *name {
	case "mbr":
		Reports.CreateMBR_Report(*path, *id)
	case "disk":
		Reports.CreateDiskReport(*path, *id)
	case "inode":
		Reports.CreateInode_Report(*path, *id)
	case "block":
		Reports.CreateBlocksReport(*path, *id)
	case "bm_inode":
		Reports.CreateBmInodeReport(*path, *id)
	case "bm_block":
		Reports.CreateBmBlockReport(*path, *id)
	case "tree":
		fmt.Println("Llamando a la función para tree")
	case "sb":
		Reports.CreateSbReport(*path, *id)
	case "file":
		Reports.CreateFileReport(*path, *id, *path_file_ls)
	case "is":
		fmt.Println("Llamando a la función para is")
	default:
		response := strings.Repeat("*", 30) + "\n" +
			"Error: Reporte no reconocido." + "\n"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)

	}

}
