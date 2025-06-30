// models/partition.go
package models

type PartitionInfo struct {
	Status         string `json:"status"`
	Type           string `json:"type"`
	Fit            string `json:"fit"`
	Start          int32  `json:"start"`
	Size           int32  `json:"size"`
	Name           string `json:"name"`
	Id             string `json:"id,omitempty"`
	Mounted        bool   `json:"mounted"`
	IsLogical      bool   `json:"is_logical,omitempty"`
	PartNext       int32  `json:"part_next,omitempty"`       // solo para lógicas
	EBRStart       int32  `json:"ebr_start,omitempty"`       // solo para lógicas
	PartitionIndex int    `json:"partition_index,omitempty"` // para primarias/ext
}
