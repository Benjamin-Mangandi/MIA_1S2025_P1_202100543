package Ext2

import "fmt"

type Superblock struct {
	S_filesystem_type   int32    // Guarda el número que identifica el sistema de archivos utilizado
	S_inodes_count      int32    // Guarda el número total de inodos
	S_blocks_count      int32    // Guarda el número total de bloques
	S_free_blocks_count int32    // Contiene el número de bloques libres
	S_free_inodes_count int32    // Contiene el número de inodos libres
	S_mtime             [17]byte // Última fecha en el que el sistema fue montado
	S_umtime            [17]byte // Última fecha en que el sistema fue desmontado
	S_mnt_count         int32    // Indica cuantas veces se ha montado el sistema
	S_magic             int32    // Valor que identifica al sistema de archivos, tendrá el valor 0xEF53
	S_inode_size        int32    // Tamaño del inodo
	S_block_size        int32    // Tamaño del bloque
	S_fist_ino          int32    // Primer inodo libre (dirección del inodo)
	S_first_blo         int32    // Primer bloque libre (dirección del inodo)
	S_bm_inode_start    int32    // Guardará el inicio del bitmap de inodos
	S_bm_block_start    int32    // Guardará el inicio del bitmap de bloques
	S_inode_start       int32    // Guardará el inicio de la tabla de inodos
	S_block_start       int32    // Guardará el inicio de la tabla de bloques
}

func PrintSuperblock(sb Superblock) {
	fmt.Println("====== Superblock ======")
	fmt.Printf("S_filesystem_type: %d\n", sb.S_filesystem_type)
	fmt.Printf("S_inodes_count: %d\n", sb.S_inodes_count)
	fmt.Printf("S_blocks_count: %d\n", sb.S_blocks_count)
	fmt.Printf("S_free_blocks_count: %d\n", sb.S_free_blocks_count)
	fmt.Printf("S_free_inodes_count: %d\n", sb.S_free_inodes_count)
	fmt.Printf("S_mtime: %s\n", string(sb.S_mtime[:]))
	fmt.Printf("S_umtime: %s\n", string(sb.S_umtime[:]))
	fmt.Printf("S_mnt_count: %d\n", sb.S_mnt_count)
	fmt.Printf("S_magic: 0x%X\n", sb.S_magic) // Usamos 0x%X para mostrarlo en formato hexadecimal
	fmt.Printf("S_inode_size: %d\n", sb.S_inode_size)
	fmt.Printf("S_block_size: %d\n", sb.S_block_size)
	fmt.Printf("S_fist_ino: %d\n", sb.S_fist_ino)
	fmt.Printf("S_first_blo: %d\n", sb.S_first_blo)
	fmt.Printf("S_bm_inode_start: %d\n", sb.S_bm_inode_start)
	fmt.Printf("S_bm_block_start: %d\n", sb.S_bm_block_start)
	fmt.Printf("S_inode_start: %d\n", sb.S_inode_start)
	fmt.Printf("S_block_start: %d\n", sb.S_block_start)
	fmt.Println("========================")
}
