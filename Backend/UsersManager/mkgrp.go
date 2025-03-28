package UsersManager

import (
	"Backend/DiskManager"
	"Backend/Globals"
	"Backend/Responsehandler"
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
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No hay un usuario activo."
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	if Globals.ActiveUser.Name != "root" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: Solo el usuario root puede crear grupos"
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
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No se pudo leer el archivo users.txt."
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Verificar si el grupo ya existe
	lines := strings.Split(fileContent, "\n")
	highestID := 0
	for i, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 3 && parts[1] == "G" {
			// Obtener el ID del grupo
			id, err := strconv.Atoi(parts[0])
			if err == nil && id > highestID {
				highestID = id
			}
			if parts[2] == name && parts[0] != "0" {
				response := strings.Repeat("*", 30) + "\n" +
					"Error: El grupo ya existe."
				Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
				return
			} else if parts[2] == name && parts[0] == "0" {
				parts[0] = strconv.Itoa(highestID + 1) // Marcar el grupo como restaurado
				lines[i] = strings.Join(parts, ",")
				// Escribir el contenido actualizado en users.txt
				newContent := strings.Join(lines, "\n")
				if err := WriteFileToInode(file, &superblock, inodeIndex, newContent, mountedPartition); err != nil {
					fmt.Println("Error: No se pudo escribir en users.txt")
					return
				}
				response := strings.Repeat("-", 40) + "\n" +
					"Grupo creado exitosamente:" + name + "\n"
				Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
				return
			}
		}
	}

	// Generar la nueva línea del grupo
	newGroup := fmt.Sprintf("%d,G,%s\n", highestID+1, name)
	fileContent += newGroup

	if err := WriteFileToInode(file, &superblock, inodeIndex, fileContent, mountedPartition); err != nil {
		fmt.Println("Error: No se pudo escribir en users.txt")
		return
	}
	response := strings.Repeat("-", 40) + "\n" +
		"Grupo creado exitosamente:" + name + "\n"
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
}

func WriteFileToInode(file *os.File, superblock *Ext2.Superblock, inodeIndex int, newContent string, mountedPartition Disk.MountedPartition) error {
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

	blockIndex := 0
	for dataLen > 0 && blockIndex < 12 {
		if inode.I_block[blockIndex] == -1 {
			newBlockIndex := AllocateNewBlock(superblock, file)
			if newBlockIndex == -1 {
				return fmt.Errorf("no hay espacio suficiente en la partición")
			}
			inode.I_block[blockIndex] = int32(newBlockIndex)
		}

		blockOffset := int64(superblock.S_block_start) + int64(inode.I_block[blockIndex])*int64(blockSize)
		writeSize := min(dataLen, blockSize)
		copyBlock := make([]byte, blockSize)
		copy(copyBlock, data[:writeSize])

		if err := Utilities.WriteObject(file, copyBlock, blockOffset); err != nil {
			return err
		}

		data = data[writeSize:]
		dataLen -= writeSize
		blockIndex++
	}

	// Manejo de bloque indirecto simple
	if dataLen > 0 {
		if inode.I_block[12] == -1 {
			newBlockIndex := AllocateNewBlock(superblock, file)
			if newBlockIndex == -1 {
				return fmt.Errorf("no hay espacio para bloque indirecto simple")
			}
			inode.I_block[12] = int32(newBlockIndex)
		}

		indirectBlockOffset := int64(superblock.S_block_start) + int64(inode.I_block[12])*int64(blockSize)
		indirectBlock := make([]int32, blockSize/4)
		if err := Utilities.ReadObject(file, &indirectBlock, indirectBlockOffset); err != nil {
			return err
		}

		for i := 0; i < len(indirectBlock) && dataLen > 0; i++ {
			if indirectBlock[i] == -1 {
				newBlockIndex := AllocateNewBlock(superblock, file)
				if newBlockIndex == -1 {
					return fmt.Errorf("no hay más espacio en el bloque indirecto")
				}
				indirectBlock[i] = int32(newBlockIndex)
			}
			blockOffset := int64(superblock.S_block_start) + int64(indirectBlock[i])*int64(blockSize)
			writeSize := min(dataLen, blockSize)
			copyBlock := make([]byte, blockSize)
			copy(copyBlock, data[:writeSize])

			if err := Utilities.WriteObject(file, copyBlock, blockOffset); err != nil {
				return err
			}
			data = data[writeSize:]
			dataLen -= writeSize
		}
		if err := Utilities.WriteObject(file, indirectBlock, indirectBlockOffset); err != nil {
			return err
		}
	}

	// Actualizar tamaño del inodo
	inode.I_size = int32(len(newContent))
	if err := Utilities.WriteObject(file, inode, inodeStart); err != nil {
		return err
	}

	// Actualizar el superbloque en disco
	if err := Utilities.WriteObject(file, superblock, int64(mountedPartition.Start)); err != nil {
		return err
	}

	return nil
}
