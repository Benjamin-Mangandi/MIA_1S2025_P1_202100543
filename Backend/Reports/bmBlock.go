package Reports

import (
	"Backend/DiskManager"
	"Backend/Responsehandler"
	Disk "Backend/Structs/disk"
	Ext2 "Backend/Structs/ext2"
	"Backend/Utilities"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func CreateBmBlockReport(path string, id string) {
	// Buscar la partición montada
	path = fixPath(path)
	err := Utilities.CreateParentDirs(path)
	if err != nil {
		response := strings.Repeat("*", 30) + "\n" +
			"Error al crear las carpetas padre" + "\n"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
	mountedPartition := DiskManager.GetMountedPartitionByID(id)
	if mountedPartition.ID == "" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No se encontró la partición montada." + "\n"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
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

	totalBlocks := superblock.S_blocks_count
	// Calcular cuántos bytes necesita el bitmap (cada byte tiene 8 bits)
	byteCount := (totalBlocks + 7) / 8

	// Variable para almacenar el contenido del reporte del bitmap de bloques
	var bitmapContent strings.Builder

	for byteIndex := int32(0); byteIndex < byteCount; byteIndex++ {
		// Mover el puntero al byte correspondiente en el bitmap de bloques
		_, err := file.Seek(int64(superblock.S_bm_block_start+byteIndex), 0)
		if err != nil {
			fmt.Println("Error al mover el puntero en el archivo:", err)
			return
		}

		// Leer un byte del bitmap
		var byteVal byte
		if err = binary.Read(file, binary.LittleEndian, &byteVal); err != nil {
			fmt.Println("Error al leer un byte del bitmap:", err)
			return
		}

		// Procesar cada bit del byte (cada bit representa un bloque)
		for bitOffset := 0; bitOffset < 8; bitOffset++ {
			// Verificar si estamos fuera del rango total de bloques
			if byteIndex*8+int32(bitOffset) >= totalBlocks {
				break
			}

			// Si el bit es 1, el bloque está ocupado, si es 0, está libre
			if (byteVal & (1 << bitOffset)) != 0 {
				bitmapContent.WriteByte('1') // Bloque ocupado
			} else {
				bitmapContent.WriteByte('0') // Bloque libre
			}

			// Añadir salto de línea cada 20 bloques
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

	if _, err = txtFile.WriteString(bitmapContent.String()); err != nil {
		fmt.Println("Error al escribir en el archivo de reporte:", err)
		return
	}
	response := strings.Repeat("-", 40) + "\n" +
		"Reporte del bitmap de bloques generado correctamente:" + path + "\n"
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
}
