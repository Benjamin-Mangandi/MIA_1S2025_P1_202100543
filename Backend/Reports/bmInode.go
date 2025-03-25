package Reports

import (
	"Backend/DiskManager"
	Disk "Backend/Structs/disk"
	Ext2 "Backend/Structs/ext2"
	"Backend/Utilities"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func CreateBmInodeReport(path string, id string) {
	// Buscar la partición montada
	path = fixPath(path)
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

	// Leer el MBR
	var mbr Disk.MBR
	if err := Utilities.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("Error al leer el MBR:", err)
		return
	}

	// Leer el superbloque
	var superblock Ext2.Superblock
	if err := Utilities.ReadObject(file, &superblock, int64(mountedPartition.Start)); err != nil {
		fmt.Println("Error al leer el superbloque:", err)
		return
	}

	// Utilizar S_inodes_count como total de inodos
	totalInodes := superblock.S_inodes_count

	// Calcular cuántos bytes necesita el bitmap (cada byte tiene 8 bits)
	byteCount := (totalInodes + 7) / 8

	// Variable para almacenar el contenido del reporte del bitmap de inodos
	var bitmapContent strings.Builder

	for byteIndex := int32(0); byteIndex < byteCount; byteIndex++ {
		// Mover el puntero al byte correspondiente en el bitmap de inodos
		_, err := file.Seek(int64(superblock.S_bm_inode_start+byteIndex), 0)
		if err != nil {
			fmt.Println("Error al mover el puntero en el archivo:", err)
			return
		}

		// Leer un byte del bitmap
		var byteVal byte
		err = binary.Read(file, binary.LittleEndian, &byteVal)
		if err != nil {
			fmt.Println("Error al leer el byte del bitmap:", err)
			return
		}

		// Procesar cada bit del byte (cada bit representa un inodo)
		for bitOffset := 0; bitOffset < 8; bitOffset++ {
			// Verificar si estamos fuera del rango total de inodos
			if byteIndex*8+int32(bitOffset) >= totalInodes {
				break
			}

			// Si el bit es 1, el inodo está ocupado, si es 0, está libre
			if (byteVal & (1 << bitOffset)) != 0 {
				bitmapContent.WriteByte('1') // Inodo ocupado
			} else {
				bitmapContent.WriteByte('0') // Inodo libre
			}

			// Añadir salto de línea cada 20 inodos
			if (byteIndex*8+int32(bitOffset)+1)%20 == 0 {
				bitmapContent.WriteString("\n")
			}
		}
	}

	// Guardar el reporte en el archivo especificado
	txtFile, err := os.Create(path)
	if err != nil {
		fmt.Println("Error al crear el archivo de reporte:", err)
		return
	}
	defer txtFile.Close()

	_, err = txtFile.WriteString(bitmapContent.String())
	if err != nil {
		fmt.Println("Error al escribir en el archivo de reporte:", err)
		return
	}

	fmt.Println("Reporte del bitmap de inodos generado correctamente:", path)
}
