package classifier

import (
	"fmt"
	"regexp"
	"strings"
)

// Classify will take a string and attempt to classify it into which products
// are discussed
func Classify(rawDesc string) []string {
	desc := stripString(rawDesc)

	var keywords = make(map[string][]string)
	keywords["Consul"] = []string{"consul"}
	keywords["Consul Connect"] = []string{"consul connect"}
	keywords["HashiConf"] = []string{"hashiconf"}
	keywords["Nomad"] = []string{"nomad"}
	keywords["Packer"] = []string{"packer"}
	keywords["Terraform"] = []string{"terraform", "statefile", "statefiles"}
	keywords["Vagrant"] = []string{"vagrant"}
	keywords["Vault"] = []string{"vault"}

	var counts = make(map[string]int)
	for k, v := range keywords {
		counts[k] = 0
		for _, s := range v {
			c := countOccurrences(desc, s)
			counts[k] = counts[k] + c
		}
	}

	var found = []string{}
	for k, c := range counts {
		if c != 0 {
			found = append(found, k)
		}
	}
	return found
}

func countOccurrences(desc, searchStr string) int {
	reg := regexp.MustCompile(fmt.Sprintf(" %s ", searchStr))
	matches := reg.FindAllStringIndex(desc, -1)
	return len(matches)
}

func stripString(s string) string {
	// Make a Regex to say we only want letters and numbers and whitespaces
	reg := regexp.MustCompile("[^a-zA-Z]+")
	return reg.ReplaceAllString(strings.ToLower(s), " ")
}
