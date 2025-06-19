package disk_management

import (
	"MIA_P2_202202410_1VAC1S2025/fs/structs"
	"MIA_P2_202202410_1VAC1S2025/fs/utils"
	"encoding/binary"
	"fmt"
	"strings"
)

func Mount(driveLetter string, name string) {
	fmt.Println("======Start MOUNT======")
	fmt.Println("Drive Letter:", driveLetter)
	fmt.Println("Name:", name)

	// Open bin file
	filepath := "./test/" + strings.ToUpper(driveLetter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	var tempMBR structs.MRB
	// Read MRB from file
	if err := utils.ReadObject(file, &tempMBR, 0); err != nil {
		fmt.Println("Error reading MRB from file:", err)
		return
	}

	// Print MRB to verify
	structs.PrintMBR(tempMBR)

	// Iterate through partitions to find the one with the given name
	var index int = -1
	var count int = 1
	var emptyId [4]byte
	for i := 0; i < 4; i++ {
		// prunt id
		if strings.Contains(string(tempMBR.Partitions[i].Name[:]), name) {
			if tempMBR.Partitions[i].Id != emptyId {
				fmt.Println("Error: Partition with name", name, "already mounted")
				return
			}
			index = i
		}
		if tempMBR.Partitions[i].Id != emptyId {
			count++
		}
	}

	if index != -1 {
		fmt.Println("Partition found:")
		structs.PrintPartition(tempMBR.Partitions[index])
	} else {
		fmt.Println("Error: Partition with name", name, "not found")
		return
	}

	// id = DriveLetter + Correlative + 19
	// My ID as Student "202202410'
	id := strings.ToUpper(driveLetter) + fmt.Sprintf("%d", count) + "10"

	// edit the partition to set the id
	copy(tempMBR.Partitions[index].Id[:], id)      // Set id
	copy(tempMBR.Partitions[index].Status[:], "1") // Set status to 1 (active)

	// Overwrite the MRB in the file
	if err := utils.WriteObject(file, tempMBR, 0); err != nil {
		fmt.Println("Error writing MRB to file:", err)
		return
	}

	// Print the updated partition
	fmt.Println("Updated Partition:")
	structs.PrintPartition(tempMBR.Partitions[index])

	// Close the bin file
	defer file.Close()

	fmt.Println("======End MOUNT======")
}

func Unmount(id string) {
	fmt.Println("======Start UNMOUNT======")
	fmt.Println("ID:", id)

	// Open bin file
	filepath := "./test/" + strings.ToUpper(id[:1]) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	var tempMBR structs.MRB
	// Read MRB from file
	if err := utils.ReadObject(file, &tempMBR, 0); err != nil {
		fmt.Println("Error reading MRB from file:", err)
		return
	}

	// Print MRB to verify
	structs.PrintMBR(tempMBR)

	// Iterate through partitions to find the one with the given id
	var found bool = false
	for i := 0; i < 4; i++ {
		if strings.Contains(string(tempMBR.Partitions[i].Id[:]), id) {
			found = true
			copy(tempMBR.Partitions[i].Status[:], "0") // Set status to 0 (inactive)
			break
		}
	}

	if !found {
		fmt.Println("Error: Partition with id", id, "not found")
		return
	}

	// Overwrite the MRB in the file
	if err := utils.WriteObject(file, tempMBR, 0); err != nil {
		fmt.Println("Error writing MRB to file:", err)
		return
	}

	fmt.Println("Partition with id", id, "unmounted successfully.")
	structs.PrintMBR(tempMBR)
	// Close the bin file
	defer file.Close()

	fmt.Println("Unmount operation completed successfully.")
	fmt.Println("======End UNMOUNT======")
}

func Fdisk(size int, driveLetter string, name string, type_ string, fit string, unit string, add int, delete string) {
	fmt.Println("======Start FDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Drive Letter:", driveLetter)
	fmt.Println("Name:", name)
	fmt.Println("Type:", type_)
	fmt.Println("Fit:", fit)
	fmt.Println("Unit:", unit)

	if fit == "" {
		fit = "wf" // Default fit if not provided
	}

	// validate fit equal to b/f/w
	if fit != "bf" && fit != "ff" && fit != "wf" {
		fmt.Println("Error: Fit must be b, f, or w")
		return
	}

	// validate size greater than 0
	if size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	// validate unit equal to k/m
	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: Unit must be b or k or m")
		return
	}

	// set the size in bytes
	if unit == "k" {
		size *= 1024
	} else if unit == "m" {
		size *= 1024 * 1024
	}

	// Open bin file
	filepath := "./test/" + strings.ToUpper(driveLetter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	var tempMBR structs.MRB
	// Read MRB from file
	if err := utils.ReadObject(file, &tempMBR, 0); err != nil {
		fmt.Println("Error reading MRB from file:", err)
		return
	}

	var gap = int32(0)
	// Iterate through partitions to calculate the gaps
	for i := 0; i < 4; i++ {
		if tempMBR.Partitions[i].Size != 0 {
			gap = tempMBR.Partitions[i].Start + tempMBR.Partitions[i].Size
		}
	}

	// Iterate through partitions to find an empty one
	var foundEmpty bool
	for i := 0; i < 4; i++ {
		if tempMBR.Partitions[i].Size == 0 {
			foundEmpty = true
			// Create a new partition
			tempMBR.Partitions[i].Size = int32(size)   // Set size
			copy(tempMBR.Partitions[i].Name[:], name)  // Set name
			copy(tempMBR.Partitions[i].Fit[:], fit)    // Set fit
			copy(tempMBR.Partitions[i].Status[:], "0") // Set status to 1 (active)
			copy(tempMBR.Partitions[i].Type[:], type_) // Set type
			tempMBR.Partitions[i].Start = 0            // Set start to 0 (for simplicity)

			if gap > 0 {
				tempMBR.Partitions[i].Start = gap
			} else {
				tempMBR.Partitions[i].Start = int32(binary.Size(tempMBR))
			}
			break
		}
	}

	// Print MRB to verify
	structs.PrintMBR(tempMBR)

	if !foundEmpty {
		fmt.Println("Error: No empty partition found")
		return
	}

	// Overwrite the MRB in the file
	if err := utils.WriteObject(file, tempMBR, 0); err != nil {
		fmt.Println("Error writing MRB to file:", err)
		return
	}

	// DELETE LOGIC
	if delete == "full" {
		for i := 0; i < 4; i++ {
			part := &tempMBR.Partitions[i]
			if strings.TrimSpace(string(part.Name[:])) == name {
				fmt.Println("Deleting partition:", name)
				part.Size = 0
				copy(part.Name[:], "")
				copy(part.Id[:], "")
				copy(part.Status[:], "0")
				copy(part.Type[:], "")
				copy(part.Fit[:], "")
				break
			}
		}
		utils.WriteObject(file, tempMBR, 0)
		fmt.Println("Partition deleted successfully")
		defer file.Close()
		return
	}

	// ADD SIZE LOGIC
	if add != 0 {
		for i := 0; i < 4; i++ {
			part := &tempMBR.Partitions[i]
			if strings.TrimSpace(string(part.Name[:])) == name {
				fmt.Println("Modifying partition:", name)

				// Calcula el cambio en bytes
				addBytes := add
				if unit == "k" {
					addBytes *= 1024
				} else if unit == "m" {
					addBytes *= 1024 * 1024
				}

				newSize := int(part.Size) + addBytes

				if newSize <= 0 {
					fmt.Println("Error: Size would be zero or negative")
					defer file.Close()
					return
				}

				// Si estás aumentando, verifica que no se solape con la siguiente partición
				if addBytes > 0 {
					end := part.Start + part.Size + int32(addBytes)
					for j := 0; j < 4; j++ {
						if i == j || tempMBR.Partitions[j].Size == 0 {
							continue
						}
						if end > tempMBR.Partitions[j].Start {
							fmt.Println("Error: Resize would overlap with another partition")
							defer file.Close()
							return
						}
					}
				}

				part.Size = int32(newSize)
				utils.WriteObject(file, tempMBR, 0)
				fmt.Println("Partition size updated successfully")
				defer file.Close()
				return
			}
		}
		fmt.Println("Error: Partition not found")
		defer file.Close()
		return
	}

	// Close the bin file
	defer file.Close()

	fmt.Println("======End FDISK======")
}

func Mkdisk(size int, fit string, unit string) {
	fmt.Println("======Start MKDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Fit:", fit)
	fmt.Println("Unit:", unit)

	// Validate input
	if fit == "" {
		fit = "ff" // Default fit if not provided
	}

	if fit != "bf" && fit != "ff" && fit != "wf" {
		fmt.Println("Error: Fit must be bf, ff, or wf")
		return
	}

	if size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}
	if unit != "k" && unit != "m" {
		fmt.Println("Error: Unit must be k or m")
		return
	}

	// Find available disk name
	diskLetter := utils.FindAvailableLetter("./test")
	if diskLetter == "" {
		fmt.Println("Error: No more available letters for disks")
		return
	}
	diskPath := fmt.Sprintf("./test/%s.bin", diskLetter)

	// Create file
	err := utils.CreateFile(diskPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	// Calculate byte size
	if unit == "k" {
		size *= 1024
	} else {
		size *= 1024 * 1024
	}

	// Open bin file
	file, err := utils.OpenFile(diskPath)
	if err != nil {
		return
	}

	// Write empty bytes
	zeroBuffer := make([]byte, 1024)
	for i := 0; i < size/1024; i++ {
		err := utils.WriteObject(file, zeroBuffer, int64(i*1024))
		if err != nil {
			return
		}
	}

	// Create and write MBR
	var newMRB structs.MRB
	newMRB.MbrSize = int32(size)
	newMRB.Signature = 10
	copy(newMRB.Fit[:], fit)
	copy(newMRB.CreationDate[:], "2025-05-01")

	if err := utils.WriteObject(file, newMRB, 0); err != nil {
		fmt.Println("Error writing MRB to file:", err)
		return
	}

	// Read and print MBR
	var tempMBR structs.MRB
	if err := utils.ReadObject(file, &tempMBR, 0); err != nil {
		fmt.Println("Error reading MRB from file:", err)
		return
	}

	fmt.Println("Disk created:", diskPath)
	fmt.Println("MRB:", tempMBR)
	fmt.Println("File size:", size)
	fmt.Println("Fit:", string(tempMBR.Fit[:]))
	fmt.Println("Creation date:", string(tempMBR.CreationDate[:]))
	fmt.Println("Signature:", tempMBR.Signature)

	defer file.Close()
	fmt.Println("======End MKDISK======")
}

func Rmdisk(driveLetter string) {
	fmt.Println("======Start RMDISK======")
	fmt.Println("Drive letter:", driveLetter)

	// Validate drive letter
	if len(driveLetter) != 1 || driveLetter[0] < 'A' || driveLetter[0] > 'Z' {
		fmt.Println("Error: Drive letter must be a single uppercase letter from A to Z")
		return
	}

	// Construct the file path
	filePath := fmt.Sprintf("./test/%s.bin", driveLetter)

	// Delete the file
	if err := utils.DeleteFile(filePath); err != nil {
		fmt.Println("Error deleting file:", err)
		return
	}

	fmt.Println("Disk removed successfully.")
	fmt.Println("======End RMDISK======")
}
