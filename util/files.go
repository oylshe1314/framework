package util

import (
	"fmt"
	"os"
)

func ReadFiles(dir string) ([]string, error) {
	des, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var files []string
	for _, de := range des {
		if de.IsDir() {
			subFiles, err := ReadFiles(dir + "/" + de.Name())
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, dir+"/"+de.Name())
		}
	}
	return files, nil
}
