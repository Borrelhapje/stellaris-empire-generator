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
	empire = chooseOrigin(empire)
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

func chooseOrigin(empire Empire) Empire {
	result := []Origin{}
	for _, origin := range allOrigins {
		if origin.isAllowed(empire) {
			result = append(result, origin)
		}
	}
	empire.origin = result[r.Intn(len(result))]
	return empire
}

type Empire struct {
	authority   string
	civics      []Civic
	ethics      []Ethic
	origin      Origin
	homeplanet  string
	mainSpecies Species
	subSpecies  Species
}

type Predicate func(empire Empire) bool

type Civic struct {
	name      string
	genocidal bool
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

func always(empire Empire) bool {
	return true
}

var allCivics = []Civic{
	{name: "Constructobot", isAllowed: auth("Machine Intelligence")},
	{name: "Delegated Functions", isAllowed: auth("Machine Intelligence")},
	{name: "Determined Exterminator", isAllowed: and(auth("Machine Intelligence"), excludeCivic("Driven Assimilator", "Rogue Servitor")), genocidal: true},
	{name: "Driven Assimilator", isAllowed: and(auth("Machine Intelligence"), excludeCivic("Determined Exterminator", "Rogue Servitor")), genocidal: true},
	{name: "Factory Overclocking", isAllowed: auth("Machine Intelligence")},
	{name: "Introspective", isAllowed: auth("Machine Intelligence")},
	{name: "Maintenance Protocols", isAllowed: auth("Machine Intelligence")},
	{name: "Memorialists", isAllowed: auth("Machine Intelligence")},
	{name: "OTA Updates", isAllowed: auth("Machine Intelligence")},
	{name: "Organic Reprocessing", isAllowed: auth("Machine Intelligence")},
	{name: "Rapid Replicator", isAllowed: auth("Machine Intelligence")},
	{name: "Rockbreakers", isAllowed: auth("Machine Intelligence")},
	{name: "Rogue Servitor", isAllowed: and(auth("Machine Intelligence"), excludeCivic("Determined Exterminator", "Driven Assimilator"))},
	{name: "Static Research Analysis", isAllowed: auth("Machine Intelligence")},
	{name: "Unitary Cohesion", isAllowed: auth("Machine Intelligence")},
	{name: "Warbots", isAllowed: auth("Machine Intelligence")},
	{name: "Zero-Waste Protocols", isAllowed: auth("Machine Intelligence")},
}

var allOrigins = []Origin{
	{name: "Prosperous Unification", isAllowed: always},
	{name: "Mechanist", isAllowed: and(includeEthic("Materialist", "Fanatic Materialist"), excludeCivic("Permanent Employment"))},
	{name: "Syncretic Evolution", isAllowed: and(excludeEthic("Gestalt Consciousness"), excludeCivic("Fanatic Purifiers"))},
	{name: "Life-Seeded", isAllowed: notAuth("Machine Intelligence")},
	{name: "Post-Apocalyptic", isAllowed: and(notAuth("Machine Intelligence"), excludeCivic("Agrarian Idyll", "Anglers"))},
	{name: "Remnants", isAllowed: excludeCivic("Agrarian Idyll")},
	{name: "Shattered Ring", isAllowed: excludeCivic("Agrarian Idyll", "Anglers")},
	{name: "Void Dwellers", isAllowed: and(excludeEthic("Machine Intelligence"), excludeCivic("Idyllic Bloom", "Agrarian Idyll", "Anglers"))},
	{name: "Scion", isAllowed: and(excludeEthic("Gestalt Consciousness", "Fanatic Xenophobe"), excludeCivic("Pompous Purists"))},
	{name: "Galactic Doorstep", isAllowed: always},
	{name: "Tree of Life", isAllowed: and(auth("Hive Mind"), excludeCivic("Devouring Swarm", "Terravore"))},
	{name: "On the Shoulders of Giants", isAllowed: excludeEthic("Gestalt Consciousness")},
	{name: "Calamitous Birth", isAllowed: excludeCivic("Catalytic Processing", "Organic Reprocessing", "Catalytic Recyclerss")},
	{name: "Resource Consolidation", isAllowed: and(auth("Machine Intelligence"), excludeCivic("Rogue Servitor", "Organic Reprocessing"))},
	{name: "Common Ground", isAllowed: and(excludeEthic("Gestalt Consciousness", "Xenophobe", "Fanatic Xenophobe"), excludeCivic("Barbaric Despoilers", "Fanatic Purifiers", "Inward Perfection"))},
	{name: "Hegemon", isAllowed: and(excludeEthic("Gestalt Consciousness", "Xenophobe", "Fanatic Xenophobe", "Egalitarian", "Fanatic Egalitarian"), excludeCivic("Fanatic Purifiers", "Inward Perfection"))},
	{name: "Doomsday", isAllowed: always},
	{name: "Lost Colony", isAllowed: excludeEthic("Gestalt Consciousness")},
	{name: "Necrophage", isAllowed: and(excludeEthic("Xenophile", "Fanatic Xenophile", "Fanatic Egalitarian"), notAuth("Machine Intelligence"), excludeCivic("Death Cult", "Corporate Death Cult", "Empath", "Permanent Employment"))},
	{name: "Clone Army", isAllowed: and(excludeEthic("Gestalt Consciousness"), excludeCivic("Permanent Employment"))},
	{name: "Here Be Dragons", isAllowed: excludeCivic("Fanatic Purifiers", "Devouring Swarm", "Terravore", "Determined Exterminator")},
	{name: "Ocean Paradise", isAllowed: notAuth("Machine Intelligence")},
	{name: "Progenitor Hive", isAllowed: auth("Hive Mind")},
	{name: "Subterrenean", isAllowed: notAuth("Machine Intelligence")},
	{name: "Slingshot to the Stars", isAllowed: always},
	{name: "Teachers of the Shroud", isAllowed: and(includeEthic("Spiritualist", "Fanatic Spiritualist"), excludeCivic("Fanatic Purifiers"))},
	{name: "Imperial Fiefdom", isAllowed: excludeCivic("Inward Perfection", "Fanatic Purifiers", "Devouring Swarm", "Terravore", "Driven Assimilator", "Determined Exterminator")},
}

func auth(s ...string) Predicate {
	return func(empire Empire) bool {
		for _, auth := range s {
			if empire.authority == auth {
				return true
			}
		}
		return false
	}
}

func notAuth(s ...string) Predicate {
	return func(empire Empire) bool {
		for _, auth := range s {
			if empire.authority == auth {
				return false
			}
		}
		return true
	}
}

func excludeCivic(s ...string) Predicate {
	return func(empire Empire) bool {
		for _, civic := range empire.civics {
			for _, excluded := range s {
				if excluded == civic.name {
					return false
				}
			}
		}
		return true
	}
}

func excludeEthic(s ...string) Predicate {
	return func(empire Empire) bool {
		for _, ethic := range empire.ethics {
			for _, excluded := range s {
				if excluded == ethic.name {
					return false
				}
			}
		}
		return true
	}
}

func includeEthic(s ...string) Predicate {
	return func(empire Empire) bool {
		for _, ethic := range empire.ethics {
			for _, included := range s {
				if included == ethic.name {
					return true
				}
			}
		}
		return false
	}
}

func and(s ...Predicate) Predicate {
	return func(empire Empire) bool {
		for _, pred := range s {
			if !pred(empire) {
				return false
			}
		}
		return true
	}
}

func or(s ...Predicate) Predicate {
	return func(empire Empire) bool {
		for _, pred := range s {
			if pred(empire) {
				return true
			}
		}
		return false
	}
}
