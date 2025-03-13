package Ext2

import "fmt"

type Fileblock struct {
	B_content [64]byte
}

func PrintFileblock(fileblock Fileblock) {
	fmt.Println("====== Fileblock ======")
	fmt.Printf("B_content: %s\n", string(fileblock.B_content[:]))
	fmt.Println("=======================")
}
