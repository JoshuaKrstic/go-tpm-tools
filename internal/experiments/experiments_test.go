package experiments

import (
	"bytes"
	"testing"
)

func TestGetBooleanExperiment(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("{\"TestFeatureForImage\":true}")

	e, err := readJSONInput(&buffer)
	if err != nil {
		t.Errorf("Failed to parse input json, err: %v", err)
	}

	val, ok := GetBooleanExperiment(e, "TestFeatureForImage")
	if !ok {
		t.Error("failed to get key 'TestFeatureForImage'")
	}
	if val != true {
		t.Errorf("expected 'TestFeatureForImage' to be true, got: %v", val)
	}
}

func TestGetBooleanExperimentNoKey(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("{}")

	e, err := readJSONInput(&buffer)
	if err != nil {
		t.Errorf("Failed to parse input json, err: %v", err)
	}

	val, ok := GetBooleanExperiment(e, "InvalidKey")
	if ok {
		t.Errorf("expected to fail to get boolean experiment, got %v", val)
	}
}

func TestGetBooleanExperimentInvalidType(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("{\"StringFeature\":\"string_value\"}")

	e, err := readJSONInput(&buffer)
	if err != nil {
		t.Errorf("Failed to parse input json, err: %v", err)
	}

	val, ok := GetBooleanExperiment(e, "StringFeature")
	if ok {
		t.Errorf("expected to fail to get boolean experiment, got %v", val)
	}
}

func TestGetStringExperiment(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("{\"StringFeature\":\"string_value\"}")

	e, err := readJSONInput(&buffer)
	if err != nil {
		t.Errorf("Failed to parse input json, err: %v", err)
	}

	val, ok := GetStringExperiment(e, "StringFeature")
	if !ok {
		t.Error("failed to get key 'StringFeature'")
	}
	if val != "string_value" {
		t.Errorf("expected 'StringFeature' to be 'string_value', got: %v", val)
	}
}

func TestGetStringExperimentNoKey(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("{}")

	e, err := readJSONInput(&buffer)
	if err != nil {
		t.Errorf("Failed to parse input json, err: %v", err)
	}

	val, ok := GetBooleanExperiment(e, "InvalidKey")
	if ok {
		t.Errorf("expected to fail to get string experiment, got %v", val)
	}
}

func TestGetStringExperimentInvalidType(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("{\"FloatFeature\":-5.6}")

	e, err := readJSONInput(&buffer)
	if err != nil {
		t.Errorf("Failed to parse input json, err: %v", err)
	}

	val, ok := GetStringExperiment(e, "FloatFeature")
	if ok {
		t.Errorf("expected to fail to get string experiment, got %v", val)
	}
}

func TestGetFloatExperiment(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("{\"FloatFeature\":-5.6}")

	e, err := readJSONInput(&buffer)
	if err != nil {
		t.Errorf("Failed to parse input json, err: %v", err)
	}

	val, ok := GetFloatExperiment(e, "FloatFeature")
	if !ok {
		t.Error("failed to get key 'FloatFeature'")
	}
	if val != -5.6 {
		t.Errorf("expected 'FloatFeature' to be -5.6, got: %v", val)
	}
}

func TestGetFloatExperimentNoKey(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("{}")

	e, err := readJSONInput(&buffer)
	if err != nil {
		t.Errorf("Failed to parse input json, err: %v", err)
	}

	val, ok := GetFloatExperiment(e, "InvalidKey")
	if ok {
		t.Errorf("expected to fail to get string experiment, got %v", val)
	}
}

func TestGetFloatExperimentInvalidType(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("{\"TestFeatureForImage\":true}")

	e, err := readJSONInput(&buffer)
	if err != nil {
		t.Errorf("Failed to parse input json, err: %v", err)
	}

	val, ok := GetFloatExperiment(e, "TestFeatureForImage")
	if ok {
		t.Errorf("expected to fail to get float experiment, got %v", val)
	}
}

func TestGetIntExperiment(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("{\"IntFeature\":10293812}")

	e, err := readJSONInput(&buffer)
	if err != nil {
		t.Errorf("Failed to parse input json, err: %v", err)
	}

	val, ok := GetIntExperiment(e, "IntFeature")
	if !ok {
		t.Error("failed to get key 'IntFeature'")
	}
	if val != 10293812 {
		t.Errorf("expected 'IntFeature' to be 10293812, got: %v", val)
	}
}

func TestGetIntExperimentNoKey(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("{}")

	e, err := readJSONInput(&buffer)
	if err != nil {
		t.Errorf("Failed to parse input json, err: %v", err)
	}

	val, ok := GetIntExperiment(e, "InvalidKey")
	if ok {
		t.Errorf("expected to fail to get int experiment, got %v", val)
	}
}

func TestGetIntExperimentInvalidType(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("{\"StringFeature\":\"string_value\"}")

	e, err := readJSONInput(&buffer)
	if err != nil {
		t.Errorf("Failed to parse input json, err: %v", err)
	}

	val, ok := GetIntExperiment(e, "StringFeature")
	if ok {
		t.Errorf("expected to fail to get int experiment, got %v", val)
	}
}

func TestGetSeveralExperiments(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("{\"TestFeatureForImage\":true,\"StringFeature\":\"string_value\",\"FloatFeature\":-5.6,\"OtherTestFeatureForImage\":false}")

	e, err := readJSONInput(&buffer)
	if err != nil {
		t.Errorf("Failed to parse input json, err: %v", err)
	}

	val, ok := GetBooleanExperiment(e, "TestFeatureForImage")
	if !ok {
		t.Error("failed to get key 'TestFeatureForImage'")
	}
	if val != true {
		t.Errorf("expected 'TestFeatureForImage' to be true, got: %v", val)
	}

	val, ok = GetBooleanExperiment(e, "OtherTestFeatureForImage")
	if !ok {
		t.Error("failed to get key 'TestFeatureForImage'")
	}
	if val != false {
		t.Errorf("expected 'OtherTestFeatureForImage' to be false, got: %v", val)
	}
}

func TestReadExperimentsFileReturnsEmptyOnNoFile(t *testing.T) {
	e, err := ReadExperimentsFile("$!bad path//")

	if err == nil {
		t.Errorf("expected err to be non nil, got nil")
	}

	if e == nil || len(e) > 0 {
		t.Errorf("expected empty experiments object, got %v", e)
	}
}

func TestReadExperimentsFileReturnsEmptyOnBadJsonFormat(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("{\"this string isnt json")

	e, err := readJSONInput(&buffer)

	if err == nil {
		t.Errorf("expected err to be non nil, got nil")
	}

	if e == nil || len(e) > 0 {
		t.Errorf("expected empty experiments object, got %v", e)
	}
}
