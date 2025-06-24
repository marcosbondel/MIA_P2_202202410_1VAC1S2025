// models/partition.go
package models

type PartitionInfo struct {
	Status string `json:"status"`
	Type   string `json:"type"`
	Fit    string `json:"fit"`
	Start  int32  `json:"start"`
	Size   int32  `json:"size"`
	Name   string `json:"name"`
}
