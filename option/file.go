package option

import (
	"encoding/json"
	"os"
)

func ReadJson(filename string) (Option, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var option = Option{}
	err = json.NewDecoder(f).Decode(&option)
	if err != nil {
		return nil, err
	}

	return option, nil
}
