package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

// var urlMap = make(map[string]int)
// var domainSet = make(map[string]bool)
// var ipMap = make(map[string]int)

type Maps struct {
	urlMap    map[string]int
	domainSet map[string]int
	ipMap     map[string]int
}

func parse_args() []string {
	totalArgs := len(os.Args[1:])

	if totalArgs < 1 {
		fmt.Println("Not enough arguments. so exiting")
		os.Exit(-1)
	}

	arguments := make([]string, totalArgs)
	copy(arguments[:], os.Args[1:])

	return arguments
}

func fillMap(m map[string]int, key string) {
	if value, ok := m[key]; ok {
		value = value + 1
		m[key] = value
	} else {
		m[key] = 1
	}

}

func fillDomainSet(maps Maps) {

	re := regexp.MustCompile(`https://([^/]+)`)

	for key, _ := range maps.urlMap {
		match := re.FindStringSubmatch(key)

		if len(match) > 1 {
			// maps.domainSet[match[1]] = true
			fillMap(maps.domainSet, match[1])
		}
	}

}

func sortMap(m map[string]int) []string {

	slicey := make([]string, 0, len(m))

	for key, _ := range m {
		slicey = append(slicey, key)
	}

	sort.SliceStable(slicey, func(i, j int) bool {
		return m[slicey[i]] > m[slicey[j]]
	})

	return slicey
}

func traverseFile(filename string, maps Maps) Maps {

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	buffer := make([]byte, 0, 64*1024)
	scanner.Buffer(buffer, 1024*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		details := strings.Fields(line)

		// fill and sort url map
		fillMap(maps.urlMap, details[3])
		// fill and sort ip map
		fillMap(maps.ipMap, details[2])

	}
	// domain set would get filled in sorted order
	fillDomainSet(maps)

	// sortMap(maps.urlMap)
	// sortMap(maps.ipMap)
	// sortMap(maps.domainSet)

	return maps
}

func displayResult(maps Maps, arguments []string, elapsed time.Duration) {

	for _, value := range arguments {
		fmt.Println("Reading ", value, "...")
	}

	fmt.Println("* Unique URLS: ", len(maps.urlMap))
	fmt.Println("* Unique Domains: ", len(maps.domainSet))
	fmt.Println("* Top 10 Websites: ")

	sortedDomains := sortMap(maps.domainSet)

	i := 0
	for _, element := range sortedDomains {

		fmt.Println("\t - ", element)

		if i >= 10 {
			break
		}

		i++
	}
	fmt.Println("* Top 5 crawlers: ")

	sortedIP := sortMap(maps.ipMap)
	i = 0
	for _, element := range sortedIP {

		fmt.Println("\t - ", element)

		if i >= 5 {
			break
		}

		i++
	}
	fmt.Println("Completed in ", elapsed)
}

func main() {

	// Parse arguments
	arguments := parse_args()
	maps := Maps{
		urlMap:    make(map[string]int),
		domainSet: make(map[string]int),
		ipMap:     make(map[string]int),
	}
	// Read from each file from an array
	start := time.Now()
	for _, value := range arguments {
		// fmt.Println(value)
		maps = traverseFile(value, maps)
	}

	elapsed := time.Since(start)
	displayResult(maps, arguments, elapsed)

}
