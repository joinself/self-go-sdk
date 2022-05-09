package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type sourcesSpec struct {
	Sources map[string][]string
}

func constNames(body []byte) (string, error) {
	var spec sourcesSpec

	err := json.Unmarshal(body, &spec)
	if err != nil {
		return "", err
	}

	output := "\n"
	sources := map[string]struct{}{}
	facts := map[string]struct{}{}
	relations := "var spec = map[string][]string{\n"
	for key, specFacts := range spec.Sources {
		source := snakeCaseToCamelCase("Source_" + key)
		relations += fmt.Sprintf("\t%s: []string{\n", source)
		if _, ok := sources[key]; !ok {
			output += fmt.Sprintf("const %s = \"%s\"", source, key) + "\n"
			sources[key] = struct{}{}
		}

		for _, f := range specFacts {
			fact := snakeCaseToCamelCase("Fact_" + f)
			relations += fmt.Sprintf("\t\t%s,\n", fact)
			if _, ok := facts[f]; !ok {
				output += fmt.Sprintf("const %s = \"%s\"", fact, f) + "\n"
				facts[f] = struct{}{}
			}
		}
		relations += fmt.Sprintf("\t},\n")
	}
	relations += "}"
	output += "\n"
	output += relations
	return output, nil
}

func main() {
	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	cts, err := constNames(content)
	if err != nil {
		log.Fatal(err)

	}
	spec := string(content)

	output := "package fact\n"
	output += "\n"
	output += "var sourceDefinition = []byte(`"
	output += spec
	output += "`)\n"

	output += cts

	fmt.Println(output)
}

func snakeCaseToCamelCase(inputUnderScoreStr string) (camelCase string) {
	//snake_case to camelCase

	isToUpper := false

	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return

}
