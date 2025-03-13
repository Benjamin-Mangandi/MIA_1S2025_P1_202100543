package Disk

// MountedPartition representa una partición montada en el sistema.
type MountedPartition struct {
	Path   string // Ruta del archivo del disco donde se encuentra la partición.
	Name   string // Nombre de la partición montada.
	ID     string // Identificador único de la partición montada.
	Status byte   // Estado de la partición: 0 = no montada, 1 = montada.
	Start  int32  // Posición de inicio de la partición dentro del archivo del disco.
}
