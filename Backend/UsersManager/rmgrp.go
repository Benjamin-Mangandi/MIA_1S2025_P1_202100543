package UsersManager

import (
	"Backend/DiskManager"
	"Backend/Globals"
	Ext2 "Backend/Structs/ext2"
	"Backend/Utilities"
	"fmt"
	"strings"
)

func Rmgrp(name string) {
	// Verificar si hay una sesión activa
	if !Globals.ActiveUser.Status {
		fmt.Println("Error: No hay un usuario activo.")
		return
	}
	if Globals.ActiveUser.Name != "root" {
		fmt.Println("Error: Solo el usuario root puede crear grupos")
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

	// 3. Leer el superbloque desde la partición
	var superblock Ext2.Superblock
	if err := Utilities.ReadObject(file, &superblock, int64(mountedPartition.Start)); err != nil {
		fmt.Println("Error al leer el superbloque.")
		return
	}

	// Buscar el inodo del archivo users.txt
	inodeIndex := FindFileInode(file, superblock, "users.txt")
	if inodeIndex == -1 {
		fmt.Println("Error: No se encontró el archivo users.txt")
		return
	}

	// Leer el contenido actual de users.txt
	fileContent := ReadFileFromInode(file, superblock, inodeIndex)
	if fileContent == "" {
		fmt.Println("Error: No se pudo leer el archivo users.txt.")
		return
	}

	// Modificar la línea del grupo
	lines := strings.Split(fileContent, "\n")
	for i, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 3 && parts[1] == "G" && parts[2] == name {
			parts[0] = "0" // Marcar el grupo como eliminado
			lines[i] = strings.Join(parts, ",")
			break
		}
	}

	// Escribir el contenido actualizado en users.txt
	newContent := strings.Join(lines, "\n")
	if err := WriteFileToInode(file, superblock, inodeIndex, newContent, mountedPartition); err != nil {
		fmt.Println("Error: No se pudo escribir en users.txt")
		return
	}

	fmt.Println("Grupo eliminado correctamente:", name)
}
