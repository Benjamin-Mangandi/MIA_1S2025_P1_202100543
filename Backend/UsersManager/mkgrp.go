package UsersManager

import (
	"Backend/DiskManager"
	"Backend/Globals"
	Disk "Backend/Structs/disk"
	Ext2 "Backend/Structs/ext2"
	"Backend/Utilities"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Mkgrp(name string) {
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

	// Verificar si el grupo ya existe
	lines := strings.Split(fileContent, "\n")
	highestID := 0
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 3 && parts[1] == "G" {
			// Obtener el ID del grupo
			id, err := strconv.Atoi(parts[0])
			if err == nil && id > highestID {
				highestID = id
			}
			if parts[2] == name {
				fmt.Println("Error: El grupo ya existe.")
				return
			}
		}
	}

	// Generar la nueva línea del grupo
	newGroup := fmt.Sprintf("%d,G,%s\n", highestID+1, name)
	fileContent += newGroup

	if err := WriteFileToInode(file, superblock, inodeIndex, fileContent, mountedPartition); err != nil {
		fmt.Println("Error: No se pudo escribir en users.txt")
		return
	}

	fmt.Println("Grupo creado exitosamente:", name)
}

func WriteFileToInode(file *os.File, superblock Ext2.Superblock, inodeIndex int, newContent string, mountedPartition Disk.MountedPartition) error {
	// Leer el inodo del archivo
	var inode Ext2.Inode
	inodeStart := int64(superblock.S_inode_start) + int64(inodeIndex)*int64(binary.Size(Ext2.Inode{}))
	if err := Utilities.ReadObject(file, &inode, inodeStart); err != nil {
		return err
	}

	data := []byte(newContent)
	dataLen := len(data)
	blockSize := binary.Size(Ext2.Fileblock{})

	// Verificar que el inodo tiene bloques asignados
	if inode.I_block[0] == -1 {
		return fmt.Errorf("el archivo no tiene bloques asignados")
	}

	// Escribir en los bloques asignados
	blockIndex := 0
	for dataLen > 0 && blockIndex < len(inode.I_block) {
		if inode.I_block[blockIndex] == -1 {
			break
		}

		blockOffset := int64(superblock.S_block_start) + int64(inode.I_block[blockIndex])*int64(blockSize)
		writeSize := min(dataLen, blockSize)

		// Crear un buffer limpio para evitar datos residuales
		copyBlock := make([]byte, blockSize)
		copy(copyBlock, data[:writeSize])

		// Escribir en el bloque
		if err := Utilities.WriteObject(file, copyBlock, blockOffset); err != nil {
			return err
		}

		// Moverse al siguiente bloque
		data = data[writeSize:]
		dataLen -= writeSize
		blockIndex++
	}

	// Si aún quedan datos, necesitamos asignar más bloques
	for dataLen > 0 && blockIndex < len(inode.I_block) {
		newBlockIndex := AllocateNewBlock(&superblock, file) // Asignar un nuevo bloque
		if newBlockIndex == -1 {
			return fmt.Errorf("no hay espacio suficiente en la partición")
		}

		// Guardar el nuevo bloque en el inodo
		inode.I_block[blockIndex] = int32(newBlockIndex)

		blockOffset := int64(superblock.S_block_start) + int64(newBlockIndex)*int64(blockSize)
		writeSize := min(dataLen, blockSize)

		// Crear un buffer limpio
		copyBlock := make([]byte, blockSize)
		copy(copyBlock, data[:writeSize])

		// Escribir en el nuevo bloque
		if err := Utilities.WriteObject(file, copyBlock, blockOffset); err != nil {
			return err
		}

		// Moverse al siguiente bloque
		data = data[writeSize:]
		dataLen -= writeSize
		blockIndex++
	}

	// Limpiar bloques restantes si el nuevo contenido es menor al anterior
	for ; blockIndex < len(inode.I_block) && inode.I_block[blockIndex] != -1; blockIndex++ {
		// Limpiar el bloque (lo sobrescribimos con ceros)
		blockOffset := int64(superblock.S_block_start) + int64(inode.I_block[blockIndex])*int64(blockSize)
		emptyBlock := make([]byte, blockSize)
		if err := Utilities.WriteObject(file, emptyBlock, blockOffset); err != nil {
			return err
		}

		// Marcar el bloque como no asignado
		inode.I_block[blockIndex] = -1
	}

	// Actualizar tamaño del inodo
	inode.I_size = int32(len(newContent))

	// Escribir el inodo actualizado en el disco
	if err := Utilities.WriteObject(file, inode, inodeStart); err != nil {
		return err
	}

	// Escribir el superbloque actualizado (para reflejar cambios en bloques)
	if err := Utilities.WriteObject(file, superblock, int64(mountedPartition.Start)); err != nil {
		return err
	}

	return nil
}
