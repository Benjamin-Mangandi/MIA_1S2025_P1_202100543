package Reports

import (
	"Backend/DiskManager"
	Ext2 "Backend/Structs/ext2"
	"Backend/Utilities"
	"fmt"
	"os"
	"os/exec"
	"unsafe"
)

func CreateInode_Report(path string, id string) {
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

	// Leer el superbloque
	var superblock Ext2.Superblock
	if err := Utilities.ReadObject(file, &superblock, int64(mountedPartition.Start)); err != nil {
		fmt.Println("Error al leer el superbloque:", err)
		return
	}

	// Iniciar la estructura Graphviz
	var dotContent string
	dotContent += "digraph Inode_Report {\n"
	dotContent += "node [shape=plaintext]\n"

	// Iterar sobre los inodos
	inodeStart := int64(superblock.S_inode_start)
	inodeSize := int64(unsafe.Sizeof(Ext2.Inode{}))

	for i := 0; i < int(superblock.S_inodes_count); i++ {
		inodePos := inodeStart + (inodeSize * int64(i))
		var inode Ext2.Inode
		if err := Utilities.ReadObject(file, &inode, inodePos); err != nil {
			fmt.Println("Error al leer el inodo:", err)
			return
		}

		if inode.I_uid == 0 {
			continue
		}
		dotContent += fmt.Sprintf("Inodo%d [label=<\n", i+1)
		dotContent += "<table border='1' cellborder='1' cellspacing='0'>\n"
		dotContent += fmt.Sprintf("<tr><td colspan='2'><b>Inodo %d</b></td></tr>\n", i+1)
		dotContent += fmt.Sprintf("<tr><td>i_uid</td><td>%d</td></tr>\n", inode.I_uid)
		dotContent += fmt.Sprintf("<tr><td>i_size</td><td>%d</td></tr>\n", inode.I_size)

		for j, block := range inode.I_block {
			if block != -1 {
				dotContent += fmt.Sprintf("<tr><td>i_block_%d</td><td>%d</td></tr>\n", j+1, block)
			}
		}

		dotContent += fmt.Sprintf("<tr><td>i_perm</td><td>%d</td></tr>\n", inode.I_perm[:])
		dotContent += "</table>>];\n"
	}
	dotContent += "}\n"

	// Guardar el código Graphviz en un archivo temporal
	tempDotPath := "/home/benjamin/inode_report.dot"
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

	fmt.Println("Reporte de inodos generado exitosamente en:", path)
}
