package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"

	"github.com/goccy/go-yaml"
)

// Root represents the top-level structure of the YAML file
type Root struct {
	Scenarios map[string]map[string][]interface{} `yaml:"scenarios"`
}

func flattenTestSuites(rawList []interface{}) map[string]bool {
	suiteSet := make(map[string]bool)

	for _, item := range rawList {
		switch v := item.(type) {
		case string:
			// Case: - gnome-agama
			suiteSet[v] = true
		case map[string]interface{}:
			// Case: - gnome: { machine: ... }
			for key := range v {
				suiteSet[key] = true
				break
			}
		}
	}
	return suiteSet
}

func findProductSuites(root *Root, productName string) ([]interface{}, bool) {
	for _, products := range root.Scenarios {
		if suites, ok := products[productName]; ok {
			return suites, true
		}
	}
	return nil, false
}

func loadYaml(filename string) (*Root, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var root Root
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, err
	}
	return &root, nil
}

func main() {
	// 1. Parse Arguments
	allsuites := flag.Bool("all", false, "List also testsuites present in both schedules")
	sourceFile := flag.String("source", "", "Path to the source YAML file")
	targetFile := flag.String("target", "", "Path to the target YAML file")
	sourceProduct := flag.String("source-product", "", "Product name in the source file")
	targetProduct := flag.String("target-product", "", "Product name in the target file (optional, defaults to source-product)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nCompares testsuites two products in openQA YAML job group files.\n")
		fmt.Fprintf(os.Stderr, "It lists suites missing in the target, and source and ptionally those present in both.\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s --source fileA.yaml --source-product opensuse-Tumbleweed-agama-installer-x86_64 --target fileB.yaml\n", os.Args[0])
	}

	flag.Parse()

	if *sourceFile == "" || *targetFile == "" || *sourceProduct == "" {
		flag.Usage()
		os.Exit(1)
	}

	// If target-product isn't specified, assume it's the same as source-product
	if *targetProduct == "" {
		*targetProduct = *sourceProduct
	}

	// 2. Load Files
	rootA, err := loadYaml(*sourceFile)
	if err != nil {
		log.Fatalf("Error loading source: %v", err)
	}

	rootB, err := loadYaml(*targetFile)
	if err != nil {
		log.Fatalf("Error loading target: %v", err)
	}

	// 3. Find the specific product lists in the nested maps
	suitesA, foundA := findProductSuites(rootA, *sourceProduct)
	if !foundA {
		log.Fatalf("Product '%s' not found in source file '%s'", *sourceProduct, *sourceFile)
	}

	suitesB, foundB := findProductSuites(rootB, *targetProduct)
	if !foundB {
		log.Fatalf("Product '%s' not found in target file '%s'", *targetProduct, *targetFile)
	}

	// 4. Flatten and Compare
	setA := flattenTestSuites(suitesA)
	setB := flattenTestSuites(suitesB)

	fmt.Printf("Diff: What is in '%s' (%s) but missing in '%s' (%s)?\n", *sourceFile, *sourceProduct, *targetFile, *targetProduct)
	fmt.Println("--------------------------------------------------------------------------------")

	foundInA, suitesInA := diff(setA, *sourceFile, setB, *targetFile, "-")
	foundInB, suitesInB := diff(setB, *sourceFile, setA, *targetFile, "+")

	if !foundInA && !foundInB {
		fmt.Println("No missing testsuites found.")
	}
	if *allsuites {
		combined := append(suitesInA, suitesInB...)
		sortedSuites := make([]string, 0, len(combined))
		for _, k := range combined {
			sortedSuites = append(sortedSuites, k)
		}
		sort.Strings(sortedSuites)
		sortedSuites = slices.Compact(sortedSuites)

		fmt.Printf("--- Present in both")
		for _, suite := range sortedSuites {
			fmt.Printf("= %s\n", suite)
		}
	}
}

func diff(source map[string]bool, sourceFile string, target map[string]bool, targetFile string, sign string) (bool, []string) {
	var suitePresentInTarget []string
	suiteListKeys := make([]string, 0, len(source))
	for k := range source {
		suiteListKeys = append(suiteListKeys, k)
	}

	sort.Strings(suiteListKeys)

	fmt.Printf("--- In '%s' but MISSING in '%s' ---\n", sourceFile, targetFile)
	diffFound := false
	for _, suite := range suiteListKeys {
		if !target[suite] {
			fmt.Printf("%s %s\n", sign, suite)
			diffFound = true
		} else {
			suitePresentInTarget = append(suitePresentInTarget, suite)
		}
	}

	return diffFound, suitePresentInTarget
}
