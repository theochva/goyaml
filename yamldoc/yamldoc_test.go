package yamldoc

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	_SampleYaml = `
a:
  b:
    c: value-c
  d:
    e: false
    f: 10
`
)

var (
	yamlText = strings.TrimSpace(_SampleYaml)
	yamlMap  = map[interface{}]interface{}{
		"a": map[interface{}]interface{}{
			"b": map[interface{}]interface{}{
				"c": "value-c",
			},
			"d": map[interface{}]interface{}{
				"e": false,
				"f": 10,
			},
		},
	}
	yamlStringMap = map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "value-c",
			},
			"d": map[string]interface{}{
				"e": false,
				"f": 10,
			},
		},
	}
)

func TestYamlDoc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "YamlDoc Test Suite")
}

func checkContainsValue(yaml YamlDoc, key string, expectedValue bool) {
	contains, err := yaml.Contains(key)
	Expect(err).To(BeNil())
	Expect(contains).To(Equal(expectedValue))
}
func checkGetValue(yaml YamlDoc, key string, expectedValue interface{}) {
	value, err := yaml.Get(key)
	Expect(err).To(BeNil())
	if expectedValue == nil {
		Expect(value).To(BeNil())
	} else {
		Expect(value).To(Equal(expectedValue))
	}
}
func checkGetStringValue(yaml YamlDoc, key string, expectedValue string) {
	value, err := yaml.GetString(key)
	Expect(err).To(BeNil())
	Expect(value).To(Equal(expectedValue))
}
func checkGetStringValueTypeErr(yaml YamlDoc, key string) {
	_, err := yaml.GetString(key)
	Expect(err).ToNot(BeNil())
	Expect(IsWrongTypeError(err)).To(BeTrue())
}
func checkGetIntValue(yaml YamlDoc, key string, expectedValue int) {
	value, err := yaml.GetInt(key)
	Expect(err).To(BeNil())
	Expect(value).To(Equal(expectedValue))
}
func checkGetIntValueTypeErr(yaml YamlDoc, key string) {
	_, err := yaml.GetInt(key)
	Expect(err).ToNot(BeNil())
	Expect(IsWrongTypeError(err)).To(BeTrue())
}
func checkGetBoolValue(yaml YamlDoc, key string, expectedValue bool) {
	value, err := yaml.GetBool(key)
	Expect(err).To(BeNil())
	if expectedValue {
		Expect(value).To(BeTrue())
	} else {
		Expect(value).To(BeFalse())
	}
}
func checkGetBoolValueTypeErr(yaml YamlDoc, key string) {
	_, err := yaml.GetBool(key)
	Expect(err).ToNot(BeNil())
	Expect(IsWrongTypeError(err)).To(BeTrue())
}
func checkDeleteValue(yaml YamlDoc, key string, valueExisted bool) {
	deleted, err := yaml.Delete(key)
	Expect(err).To(BeNil())
	Expect(deleted).To(Equal(valueExisted))

	checkGetValue(yaml, key, nil)
	checkContainsValue(yaml, key, false)
}
func checkSetValue(yaml YamlDoc, key string, value interface{}) {
	valueSet, err := yaml.Set(key, value)
	Expect(err).To(BeNil())
	Expect(valueSet).To(BeTrue())

	checkGetValue(yaml, key, value)
}
func checkText(yaml YamlDoc, expectedText string) {
	text, err := yaml.Text()
	Expect(err).To(BeNil())
	Expect(text).To(Equal(expectedText))
}
func checkSampleYaml(yaml YamlDoc) {
	// Check data map has one key
	Expect(len(yaml.Data())).To(Equal(1))
	// Check that it contains the single entry
	checkContainsValue(yaml, "a.b.c", true)
	// Check when serialized it produces what we expect
	checkText(yaml, yamlText)
	// Check we get a nested element
	checkGetValue(yaml, "a.b.c", "value-c")

	// Check we get a nested string element
	checkGetStringValue(yaml, "a.b.c", "value-c")
	// Check we get an wrong type error when getting string value
	checkGetStringValueTypeErr(yaml, "a.d.e")
	// Check we get an wrong type error when getting int value
	checkGetIntValueTypeErr(yaml, "a.b.c")
	// Check we get an wrong type error when getting bool value
	checkGetBoolValueTypeErr(yaml, "a.d.f")
	// Check we get a nested string element
	checkGetBoolValue(yaml, "a.d.e", false)
	// Check we get a nested string element
	checkGetIntValue(yaml, "a.d.f", 10)
	// Check it returns nil when entry not found
	checkGetValue(yaml, "foo.bar", nil)
	// Check that it returns false when deleting a non-existant entry
	checkDeleteValue(yaml, "foo.bar", false)
	// Check that it returns false when deleting a non-existant entry
	checkDeleteValue(yaml, "foo", false)
	// Check we can add/remove a child entry
	newKey := "g"
	newValue := "value-g"
	newText := fmt.Sprintf("%s\n%s: %s", yamlText, newKey, newValue)
	checkSetValue(yaml, newKey, newValue)
	checkText(yaml, newText)

	// Remove the new key
	checkDeleteValue(yaml, newKey, true)

	// After deleting the entry, verify what
	// is left is what we started with
	checkText(yaml, yamlText)
}

var _ = Describe("Yaml functions", func() {
	var yaml YamlDoc

	Context("Read existing string text", func() {
		BeforeEach(func() {
			var err error
			yaml, err = FromString(yamlText)
			Expect(err).To(BeNil())
			Expect(yaml).ToNot(BeNil())
		})

		It("Check sample yaml", func() {
			checkSampleYaml(yaml)
		})
	})
	Context("Start with empty then set values from map", func() {
		BeforeEach(func() {
			var err error
			yaml, err = New(nil)
			Expect(err).To(BeNil())
			Expect(yaml).ToNot(BeNil())
		})
		// Check the data in the yaml is the same as the map we initialized with
		JustBeforeEach(func() {
			yaml.SetData(yamlMap)
		})

		It("is the same as the original map", func() {
			Expect(reflect.DeepEqual(yamlMap, yaml.Data())).To(BeTrue())
		})

		It("Check sample yaml", func() {
			checkSampleYaml(yaml)
		})

		It("can convert to map[string]interface{}", func() {
			// Get the data as a map[string]interface{} is the same as the one we initialized with
			converted, strYamlMap := yaml.Map()
			Expect(converted).To(BeTrue())
			Expect(reflect.DeepEqual(strYamlMap, yamlStringMap)).To(BeTrue())
		})
	})
})
