package models

type DiskResponse struct {
	Disks []DiskInfo `json:"disks"`
}

type DiskInfo struct {
	Name       string          `json:"name"`
	Size       int32           `json:"size"`
	Creation   string          `json:"creation"`
	Fit        string          `json:"fit"`
	Signature  int32           `json:"signature"`
	Partitions []PartitionInfo `json:"partitions"`
}
