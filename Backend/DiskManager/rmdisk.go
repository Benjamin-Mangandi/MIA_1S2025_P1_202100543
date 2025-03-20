package DiskManager

import (
	"fmt"
	"os"
)

func RmDisk(path string) {

	// Verificar si el archivo existe antes de eliminarlo
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Error: No se encontró el disco en la ruta '%s'.\n", path)
		return
	}

	// Confirmar con el usuario antes de eliminar
	fmt.Printf("¿Está seguro de que desea eliminar el disco en '%s'? (s/n): ", path)
	var response string
	fmt.Scanln(&response)

	if response != "s" && response != "S" {
		fmt.Println("Operación cancelada.")
		return
	}

	// Intentar eliminar el archivo
	err := os.Remove(path)
	if err != nil {
		fmt.Printf("Error: No se pudo eliminar el disco. %v\n", err)
		return
	}

	fmt.Println("Disco eliminado con éxito.")
}
