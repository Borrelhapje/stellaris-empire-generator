package main

import (
	"fmt"
	"math/rand"
	"time"
)

var seed = time.Now().UnixNano()
var r = rand.New(rand.NewSource(seed))

func main() {
	fmt.Println("Seed is " + fmt.Sprint(seed))
	empire := Empire{}
	empire = chooseAuthority(empire)
	empire = chooseCivic(empire)
	empire = chooseCivic(empire)
	for i := 0; i < 3; i++ {
		empire = chooseEthic(empire)
	}
	fmt.Println(empire)
}

func chooseAuthority(empire Empire) Empire {
	var authorities = []string{"Democratic", "Oligarchy", "Dictatorial", "Imperial", "Hive Mind", "Machine Intelligence", "Corporate"}
	empire.authority = authorities[r.Intn(len(authorities))]
	return empire
}

func chooseCivic(empire Empire) Empire {
	civicList := getCivicList(empire)
	empire.civics = append(empire.civics, civicList[r.Intn(len(civicList))])
	return empire
}

func getCivicList(empire Empire) []Civic {
	allCivics := []Civic{}
	result := []Civic{}
	for _, civic := range allCivics {
		if civic.isAllowed(empire) {
			result = append(result, civic)
		}
	}
	return result
}

func chooseEthic(empire Empire) Empire {
	ethicList := getEthicList(empire)
	empire.ethics = append(empire.ethics, ethicList[r.Intn(len(ethicList))])
	return empire
}

func getEthicList(empire Empire) []Ethic {
	allEthics := []Ethic{}
	result := []Ethic{}
	for _, ethic := range allEthics {
		if ethic.isAllowed(empire) {
			result = append(result, ethic)
		}
	}
	return result
}

type Empire struct {
	authority   string
	civics      []Civic
	ethics      []Ethic
	origin      string
	homeplanet  string
	mainSpecies Species
	subSpecies  Species
}

type Predicate func(empire Empire) bool

type Civic struct {
	name      string
	isAllowed Predicate // should only check for other civics and authority
}

type Ethic struct {
	name      string
	isAllowed Predicate // checks if valid for civics, authority and other ethics
}

type Origin struct {
	name      string
	isAllowed Predicate // checks if valid for civics, authority and ethics
}

type Species struct {
	popType string
	traits  []string
}

var allCivics = []Civic{}
