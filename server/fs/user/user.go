package user

import (
	//  "os"
	"MIA_P1_202202410_1VAC1S2025/global"
	"MIA_P1_202202410_1VAC1S2025/structs"
	"MIA_P1_202202410_1VAC1S2025/utils"
	"MIA_P1_202202410_1VAC1S2025/utils_inodes"
	"encoding/binary"
	"fmt"
	"strings"
)

// // login -user=root -pass=123 -id=A119
func Login(user string, pass string, id string) {
	fmt.Println("======Start LOGIN======")
	fmt.Println("User:", user)
	fmt.Println("Pass:", pass)
	fmt.Println("Id:", id)

	if global.CurrentUser.Status {
		fmt.Println("User already logged in")
		return
	}

	var login bool = false
	driveletter := string(id[0])

	// Open bin file
	filepath := "./test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR structs.MRB
	// Read object from bin file
	if err := utils.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	structs.PrintMBR(TempMBR)

	fmt.Println("-------------")

	var index int = -1
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) {
				fmt.Println("Partition found")
				if strings.Contains(string(TempMBR.Partitions[i].Status[:]), "1") {
					fmt.Println("Partition is mounted")
					index = i
				} else {
					fmt.Println("Partition is not mounted")
					return
				}
				break
			}
		}
	}

	if index != -1 {
		structs.PrintPartition(TempMBR.Partitions[index])
	} else {
		fmt.Println("Partition not found")
		return
	}

	var tempSuperblock structs.Superblock
	// Read object from bin file
	if err := utils.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		return
	}

	// initSearch /users.txt -> regresa no Inodo
	// initSearch -> 1
	indexInode := utils_inodes.InitSearch("/users.txt", file, tempSuperblock)

	// indexInode := int32(1)

	var crrInode structs.Inode
	// Read object from bin file
	if err := utils.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		return
	}

	// read file data
	data := utils_inodes.GetInodeFileData(crrInode, file, tempSuperblock)

	fmt.Println("Fileblock------------")
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	// login -user=root -pass=123 -id=A119

	// Iterar a través de las líneas
	for _, line := range lines {
		// Imprimir cada línea
		// fmt.Println(line)
		words := strings.Split(line, ",")
		fmt.Println("Words:", words)
		if len(words) == 5 {
			if (strings.Contains(words[3], user)) && (strings.Contains(words[4], pass)) {
				login = true

				break
			}
		}
	}

	// Print object
	fmt.Println("Inode", crrInode.I_block)

	// Close bin file
	defer file.Close()

	if login {
		fmt.Println("User logged in")
		global.CurrentUser.ID = id
		global.CurrentUser.Status = true
		global.CurrentUser.User = user
	} else {
		fmt.Println("User not found or invalid credentials")
	}

	fmt.Println("======End LOGIN======")
}

func Logout() {
	fmt.Println("======Start LOGOUT======")
	if global.CurrentUser.Status {
		global.CurrentUser.ID = ""
		global.CurrentUser.Status = false
		fmt.Println("User logged out")
	} else {
		fmt.Println("No user logged in")
	}
	fmt.Println("======End LOGOUT======")
}
