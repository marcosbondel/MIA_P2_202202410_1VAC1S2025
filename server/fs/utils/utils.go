package utils

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Funtion to create bin file
func CreateFile(name string) error {
	//Ensure the directory exists
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	// Create the file
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return err
		}
		defer file.Close()
	}
	return nil
}

// Funtion to open bin file in read/write mode
func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Errpr open file:", err)
		return nil, err
	}
	return file, nil
}

// Function to write and object to a bin file
func WriteObject(file *os.File, data interface{}, position int64) error {
	_, err := file.Seek(position, 0) // mejor usar io.SeekStart
	if err != nil {
		fmt.Println("Error seeking file:", err)
		return err
	}

	if err := binary.Write(file, binary.LittleEndian, data); err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	return nil
}

// Function to read an object from a bin file
func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Error reading the object", err)
		return err
	}
	return nil
}

func DeleteFile(name string) error {
	if err := os.Remove(name); err != nil {
		fmt.Println("Error deleting file:", err)
		return err
	}
	fmt.Println("File deleted successfully:", name)
	return nil
}

func FindAvailableLetter(directory string) string {
	entries, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return ""
	}

	usedLetters := make(map[string]bool)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".bin") {
			letter := strings.ToUpper(string(entry.Name()[0]))
			usedLetters[letter] = true
		}
	}

	for i := 'A'; i <= 'Z'; i++ {
		letter := string(i)
		if !usedLetters[letter] {
			return letter
		}
	}
	return ""
}

// Busca el primer bit libre en un bitmap (retorna el índice o -1 si no hay)
func FindFreeBit(file *os.File, bitmapStart int32, totalBits int) int {
	file.Seek(int64(bitmapStart), 0)

	bitmap := make([]byte, totalBits)
	if _, err := file.Read(bitmap); err != nil {
		fmt.Println("Error leyendo bitmap:", err)
		return -1
	}

	for i := 0; i < totalBits; i++ {
		if bitmap[i] == 0 {
			return i
		}
	}
	return -1
}

// Escribe un bit en el bitmap (valor debe ser 0 o 1)
func WriteBit(file *os.File, bitmapStart int32, index int, value byte) error {
	pos := int64(bitmapStart) + int64(index)

	file.Seek(pos, 0)
	if _, err := file.Write([]byte{value}); err != nil {
		fmt.Println("Error escribiendo bit:", err)
		return err
	}
	return nil
}
