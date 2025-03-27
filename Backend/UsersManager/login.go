package UsersManager

import (
	"Backend/DiskManager"
	"Backend/Globals"
	"Backend/Responsehandler"
	Ext2 "Backend/Structs/ext2"
	"Backend/Utilities"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func Login(user, pass, id string) bool {

	if Globals.ActiveUser.Status {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: Se debe cerrar la sesion anterior primero"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
	}
	// 1. Buscar la partición montada por ID
	mountedPartition := DiskManager.GetMountedPartitionByID(id)

	// Verificar si la partición es inválida (ID vacío significa que no se encontró)
	if mountedPartition.ID == "" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: Partición no encontrada o no montada."
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return false
	}

	// 2. Abrir el archivo binario de la partición
	file, err := Utilities.OpenFile(mountedPartition.Path)
	if err != nil {
		fmt.Println("Error al abrir el archivo del disco.")
		return false
	}
	defer file.Close()

	// 3. Leer el superbloque desde la partición
	var superblock Ext2.Superblock
	if err := Utilities.ReadObject(file, &superblock, int64(mountedPartition.Start)); err != nil {
		fmt.Println("Error al leer el superbloque.")
		return false
	}

	// 4. Buscar el inodo del archivo users.txt
	inodeIndex := FindFileInode(file, superblock, "/users.txt")
	if inodeIndex == -1 {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: Archivo users.txt no encontrado."
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return false
	}

	// 5. Leer el contenido del archivo users.txt
	fileContent := ReadFileFromInode(file, superblock, inodeIndex)
	if fileContent == "" {
		fmt.Println("Error: No se pudo leer el archivo users.txt.")
		return false
	}

	// 6. Verificar las credenciales del usuario
	lines := strings.Split(fileContent, "\n")
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 5 && parts[1] == "U" { // Verifica que es un usuario
			userID, userType, group, username, password := parts[0], parts[1], parts[2], parts[3], parts[4]

			// Verificar si el usuario está eliminado (ID en 0)
			if userID == "0" {
				continue // Saltar usuarios eliminados
			}
			if username == user && password == pass {
				// Guardar usuario en la sesión activa
				Globals.ActiveUser = Ext2.User{
					ID:          userID,
					Type:        userType,
					Group:       group,
					Name:        username,
					Password:    password,
					Status:      true,
					PartitionID: mountedPartition.ID,
				}
				Ext2.PrintUser(Globals.ActiveUser)
				return true
			}
		}
	}
	response := strings.Repeat("*", 30) + "\n" +
		"Error: Usuario o contraseña incorrectos o usuario eliminado."
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
	return false
}

func FindFileInode(file *os.File, superblock Ext2.Superblock, filePath string) int {
	// Normalizar la ruta
	filePath = strings.TrimLeft(filePath, "/") // Elimina el "/" inicial
	pathParts := strings.Split(filePath, "/")  // Divide la ruta en partes

	// Comenzar en el inodo raíz
	var currentInode Ext2.Inode
	inodeIndex := 0 // Inodo raíz

	// Iterar sobre cada parte del path
	for _, part := range pathParts {
		// Leer el inodo actual
		inodeOffset := int64(superblock.S_inode_start) + int64(inodeIndex)*int64(binary.Size(Ext2.Inode{}))
		if err := Utilities.ReadObject(file, &currentInode, inodeOffset); err != nil {
			fmt.Println("Error al leer inodo en la ruta:", part)
			return -1
		}

		// Buscar la siguiente parte en los bloques de datos del inodo actual
		found := false
		for _, blockIndex := range currentInode.I_block {
			if blockIndex == -1 {
				continue
			}

			// Leer el bloque de carpetas
			var folderBlock Ext2.Folderblock
			blockOffset := int64(superblock.S_block_start + int32(blockIndex)*int32(binary.Size(Ext2.Folderblock{})))
			if err := Utilities.ReadObject(file, &folderBlock, blockOffset); err != nil {
				fmt.Println("Error al leer el bloque de carpetas.")
				continue
			}

			// Buscar el nombre dentro del bloque
			for _, entry := range folderBlock.B_content {
				entryName := strings.TrimRight(string(entry.B_name[:]), "\x00") // Convertir a string
				if entryName == part {
					// Avanzar al siguiente inodo
					inodeIndex = int(entry.B_inodo)
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			response := strings.Repeat("*", 30) + "\n" +
				"Archivo no encontrado en el sistema de archivos:" + filePath
			Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
			return -1
		}
	}

	return inodeIndex
}

// ReadFileFromInode lee el contenido de un archivo dado su inodo
func ReadFileFromInode(file *os.File, superblock Ext2.Superblock, inodeIndex int) string {
	var inode Ext2.Inode
	inodeOffset := int64(superblock.S_inode_start + int32(inodeIndex)*int32(binary.Size(Ext2.Inode{})))
	if err := Utilities.ReadObject(file, &inode, inodeOffset); err != nil {
		fmt.Println("Error al leer el inodo del archivo.")
		return ""
	}

	var buffer bytes.Buffer

	// Leer los bloques de datos asignados al inodo
	for _, blockIndex := range inode.I_block {
		if blockIndex == -1 {
			continue
		}

		// Leer el bloque de datos
		var fileBlock Ext2.Fileblock
		blockOffset := int64(superblock.S_block_start + int32(blockIndex)*int32(binary.Size(Ext2.Fileblock{})))
		if err := Utilities.ReadObject(file, &fileBlock, blockOffset); err != nil {
			fmt.Println("Error al leer el bloque de archivo.")
			continue
		}

		// Convertir los datos del bloque a texto y agregarlos al buffer
		buffer.Write(bytes.Trim(fileBlock.B_content[:], "\x00")) // Eliminar bytes nulos
	}

	return buffer.String()
}
