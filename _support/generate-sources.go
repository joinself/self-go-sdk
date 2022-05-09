package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
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

	output := ""
	sources := map[string]struct{}{}
	facts := map[string]struct{}{}
	definitions := []string{}
	relations := "var spec = map[string][]string{\n"

	keys := make([]string, 0, len(spec.Sources))
	for k := range spec.Sources {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		specFacts := spec.Sources[key]
		source := snakeCaseToCamelCase("Source_" + key)
		relations += fmt.Sprintf("\t%s: []string{\n", source)
		if _, ok := sources[key]; !ok {
			definitions = append(definitions, fmt.Sprintf("const %s = \"%s\"", source, key))
			sources[key] = struct{}{}
		}

		for _, f := range specFacts {
			fact := snakeCaseToCamelCase("Fact_" + f)
			relations += fmt.Sprintf("\t\t%s,\n", fact)
			if _, ok := facts[f]; !ok {
				definitions = append(definitions, fmt.Sprintf("const %s = \"%s\"", fact, f))
				facts[f] = struct{}{}
			}
		}
		relations += fmt.Sprintf("\t},\n")
	}
	sort.Strings(definitions)
	relations += "}"

	output += strings.Join(definitions, "\n") + "\n\n"
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

	output := "package fact\n"
	output += "\n"
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
