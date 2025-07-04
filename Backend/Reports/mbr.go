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

func CreateMBR_Report(path string, id string) {
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

	// Iniciar la estructura Graphviz
	var dotContent string
	dotContent += "digraph MBR_Report {\n"
	dotContent += "node [shape=plaintext]\n"
	dotContent += "Reporte [label=<\n"
	dotContent += "<table border='1' cellborder='1' cellspacing='0'>\n"

	dotContent += "<tr><td colspan='2' bgcolor='#14cd06'><b>REPORTE DE MBR</b></td></tr>\n"
	dotContent += fmt.Sprintf("<tr><td>mbr_tamaño</td><td>%d</td></tr>\n", mbr.Size)
	creationDateStr := strings.TrimRight(string(mbr.CreationDate[:]), "\x00")
	dotContent += fmt.Sprintf("<tr><td>mbr_fecha_creacion</td><td>%s</td></tr>\n", creationDateStr)
	dotContent += fmt.Sprintf("<tr><td>mbr_disk_signature</td><td>%d</td></tr>\n", mbr.Signature)
	dotContent += fmt.Sprintf("<tr><td>mbr_disk_fit</td><td>%c</td></tr>\n", mbr.Fit)

	// Particiones
	for _, part := range mbr.Partitions {
		if part.Type == 'p' {
			dotContent += "<tr><td colspan='2' bgcolor='#10a5ca'><b>Particion</b></td></tr>\n"
			dotContent += fmt.Sprintf("<tr><td>part_status</td><td>%c</td></tr>\n", part.Status)
			dotContent += fmt.Sprintf("<tr><td>part_type</td><td>%c</td></tr>\n", part.Type)
			dotContent += fmt.Sprintf("<tr><td>part_fit</td><td>%c</td></tr>\n", part.Fit)
			dotContent += fmt.Sprintf("<tr><td>part_start</td><td>%d</td></tr>\n", part.Start)
			dotContent += fmt.Sprintf("<tr><td>part_size</td><td>%d</td></tr>\n", part.Size)
			partNameClean := strings.TrimRight(string(part.Name[:]), "\x00")
			dotContent += fmt.Sprintf("<tr><td>part_name</td><td>%s</td></tr>\n", partNameClean)
			dotContent += fmt.Sprintf("<tr><td>part_correlative</td><td>%d</td></tr>\n", part.Correlative)
		}
		// Si es extendida, mostrar particiones lógicas
		if part.Type == 'e' {
			dotContent += "<tr><td colspan='2' bgcolor='#10a5ca'><b>Particion</b></td></tr>\n"
			dotContent += fmt.Sprintf("<tr><td>part_status</td><td>%c</td></tr>\n", part.Status)
			dotContent += fmt.Sprintf("<tr><td>part_type</td><td>%c</td></tr>\n", part.Type)
			dotContent += fmt.Sprintf("<tr><td>part_fit</td><td>%c</td></tr>\n", part.Fit)
			dotContent += fmt.Sprintf("<tr><td>part_start</td><td>%d</td></tr>\n", part.Start)
			dotContent += fmt.Sprintf("<tr><td>part_size</td><td>%d</td></tr>\n", part.Size)
			partNameClean := strings.TrimRight(string(part.Name[:]), "\x00")
			dotContent += fmt.Sprintf("<tr><td>part_name</td><td>%s</td></tr>\n", partNameClean)
			dotContent += fmt.Sprintf("<tr><td>part_correlative</td><td>%d</td></tr>\n", part.Correlative)
			ebrStart := part.Start
			for ebrStart != -1 {
				var ebr Disk.EBR
				if err := Utilities.ReadObject(file, &ebr, int64(ebrStart)); err != nil {
					break
				}
				dotContent += "<tr><td colspan='2' bgcolor='#FF6600'><b>Particion Logica</b></td></tr>\n"
				dotContent += fmt.Sprintf("<tr><td>part_next</td><td>%d</td></tr>\n", ebr.PartNext)
				dotContent += fmt.Sprintf("<tr><td>part_fit</td><td>%c</td></tr>\n", ebr.PartFit)
				dotContent += fmt.Sprintf("<tr><td>part_start</td><td>%d</td></tr>\n", ebr.PartStart)
				dotContent += fmt.Sprintf("<tr><td>part_size</td><td>%d</td></tr>\n", ebr.PartSize)
				partNameClean := strings.TrimRight(string(ebr.PartName[:]), "\x00")
				dotContent += fmt.Sprintf("<tr><td>part_name</td><td>%s</td></tr>\n", partNameClean)
				ebrStart = ebr.PartNext
			}
		}
	}

	dotContent += " </table>>];\n"
	dotContent += "}\n"

	// Guardar el código Graphviz en un archivo temporal
	tempDotPath := "/home/user/mbr_report.dot"
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
		"Reporte del mbr generado exitosamente en: " + path + "\n"
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
}

func fixPath(path string) string {
	homeDir, _ := os.UserHomeDir() // Obtiene el home del usuario actual
	return strings.Replace(path, "/home/user/", homeDir+"/", 1)
}
