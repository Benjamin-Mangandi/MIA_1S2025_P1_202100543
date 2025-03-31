package Reports

import (
	"Backend/DiskManager"
	"Backend/Responsehandler"
	Ext2 "Backend/Structs/ext2"
	"Backend/Utilities"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CreateBlocksReport(path string, id string) {
	// Buscar la partición montada
	path = fixPath(path)
	err := Utilities.CreateParentDirs(path)
	if err != nil {
		fmt.Println("Error al crear directorios:", err)
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

	dotContent := "digraph G {\n"
	dotContent += "node [shape=record];\n"
	dotContent += "rankdir=TB;\n" // Mejor disposición para bloques

	// Mapa para almacenar información de inodos
	inodeTypes := make(map[int32]byte) // Almacena el tipo de inodo (0: carpeta, 1: archivo)

	// Primera pasada: recolectar tipos de inodos
	for i := int32(0); i < superblock.S_inodes_count; i++ {
		inodePos := superblock.S_inode_start + (i * superblock.S_inode_size)
		var inode Ext2.Inode
		if err := Utilities.ReadObject(file, &inode, int64(inodePos)); err == nil {
			inodeTypes[i] = inode.I_type[0]
		}
	}

	// Segunda pasada: procesar bloques
	for i := int32(0); i < superblock.S_blocks_count; i++ {
		blockPos := superblock.S_block_start + (i * superblock.S_block_size)

		// Verificar si el bloque está en uso (consultando bitmap)
		if !isBlockUsed(file, superblock.S_bm_block_start, i) {
			continue
		}

		// Determinar el tipo de bloque consultando los inodos que lo referencian
		isFolder := false
		isFile := false

		// Buscar inodos que apunten a este bloque
		for inodeNum, inodeType := range inodeTypes {
			inodePos := superblock.S_inode_start + (inodeNum * superblock.S_inode_size)
			var inode Ext2.Inode
			if err := Utilities.ReadObject(file, &inode, int64(inodePos)); err != nil {
				continue
			}

			for _, blockPtr := range inode.I_block {
				if blockPtr == i {
					if inodeType == '0' {
						isFolder = true
					} else if inodeType == '1' {
						isFile = true
					}
					break
				}
			}
		}

		// Procesar según el tipo detectado
		if isFolder {
			var folderBlock Ext2.Folderblock
			if err := Utilities.ReadObject(file, &folderBlock, int64(blockPos)); err == nil {
				dotContent += fmt.Sprintf("block%d [label=\"{Bloque Carpeta %d|", i, i)
				for _, content := range folderBlock.B_content {
					name := bytes.TrimRight(content.B_name[:], "\x00")
					if content.B_inodo != -1 && string(name) != "-1" {
						name := strings.TrimRight(string(content.B_name[:]), "\x00")
						dotContent += fmt.Sprintf("<f%d> %s : %d|", content.B_inodo, name, content.B_inodo)
					}
				}
				dotContent += "}\"]\n"
			}
		} else if isFile {
			var fileBlock Ext2.Fileblock
			if err := Utilities.ReadObject(file, &fileBlock, int64(blockPos)); err == nil {
				content := strings.TrimRight(string(fileBlock.B_content[:]), "\x00")
				// Escapar caracteres especiales para DOT
				content = strings.ReplaceAll(content, "\"", "\\\"")
				content = strings.ReplaceAll(content, "\n", "\\n")
				content = strings.ReplaceAll(content, "\r", "\\r")
				dotContent += fmt.Sprintf("block%d [label=\"{Bloque Archivo %d|%s}\"]\n", i, i, content)
			}
		} else {
			// Podría ser un bloque de punteros
			var pointerBlock Ext2.Pointerblock
			if err := Utilities.ReadObject(file, &pointerBlock, int64(blockPos)); err == nil {
				hasPointers := false
				for _, ptr := range pointerBlock.B_pointers {
					if ptr != -1 {
						hasPointers = true
						break
					}
				}
				if hasPointers {
					dotContent += fmt.Sprintf("block%d [label=\"{Bloque Punteros %d|", i, i)
					for j, ptr := range pointerBlock.B_pointers {
						if ptr != -1 {
							dotContent += fmt.Sprintf("<f%d> %d|", j, ptr)
						}
					}
					dotContent += "}\"]\n"
				}
			}
		}
	}

	// Tercera pasada: agregar relaciones
	for inodeNum, _ := range inodeTypes {
		inodePos := superblock.S_inode_start + (inodeNum * superblock.S_inode_size)
		var inode Ext2.Inode
		if err := Utilities.ReadObject(file, &inode, int64(inodePos)); err != nil {
			continue
		}

		for _, blockPtr := range inode.I_block {
			if blockPtr != -1 && blockPtr != 0 {
				dotContent += fmt.Sprintf("block%d -> block%d;\n", blockPtr-1, inodeNum)
			}
		}
	}

	dotContent += "}" // Cerrar el DOT

	// Guardar el código Graphviz en un archivo temporal
	tempDotPath := "/home/user/block_report.dot"
	tempDotPath = fixPath(tempDotPath)
	if err := os.WriteFile(tempDotPath, []byte(dotContent), 0644); err != nil {
		fmt.Println("Error al escribir el archivo .dot:", err)
		return
	}

	cmd := exec.Command("dot", "-Tjpg", tempDotPath, "-o", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error al ejecutar Graphviz:", err)
		fmt.Println("Salida del comando:", string(output))
		return
	}

	response := strings.Repeat("-", 40) + "\n" +
		"Reporte de bloques generado exitosamente en: " + path + "\n"
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
}

// Función auxiliar para verificar si un bloque está en uso
func isBlockUsed(file *os.File, bmStart int32, blockNum int32) bool {
	bytePosition := bmStart + (blockNum / 8)
	bitPosition := blockNum % 8

	var byteVal byte
	_, err := file.Seek(int64(bytePosition), 0)
	if err != nil {
		return false
	}
	if err := binary.Read(file, binary.LittleEndian, &byteVal); err != nil {
		return false
	}

	return (byteVal & (1 << bitPosition)) != 0
}
