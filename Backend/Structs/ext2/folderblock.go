package Ext2

import "fmt"

type Folderblock struct {
	B_content [4]Content
}

func PrintFolderblock(folderblock Folderblock) {
	fmt.Println("====== Folderblock ======")
	for i, content := range folderblock.B_content {
		fmt.Printf("Content %d: Name: %s, Inodo: %d\n", i, string(content.B_name[:]), content.B_inodo)
	}
	fmt.Println("=========================")
}
