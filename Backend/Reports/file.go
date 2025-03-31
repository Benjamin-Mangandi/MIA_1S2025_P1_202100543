package Reports

import (
	"Backend/DiskManager"
	"Backend/Globals"
	"Backend/Responsehandler"
	Ext2 "Backend/Structs/ext2"
	"Backend/UsersManager"
	"Backend/Utilities"
	"fmt"
	"os"
	"strings"
)

func CreateFileReport(path string, id string, filePath string) {
	// Buscar la partición montada
	path = fixPath(path)
	err := Utilities.CreateParentDirs(path)
	if err != nil {
		response := strings.Repeat("*", 30) + "\n" +
			"Error al crear las carpetas padre"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	mountedPartition := DiskManager.GetMountedPartitionByID(id)
	if mountedPartition.ID == "" {
		fmt.Println("Error: No se encontró la partición montada.")
		return
	}

	// Abrir el archivo del disco
	file, err := os.Open(mountedPartition.Path)
	if err != nil {
		fmt.Println("Error al abrir el disco:", err)
		return
	}
	defer file.Close()

	// Leer el superbloque
	var superblock Ext2.Superblock
	if err := Utilities.ReadObject(file, &superblock, int64(mountedPartition.Start)); err != nil {
		fmt.Println("Error al leer el superbloque:", err)
		return
	}

	// Buscar el inodo del archivo en el sistema de archivos
	inodeIndex := UsersManager.FindFileInode(file, superblock, filePath)
	if inodeIndex == -1 {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: Archivo no encontrado."
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Leer el contenido del archivo
	content := Globals.ReadFileFromInode(file, superblock, inodeIndex, Globals.ActiveUser)
	if content == "" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No se tienen permisos para leer el archivo"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	// Escribir el contenido en un archivo de texto en la ruta especificada
	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error al escribir el archivo de reporte:", err)
		return
	}
	response := strings.Repeat("-", 40) + "\n" +
		"Reporte File creado exitosamente en:" + path + "\n"
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
}
