package Structs

type MountedPartition struct {
	Path   string
	Name   string
	ID     string
	Status byte // 0: no montada, 1: montada
	Start  int32
}
