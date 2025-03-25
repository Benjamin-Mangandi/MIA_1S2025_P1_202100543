package DiskManager

import (
	"Backend/Responsehandler"
	"fmt"
	"os"
)

func RmDisk(path string) {

	// Verificar si el archivo existe antes de eliminarlo
	if _, err := os.Stat(path); os.IsNotExist(err) {
		response := fmt.Sprintf("---------------------\n"+
			"Error: No se encontró el disco en la ruta '%s'.\n", path)
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)

		return
	}

	// Intentar eliminar el archivo
	err := os.Remove(path)
	if err != nil {
		response := "---------------------\n" +
			"Error: No se pudo eliminar el disco."
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	response := fmt.Sprintf("---------------------\n"+
		"Disco Eliminado Correctamente: '%s'.\n", path)
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
}
