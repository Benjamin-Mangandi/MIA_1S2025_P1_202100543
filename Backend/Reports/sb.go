package Reports

import (
	"Backend/DiskManager"
	Ext2 "Backend/Structs/ext2"
	"Backend/Utilities"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func CreateSbReport(path string, id string) {
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

	fmt.Println(superblock.S_mtime)
	fmt.Println(superblock.S_umtime)
	mtimeStr := strings.Trim(string(superblock.S_mtime[:]), "\x00")
	var mtime string
	if mtimeStr == "" {
		mtime = "N/A"
	} else {
		parsedMTime, err := time.Parse("2006-01-02 15:04", mtimeStr)
		if err != nil {
			fmt.Println("Error al parsear mtime:", err)
			return
		}
		mtime = parsedMTime.Format(time.RFC3339)
	}

	umtimeStr := strings.Trim(string(superblock.S_umtime[:]), "\x00")
	var umtime string
	if umtimeStr == "" {
		umtime = "N/A"
	} else {
		parsedUMTime, err := time.Parse("2006-01-02 15:04", umtimeStr)
		if err != nil {
			fmt.Println("Error al parsear umtime:", err)
			return
		}
		umtime = parsedUMTime.Format(time.RFC3339)
	}

	dotFormat := `
digraph G {
	fontname="Helvetica,Arial,sans-serif"
	node [fontname="Helvetica,Arial,sans-serif", shape=plain, fontsize=12];
	edge [fontname="Helvetica,Arial,sans-serif", color="#FF7043", arrowsize=0.8];
	bgcolor="#FAFAFA";
	rankdir=TB;

	superblockTable [label=<
		<table border="0" cellborder="1" cellspacing="0" cellpadding="10" bgcolor="#FFF9C4" style="rounded">
			<tr><td colspan="2" bgcolor="#4CAF50" align="center"><b>REPORTE DEL SUPERBLOQUE</b></td></tr>
			<tr><td><b>Tipo de Sistema de Archivos</b></td><td>%d (EXT2)</td></tr>
			<tr><td><b>Cantidad de Inodos</b></td><td>%d</td></tr>
			<tr><td><b>Cantidad de Bloques</b></td><td>%d</td></tr>
			<tr><td><b>Inodos Libres</b></td><td>%d</td></tr>
			<tr><td><b>Bloques Libres</b></td><td>%d</td></tr>
			<tr><td><b>Montajes Realizados</b></td><td>%d</td></tr>
			<tr><td><b>Última Modificación</b></td><td>%s</td></tr>
			<tr><td><b>Último Montaje</b></td><td>%s</td></tr>
			<tr><td><b>Identificador Mágico</b></td><td>0x%X</td></tr>
			<tr><td><b>Tamaño de Inodo</b></td><td>%d bytes</td></tr>
			<tr><td><b>Tamaño de Bloque</b></td><td>%d bytes</td></tr>
			<tr><td><b>Primer Inodo Libre</b></td><td>%d</td></tr>
			<tr><td><b>Primer Bloque Libre</b></td><td>%d</td></tr>
			<tr><td><b>Inicio Bitmap de Inodos</b></td><td>%d</td></tr>
			<tr><td><b>Inicio Bitmap de Bloques</b></td><td>%d</td></tr>
			<tr><td><b>Inicio Tabla de Inodos</b></td><td>%d</td></tr>
			<tr><td><b>Inicio Bloques de Datos</b></td><td>%d</td></tr>
		</table>>];
}
`

	// Formatear el contenido con los datos del superbloque
	dotContent := fmt.Sprintf(dotFormat,
		superblock.S_filesystem_type,
		superblock.S_inodes_count,
		superblock.S_blocks_count,
		superblock.S_free_inodes_count,
		superblock.S_free_blocks_count,
		superblock.S_mnt_count,
		mtime,
		umtime,
		superblock.S_magic,
		superblock.S_inode_size,
		superblock.S_block_size,
		superblock.S_fist_ino,
		superblock.S_first_blo,
		superblock.S_bm_inode_start,
		superblock.S_bm_block_start,
		superblock.S_inode_start,
		superblock.S_block_start,
	)

	// Guardar el código Graphviz en un archivo temporal
	tempDotPath := "/home/benjamin/inode_report.dot"
	if err := os.WriteFile(tempDotPath, []byte(dotContent), 0644); err != nil {
		fmt.Println("Error al escribir el archivo .dot:", err)
		return
	}

	// Ejecutar Graphviz para generar la imagen
	cmd := exec.Command("dot", "-Tjpg", tempDotPath, "-o", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error al ejecutar Graphviz:", err)
		fmt.Println("Salida del comando:", string(output))
		return
	}

	fmt.Println("Reporte del superbloque generado exitosamente en:", path)
}
