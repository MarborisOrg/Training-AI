package utils

import (
	"fmt"
	"math"
	"os"
	"os/user"
	"path/filepath"
)

func MultipliesByTwo(x float64) float64 {
	return 2 * x
}

func SubtractsOne(x float64) float64 {
	return x - 1
}

func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func Index(slice []string, text string) int {
	for i, item := range slice {
		if item == text {
			return i
		}
	}

	return 0
}

func Difference(slice []string, slice2 []string) (difference []string) {
	for i := 0; i < 2; i++ {
		for _, s1 := range slice {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found {
				difference = append(difference, s1)
			}
		}
		if i == 0 {
			slice, slice2 = slice2, slice
		}
	}
	return difference
}

func ReadFile(path string) (bytes []byte) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		bytes, err = os.ReadFile("../" + path)
	}

	if err != nil {
		panic(err)
	}

	return bytes
}

func Contains(slice []string, text string) bool {
	for _, item := range slice {
		if item == text {
			return true
		}
	}

	return false
}

func GetResDir(dir string, file string, dir2 ...string) (filePath string) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		panic("[Training] No user.. (419)")
	}
	homeDir := usr.HomeDir

	if len(dir2) == 0 || dir2[0] == "" {
		return filepath.Join(homeDir, ".marboris", "res", dir, file)
	}

	return filepath.Join(homeDir, ".marboris", "res", dir, dir2[0], file)
}
