package Globals

import (
	Disk "Backend/Structs/disk"
	Ext2 "Backend/Structs/ext2"
	"regexp"
)

var ActiveUser Ext2.User // Variable global para la sesión activa

var MountedPartitions = make(map[string][]Disk.MountedPartition)

// Expresión regular para capturar parámetros en el formato -key=value
var Regex = regexp.MustCompile(`-(\w+)(?:=("[^"]+"|\S+))?`)
