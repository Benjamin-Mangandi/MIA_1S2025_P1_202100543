package Disk

import (
	"Backend/Responsehandler"
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
	answer := "---------------------\n" +
		"Disco creado correctamente\n" +
		"Tamaño: " + fmt.Sprintf("%d", data.Size) + " bytes" + "\n" +
		"Fecha de creación: " + string(data.CreationDate[:]) + "\n" +
		"Signature: " + fmt.Sprintf("%d", data.Signature) + "\n" +
		"Fit: " + string(data.Fit) + "\n" +
		"---------------------"
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, answer)
}
