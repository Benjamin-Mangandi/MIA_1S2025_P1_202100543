package Disk

import (
	"fmt"
)

type MBR struct {
	Size         int32        //Tamaño del disco
	CreationDate [16]byte     //Fecha de creacion del disco
	Signature    int32        //Numero random que lo identifica
	Fit          byte         //Tipo de ajuste: B, F o w
	Partitions   [4]Partition //4 Posibles Particiones
}

func PrintMBR(data MBR) {

	fmt.Println(fmt.Sprintf("Fecha de Creacion: %s, Fit: %s, Size: %d, Signature: %d ",
		string(data.CreationDate[:]),
		string(data.Fit),
		data.Size,
		data.Signature))

}
