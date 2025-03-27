package FileSystem

import (
	"Backend/DiskManager"
	"Backend/Globals"
	"Backend/Responsehandler"
	Ext2 "Backend/Structs/ext2"
	"Backend/UsersManager"
	"Backend/Utilities"
	"fmt"
	"strings"
)

// Cat busca y muestra el contenido de los archivos en el sistema EXT2
func Cat(files map[string]string) {
	// Verificar si hay una sesión activa
	if !Globals.ActiveUser.Status {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No hay un usuario activo en el sistema."
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Obtener la partición montada asociada a la sesión
	mountedPartition := DiskManager.GetMountedPartitionByID(Globals.ActiveUser.PartitionID)
	if mountedPartition.ID == "" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No se encontró la partición montada."
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Abrir el archivo del disco
	file, err := Utilities.OpenFile(mountedPartition.Path)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo del disco.")
		return
	}
	defer file.Close()

	// Leer el superbloque de la partición
	var superblock Ext2.Superblock
	if err := Utilities.ReadObject(file, &superblock, int64(mountedPartition.Start)); err != nil {
		fmt.Println("Error: No se pudo leer el superbloque.")
		return
	}

	// Iterar sobre los archivos solicitados
	for key, filePath := range files {
		response := strings.Repeat("-", 40) +
			fmt.Sprintf("\nArchivo [%s]: %s\n", key, filePath)

		// Buscar el inodo del archivo en el sistema de archivos
		inodeIndex := UsersManager.FindFileInode(file, superblock, filePath)
		if inodeIndex == -1 {
			response := strings.Repeat("*", 30) + "\n" +
				"Error: Archivo no encontrado."
			Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
			continue
		}

		// Leer el contenido del archivo
		content := Globals.ReadFileFromInode(file, superblock, inodeIndex, Globals.ActiveUser)
		if content == "" {
			response := strings.Repeat("*", 30) + "\n" +
				"Error: No se tienen permisos para leer el archivo"
			Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
			continue
		}

		// Imprimir el contenido del archivo
		response += "  Contenido del archivo:\n" +
			strings.TrimSpace(content) + "\n"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
	}
}
