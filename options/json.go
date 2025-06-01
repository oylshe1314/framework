package options

import (
	json "github.com/json-iterator/go"
	"os"
)

func ReadJson(filename string) (Options, error) {
	var options = map[string]interface{}{} //Remind me, don't use Options{} to init the variable
	if len(filename) == 0 {
		return options, nil
	}

	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, &options)
	if err != nil {
		return nil, err
	}

	return options, nil
}
