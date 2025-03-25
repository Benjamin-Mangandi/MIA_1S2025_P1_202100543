package Ext2

import "fmt"

type Inode struct {
	I_uid   int32
	I_gid   int32
	I_size  int32
	I_atime [16]byte
	I_ctime [16]byte
	I_mtime [16]byte
	I_block [15]int32
	I_type  [1]byte
	I_perm  [3]byte
}

func PrintInode(inode Inode) {
	fmt.Println("====== Inode ======")
	fmt.Printf("I_uid: %d\n", inode.I_uid)
	fmt.Printf("I_gid: %d\n", inode.I_gid)
	fmt.Printf("I_size: %d\n", inode.I_size)
	fmt.Printf("I_atime: %s\n", string(inode.I_atime[:]))
	fmt.Printf("I_ctime: %s\n", string(inode.I_ctime[:]))
	fmt.Printf("I_mtime: %s\n", string(inode.I_mtime[:]))
	fmt.Printf("I_type: %s\n", string(inode.I_type[:]))
	fmt.Printf("I_perm: %s\n", string(inode.I_perm[:]))
	fmt.Printf("I_block: %v\n", inode.I_block)
	fmt.Println("===================")
}
