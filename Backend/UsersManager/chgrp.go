package UsersManager

import (
	"Backend/DiskManager"
	"Backend/Globals"
	"Backend/Responsehandler"
	Ext2 "Backend/Structs/ext2"
	"Backend/Utilities"
	"fmt"
	"strings"
)

func Chgrp(user, group string) {
	// Verificar si hay una sesión activa
	if !Globals.ActiveUser.Status {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No hay un usuario activo." + "\n"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	if Globals.ActiveUser.Name != "root" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: Solo el usuario root puede eliminar usuarios" + "\n"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Obtener la partición montada asociada a la sesión
	mountedPartition := DiskManager.GetMountedPartitionByID(Globals.ActiveUser.PartitionID)
	if mountedPartition.ID == "" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No se encontró la partición montada." + "\n"
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
			"Error: No se pudo leer el archivo users.txt." + "\n"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	// Verificar si el usuario existe y obtener su línea
	lines := strings.Split(fileContent, "\n")
	userFound := false
	groupExists := false

	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 3 && parts[1] == "G" && parts[2] == group {
			groupExists = true
		}
		if len(parts) == 5 && parts[1] == "U" && parts[3] == user {
			userFound = true
		}
	}

	if !userFound {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: El usuario especificado no existe." + "\n"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	if !groupExists {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: El grupo especificado no existe."
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Modificar la línea del usuario para cambiar su grupo sin perder la contraseña
	for i, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 5 && parts[1] == "U" && parts[3] == user {
			lines[i] = fmt.Sprintf("%s,U,%s,%s,%s", parts[0], user, group, parts[4]) // Se mantiene la contraseña
			break
		}
	}

	// Unir el contenido actualizado y escribirlo de vuelta en users.txt
	newContent := strings.Join(lines, "\n")
	if err := WriteFileToInode(file, &superblock, inodeIndex, newContent, mountedPartition); err != nil {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No se pudo actualizar el archivo users.txt."
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	response := strings.Repeat("*", 30) + "\n" +
		"El grupo del usuario ha sido actualizado correctamente: " + user + "\n"
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)

}
