package FileSystem

import (
	"Backend/DiskManager"
	"Backend/Globals"
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
		fmt.Println("Error: No hay un usuario activo.")
		return
	}

	// Obtener la partición montada asociada a la sesión
	fmt.Println(Globals.ActiveUser.ID)
	mountedPartition := DiskManager.GetMountedPartitionByID(Globals.ActiveUser.PartitionID)
	if mountedPartition.ID == "" {
		fmt.Println("Error: No se encontró la partición montada.")
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
		fmt.Printf("\nArchivo [%s]: %s\n", key, filePath)

		// Buscar el inodo del archivo en el sistema de archivos
		inodeIndex := UsersManager.FindFileInode(file, superblock, filePath)
		if inodeIndex == -1 {
			fmt.Println("  Error: Archivo no encontrado.")
			continue
		}

		// Leer el contenido del archivo
		content := UsersManager.ReadFileFromInode(file, superblock, inodeIndex)
		if content == "" {
			fmt.Println("  Error: No se pudo leer el contenido del archivo.")
			continue
		}

		// Imprimir el contenido del archivo
		fmt.Println("  Contenido del archivo:")
		fmt.Println(strings.TrimSpace(content))
	}
}
