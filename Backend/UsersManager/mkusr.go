package UsersManager

import (
	"Backend/DiskManager"
	"Backend/Globals"
	Ext2 "Backend/Structs/ext2"
	"Backend/Utilities"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Mkusr(name string, pass string, group string) {
	// Verificar si hay una sesión activa
	if !Globals.ActiveUser.Status {
		fmt.Println("Error: No hay un usuario activo.")
		return
	}
	if Globals.ActiveUser.Name != "root" {
		fmt.Println("Error: Solo el usuario root puede crear usuarios.")
		return
	}

	// Obtener la partición montada asociada a la sesión
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

	// Leer el superbloque desde la partición
	var superblock Ext2.Superblock
	if err := Utilities.ReadObject(file, &superblock, int64(mountedPartition.Start)); err != nil {
		fmt.Println("Error al leer el superbloque.")
		return
	}

	// Buscar el inodo del archivo users.txt
	inodeIndex := FindFileInode(file, superblock, "users.txt")
	if inodeIndex == -1 {
		fmt.Println("Error: No se encontró el archivo users.txt.")
		return
	}

	// Leer el contenido actual de users.txt
	fileContent := ReadFileFromInode(file, superblock, inodeIndex)
	if fileContent == "" {
		fmt.Println("Error: No se pudo leer el archivo users.txt.")
		return
	}

	// Verificar si el grupo existe
	lines := strings.Split(fileContent, "\n")
	groupExists := false
	highestID := 0

	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 3 && parts[1] == "G" {
			// Es un grupo
			if parts[2] == group && parts[0] != "0" { // Verificar que el grupo no esté eliminado
				groupExists = true
			}
		} else if len(parts) == 5 && parts[1] == "U" {
			// Es un usuario, verificar el ID más alto
			id, err := strconv.Atoi(parts[0])
			if err == nil && id > highestID {
				highestID = id
			}
			// Verificar si el usuario ya existe
			if parts[3] == name {
				fmt.Println("Error: El usuario ya existe.")
				return
			}
		}
	}

	if !groupExists {
		fmt.Println("Error: El grupo especificado no existe.")
		return
	}

	// Generar la nueva línea del usuario
	newUser := fmt.Sprintf("%d,U,%s,%s,%s\n", highestID+1, group, name, pass)
	// Escribir el contenido actualizado en users.txt
	newContent := fileContent + newUser
	if err := WriteFileToInode(file, superblock, inodeIndex, newContent, mountedPartition); err != nil {
		fmt.Println("Error: No se pudo escribir en users.txt.")
		return
	}

	fmt.Println("Usuario creado exitosamente:", name)
}

func AllocateNewBlock(superblock *Ext2.Superblock, file *os.File) int {
	// Leer el bitmap de bloques
	bitmapStart := int64(superblock.S_bm_block_start)
	bitmapSize := int(superblock.S_blocks_count)
	bitmap := make([]byte, bitmapSize)

	if _, err := file.ReadAt(bitmap, bitmapStart); err != nil {
		fmt.Println("Error: No se pudo leer el bitmap de bloques.")
		return -1
	}

	// Buscar el primer bloque libre (byte con valor 0)
	for i := 0; i < bitmapSize; i++ {
		if bitmap[i] == 0 {
			// Marcar el bloque como ocupado (1)
			bitmap[i] = 1
			if _, err := file.WriteAt([]byte{1}, bitmapStart+int64(i)); err != nil {
				fmt.Println("Error: No se pudo actualizar el bitmap de bloques.")
				return -1
			}

			// Actualizar el contador de bloques libres
			superblock.S_free_blocks_count--

			return i // Retorna el índice del bloque asignado
		}
	}

	fmt.Println("Error: No hay bloques disponibles.")
	return -1
}
