package Reports

import (
	"Backend/DiskManager"
	"Backend/Responsehandler"
	Disk "Backend/Structs/disk"
	"Backend/Utilities"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CreateDiskReport(path string, id string) {
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

	dotContent := `digraph G {
		fontname="Helvetica,Arial,sans-serif"
		node [fontname="Helvetica,Arial,sans-serif"]
		edge [fontname="Helvetica,Arial,sans-serif"]
		concentrate=True;
		rankdir=TB;
		node [shape=record];
		title [label="Reporte DISK" shape=plaintext fontname="Helvetica,Arial,sans-serif"];
		dsk [label="`

	// Calcular el tamaño total del disco y el tamaño usado
	totalSize := mbr.Size
	usedSize := int32(0)

	// Agregar MBR al reporte
	dotContent += "{MBR}"

	// Recorrer las particiones del MBR y generar el contenido DOT
	for _, part := range mbr.Partitions {
		if part.Size > 0 {
			// Calcular el porcentaje de uso
			percentage := (float64(part.Size) / float64(totalSize)) * 100
			usedSize += part.Size

			// Convertir Part_name a string y eliminar los caracteres nulos
			partName := strings.TrimRight(string(part.Name[:]), "\x00")
			if part.Type == 'p' {
				// Partición primaria
				dotContent += fmt.Sprintf("|{Primaria\\n%s\\n%.2f%%}", partName, percentage)
			} else if part.Type == 'e' {
				// Partición extendida
				dotContent += fmt.Sprintf("|{Extendida\\n%.2f%%|{", percentage)
				ebrStart := part.Start
				ebrCount := 0
				ebrUsedSize := int32(0)
				var ebr Disk.EBR
				for ebrStart != -1 {
					if err := Utilities.ReadObject(file, &ebr, int64(ebrStart)); err != nil {
						fmt.Println("Error al leer el EBR:", err)
						return
					}
					ebrName := strings.TrimRight(string(ebr.PartName[:]), "\x00")
					ebrPercentage := (float64(ebr.PartSize) / float64(totalSize)) * 100
					ebrUsedSize += ebr.PartSize

					// Agregar EBR y partición lógica al reporte
					if ebrCount > 0 {
						dotContent += "|"
					}
					dotContent += fmt.Sprintf("{EBR|Lógica\\n %s\\n%.2f%%}", ebrName, ebrPercentage)

					// Actualizar el inicio para el próximo EBR
					ebrStart = ebr.PartNext
					ebrCount++
				}

				// Calcular espacio libre dentro de la partición extendida
				extendedFreeSize := part.Size - ebrUsedSize
				if extendedFreeSize > 0 {
					extendedFreePercentage := (float64(extendedFreeSize) / float64(totalSize)) * 100
					dotContent += fmt.Sprintf("|Libre\\n %.2f%%", extendedFreePercentage)
				}

				dotContent += "}}"
			}
		}
	}

	// Calcular espacio libre restante y añadirlo si es necesario
	freeSize := totalSize - int64(usedSize)
	if freeSize > 0 {
		freePercentage := (float64(freeSize) / float64(totalSize)) * 100
		dotContent += fmt.Sprintf("|Libre\\n %.2f%%", freePercentage)
	}

	// Cerrar el nodo de disco y completar el DOT
	dotContent += `"];

		title -> dsk [style=invis];
	}`

	// Guardar el código Graphviz en un archivo temporal
	tempDotPath := "/home/user/disk_report.dot"
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

	response := strings.Repeat("*", 30) + "\n" +
		"Reporte generado exitosamente en: " + path + "\n"
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
}
