package Globals

import (
	Disk "Backend/Structs/disk"
	Ext2 "Backend/Structs/ext2"
	"Backend/Utilities"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"regexp"
)

var ActiveUser Ext2.User // Variable global para la sesión activa

var MountedPartitions = make(map[string][]Disk.MountedPartition)

// Expresión regular para capturar parámetros en el formato -key=value
var Regex = regexp.MustCompile(`-(\w+)(?:=("[^"]+"|\S+))?`)

func HasPermission(user Ext2.User, inode Ext2.Inode, permission byte) bool {
	// Verifica si el usuario es root (tiene acceso completo)
	if user.ID == "1" {
		return true
	}

	// Convertir permisos ASCII a números enteros
	ownerPerms := int(inode.I_perm[0] - '0')
	groupPerms := int(inode.I_perm[1] - '0')
	otherPerms := int(inode.I_perm[2] - '0')

	// Si el usuario es el dueño del inodo
	if user.ID == fmt.Sprint(inode.I_uid) {
		return (ownerPerms & int(permission)) != 0
	}

	// Si el usuario pertenece al grupo del inodo
	if user.Group == fmt.Sprint(inode.I_gid) {
		return (groupPerms & int(permission)) != 0
	}

	// Si es un usuario "otro"
	return (otherPerms & int(permission)) != 0
}

// ReadFileFromInode lee el contenido de un archivo dado su inodo
func ReadFileFromInode(file *os.File, superblock Ext2.Superblock, inodeIndex int, activeUser Ext2.User) string {
	var inode Ext2.Inode
	inodeOffset := int64(superblock.S_inode_start + int32(inodeIndex)*int32(binary.Size(Ext2.Inode{})))
	if err := Utilities.ReadObject(file, &inode, inodeOffset); err != nil {
		fmt.Println("Error al leer el inodo del archivo.")
		return ""
	}

	var buffer bytes.Buffer

	if !HasPermission(activeUser, inode, 4) {
		return ""
	}
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
