package Structs

import (
	"fmt"
)

type MBR struct {
	Size         int32
	CreationDate [16]byte
	Signature    int32
	Fit          byte
	Partitions   [4]Partition
}

func PrintMBR(data MBR) {

	fmt.Println(fmt.Sprintf("Fecha de Creacion: %s, Fit: %s, Size: %d, Signature: %d ",
		string(data.CreationDate[:]),
		string(data.Fit),
		data.Size,
		data.Signature))

}
