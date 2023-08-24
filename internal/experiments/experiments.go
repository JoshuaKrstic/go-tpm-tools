// Package experiments contains functionalities to retrieve synced experiments
package experiments

import (
	"encoding/json"
	"io"
	"os"
)

// GetBooleanExperiment attempts to read the key from the map provided and does type checking on the
// returned value, if any. Value is false on failure.
func GetBooleanExperiment(e map[string]any, key string) (bool, bool) {
	v, err := e[key]
	if !err {
		return false, false
	}

	f, ok := v.(bool)
	if !ok {
		return false, false
	}

	return f, true
}

// GetStringExperiment attempts to read the key from the map provided and does type checking on the
// returned value, if any. Value is an empty string on failure.
func GetStringExperiment(e map[string]any, key string) (string, bool) {
	v, err := e[key]
	if !err {
		return "", false
	}

	f, ok := v.(string)
	if !ok {
		return "", false
	}

	return f, true
}

// GetIntExperiment attempts to read the key from the map provided and does type checking on the
// returned value, if any. Value is -1 on failure.
// encoding/JSON doesn't differentiate between ints and floats. This method casts the float from
// the unmarshalled json to an int.
func GetIntExperiment(e map[string]any, key string) (int, bool) {
	v, err := e[key]
	if !err {
		return -1, false
	}

	f, ok := v.(float64)
	if !ok {
		return -1, false
	}

	return int(f), true
}

// GetFloatExperiment attempts to read the key from the map provided and does type checking on the
// returned value, if any. Value is -1 on failure.
func GetFloatExperiment(e map[string]any, key string) (float64, bool) {
	v, err := e[key]
	if !err {
		return -1, false
	}

	f, ok := v.(float64)
	if !ok {
		return -1, false
	}

	return f, true
}

// ReadExperimentsFile takes a filepath, opens the file, and calls ReadJsonInput with the contents
// of the file.
// If the file cannot be opened, the experiments map is set to an empty map.
func ReadExperimentsFile(fpath string) (map[string]any, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return map[string]any{}, err
	}

	return readJSONInput(f)
}

// ReadJSONInput  takes a reader and unmarshalls the contents into the experiments map.
// If the unmarsahlling fails, the experiments map is set to an empty map.
func readJSONInput(read io.Reader) (map[string]any, error) {
	jsondata := make(map[string]any)
	if err := json.NewDecoder(read).Decode(&jsondata); err != nil {
		return map[string]any{}, err
	}

	return jsondata, nil
}
