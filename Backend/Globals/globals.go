package Globals

import (
	"Backend/Structs"
	Ext2 "Backend/Structs/ext2"
	"regexp"
)

var ActiveUser Ext2.User // Variable global para la sesión activa

var MountedPartitions = make(map[string][]Structs.MountedPartition)

var Regex = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
