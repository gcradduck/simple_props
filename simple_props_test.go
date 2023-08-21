package simple_props

import (
	"testing"
)

const PROPERTY_COUNT = 16

type propsTest struct {
	key            string
	expectedString string
	expectedInt    int
	expectedBool   bool
}

// TestLoadProps tests LoadProps function
func TestLoadProps(t *testing.T) {

	props, gotErr := LoadProps("test_properties.properties")

	if gotErr != nil {
		t.Error("got error when calling LoadProps", gotErr)
	}

	expectedPropCount := PROPERTY_COUNT
	gotPropCount := len(props.Props)

	if gotPropCount != expectedPropCount {
		t.Errorf("expected property count of %d but got %d", expectedPropCount, gotPropCount)
	}

	_, gotErr = LoadProps(("some_nonexistent_filename.properties"))

	if gotErr == nil {
		t.Error("load of nonexistent properties file should have thrown an error but did not")
	}

}

// TestGet tests the Get function
func TestGet(t *testing.T) {

	props, gotErr := LoadProps("test_properties.properties")

	if gotErr != nil {
		t.Error("got error when calling LoadProps", gotErr)
	}

	tests := []propsTest{
		{key: "stringKey1", expectedString: "1"},
		{key: "stringKey2", expectedString: "two"},
		{key: "stringKey3", expectedString: `"three"`},
	}

	for _, test := range tests {

		gotValue := props.Get(test.key)

		if gotValue != test.expectedString {
			t.Errorf("key %s expected string value of %s but got %s", test.key, test.expectedString, gotValue)
		}
	}
}

// TestGetInt tests the GetInt function
func TestGetInt(t *testing.T) {

	props, gotErr := LoadProps("test_properties.properties")
	if gotErr != nil {
		t.Error("got error when calling LoadProps", gotErr)
	}

	test := propsTest{
		key:         "intKey1",
		expectedInt: 1,
	}

	gotValue := props.GetInt(test.key, 0)

	if gotValue != test.expectedInt {
		t.Errorf("key %s expected value %d but got %d", "intKey1", test.expectedInt, gotValue)
	}

	key := "someNonExistentKey"
	expectedValue := -1
	gotValue = props.GetInt(key, -1)

	if gotValue != expectedValue {
		t.Errorf("key %s expected value %d but got %d", key, expectedValue, gotValue)
	}

	key = "intKey2"
	expectedValue = 1
	gotValue = props.GetInt(key, 0)

	if gotValue == expectedValue {
		t.Errorf("key %s should have returned %d but got %d", key, 0, expectedValue)
	}

}

// TestGetCleanFilePath tests the GetCleanFilePath function
func TestGetCleanFilePath(t *testing.T) {

	props, gotErr := LoadProps("test_properties.properties")
	if gotErr != nil {
		t.Error("got error when calling LoadProps", gotErr)
	}

	tests := []propsTest{
		{key: "filePathKey1", expectedString: `C:\dir1\dir2`},
		{key: "filePathKey2", expectedString: `C:\dir1\dir2`},
		{key: "filePathKey3", expectedString: `C:\dir1\dir2\dir3\dir4\dir5\`},
		{key: "filePathKey4", expectedString: ``},
	}

	for _, test := range tests {

		expectedValue := test.expectedString
		gotValue := props.GetCleanFilePath(test.key)

		if gotValue != expectedValue {
			t.Errorf("key %s expected filepath %s but got %s", test.key, expectedValue, gotValue)
		}
	}
}

// TestGetBool tests the GetBool function
func TestGetBool(t *testing.T) {

	props, gotErr := LoadProps("test_properties.properties")
	if gotErr != nil {
		t.Error("got error when calling LoadProps", gotErr)
	}

	tests := []propsTest{
		{key: "boolKey1", expectedBool: true},
		{key: "boolKey2", expectedBool: true},
		{key: "boolKey3", expectedBool: true},
		{key: "boolKey4", expectedBool: true},
		{key: "boolKey5", expectedBool: true},

		{key: "boolKey6", expectedBool: false},
		{key: "boolKey7", expectedBool: false},
		{key: "boolKey8", expectedBool: false},
	}

	for _, test := range tests {

		expectedValue := test.expectedBool
		gotValue := props.GetBool(test.key, false)

		if gotValue != expectedValue {
			t.Errorf("key %s expected bool %v but got %v", test.key, expectedValue, gotValue)
		}
	}

	key := "nonexistentboolkey"
	expectedValue := true
	gotValue := props.GetBool(key, true)

	if gotValue != expectedValue {
		t.Errorf("key %s expected bool %v but got %v", "nonexistentboolkey", expectedValue, gotValue)
	}
}
