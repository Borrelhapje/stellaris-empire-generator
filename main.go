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
	empire = chooseEthic(empire)
	empire = chooseAuthority(empire)
	empire = chooseCivic(empire)
	empire = chooseCivic(empire)

	empire = chooseOrigin(empire)
	fmt.Println(empire)
}

func chooseAuthority(empire Empire) Empire {
	result := []Authority{}
	for _, auth := range allAuthorities {
		if auth.isAllowed(empire) {
			result = append(result, auth)
		}
	}
	empire.authority = result[r.Intn(len(result))].name
	return empire
}

func chooseCivic(empire Empire) Empire {
	civicList := getCivicList(empire)
	empire.civics = append(empire.civics, civicList[r.Intn(len(civicList))])
	return empire
}

func getCivicList(empire Empire) []Civic {
	result := []Civic{}
	for _, civic := range allCivics {
		if civic.isAllowed(empire) {
			result = append(result, civic)
		}
	}
	return result
}

func chooseEthic(empire Empire) Empire {
	firstFanatic := r.Intn(2) == 1
	ethicList := getEthicList(empire)
	firstDraw := ethicList[r.Intn(len(ethicList))]
	if firstDraw.name == "Gestalt Consciousness" {
		empire.ethics = []Ethic{firstDraw}
		return empire
	}
	if firstFanatic {
		empire.ethics = append(empire.ethics, Ethic{name: "Fanatic " + firstDraw.name, isAllowed: firstDraw.isAllowed})
		ethicList := getEthicList(empire)
		nextDraw := ethicList[r.Intn(len(ethicList))]
		empire.ethics = append(empire.ethics, nextDraw)
	} else {
		empire.ethics = append(empire.ethics, firstDraw)
		for i := 0; i < 2; i++ {
			ethicList := getEthicList(empire)
			nextDraw := ethicList[r.Intn(len(ethicList))]
			empire.ethics = append(empire.ethics, nextDraw)
		}
	}
	return empire
}

func getEthicList(empire Empire) []Ethic {
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

type Authority struct {
	name      string
	isAllowed Predicate
}

type Species struct {
	popType string
	traits  []string
}

func always(empire Empire) bool {
	return true
}

func normalAuth() Predicate {
	return auth("Democratic", "Oligarchy", "Dictatorial", "Imperial")
}

func onlyGestalt(empire Empire) bool {
	return len(empire.ethics) == 0
}

var allEthics = []Ethic{
	{name: "Authoritarian", isAllowed: excludeEthic("Egalitarian", "Gestalt Consciousness")},
	{name: "Spiritualist", isAllowed: excludeEthic("Materialist", "Gestalt Consciousness")},
	{name: "Militarist", isAllowed: excludeEthic("Pacifist", "Gestalt Consciousness")},
	{name: "Xenophobe", isAllowed: excludeEthic("Xenophile", "Gestalt Consciousness")},
	{name: "Egalitarian", isAllowed: excludeEthic("Authoritarian", "Gestalt Consciousness")},
	{name: "Materialist", isAllowed: excludeEthic("Materialist", "Gestalt Consciousness")},
	{name: "Pacifist", isAllowed: excludeEthic("Militarist", "Gestalt Consciousness")},
	{name: "Xenophile", isAllowed: excludeEthic("Xenophobe", "Gestalt Consciousness")},
	{name: "Gestalt Consciousness", isAllowed: onlyGestalt},
	{name: "Gestalt Consciousness", isAllowed: onlyGestalt},
	{name: "Gestalt Consciousness", isAllowed: onlyGestalt},
}

var allAuthorities = []Authority{
	{name: "Democratic", isAllowed: excludeEthic("Authoritarian", "Fanatic Authoritarian", "Gestalt Consciousness")},
	{name: "Oligarchy", isAllowed: excludeEthic("Fanatic Authoritarian", "Fanatic Egalitarian", "Gestalt Consciousness")},
	{name: "Dictatorial", isAllowed: excludeEthic("Egalitarian", "Fanatic Egalitarian", "Gestalt Consciousness")},
	{name: "Imperial", isAllowed: excludeEthic("Egalitarian", "Fanatic Egalitarian", "Gestalt Consciousness")},
	{name: "Corporate", isAllowed: excludeEthic("Fanatic Egalitarian", "Fanatic Authoritarian", "Gestalt Consciousness")},
	{name: "Hive Mind", isAllowed: includeEthic("Gestalt Consciousness")},
	{name: "Machine Intelligence", isAllowed: includeEthic("Gestalt Consciousness")},
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
	{name: "Ascetic", isAllowed: auth("Hive Mind")},
	{name: "Devouring Swarm", isAllowed: and(auth("Hive Mind"), excludeCivic("Terravore")), genocidal: true},
	{name: "Terravore", isAllowed: and(auth("Hive Mind"), excludeCivic("Devouring Swarm")), genocidal: true},
	{name: "Divided Attention", isAllowed: auth("Hive Mind")},
	{name: "Empath", isAllowed: auth("Hive Mind")},
	{name: "Idyllic Bloom", isAllowed: auth("Hive Mind")},
	{name: "Memorialist", isAllowed: auth("Hive Mind")},
	{name: "Natural NeuralNetwork", isAllowed: auth("Hive Mind")},
	{name: "One Mind", isAllowed: auth("Hive Mind")},
	{name: "Organic Reprocessing", isAllowed: auth("Hive Mind")},
	{name: "Pooled Knowledge", isAllowed: auth("Hive Mind")},
	{name: "Strength of Legions", isAllowed: auth("Hive Mind")},
	{name: "Subspacce Ephase", isAllowed: auth("Hive Mind")},
	{name: "Subusmed Will", isAllowed: auth("Hive Mind")},
	{name: "Brand Loyalty", isAllowed: auth("Corporate")},
	{name: "Catalytic Recyclers", isAllowed: auth("Corporate")},
	{name: "Corporate Hedonism", isAllowed: and(auth("Corporate"), excludeCivic("Indentured Assets"))},
	{name: "Criminal Heritage", isAllowed: auth("Corporate")},
	{name: "Franchising", isAllowed: auth("Corporate")},
	{name: "Free Traders", isAllowed: auth("Corporate")},
	{name: "Mastercraft Inc.", isAllowed: auth("Corporate")},
	{name: "Media Conglomerate", isAllowed: auth("Corporate")},
	{name: "Permanent Employment", isAllowed: and(auth("Corporate"))},
	{name: "Private Prospectors", isAllowed: auth("Corporate")},
	{name: "Public Relations Specialists", isAllowed: auth("Corporate")},
	{name: "Ruthless Competition", isAllowed: auth("Corporate")},
	{name: "Trading Posts", isAllowed: auth("Corporate")},
	{name: "Corporate Death Cult", isAllowed: auth("Corporate")},
	{name: "Gospel of the Masses", isAllowed: auth("Corporate")},
	{name: "Indentured Assets", isAllowed: and(auth("Corporate"), excludeCivic("Corporate Hedonism"))},
	{name: "Naval Contractors", isAllowed: auth("Corporate")},
	{name: "Private Military Companies", isAllowed: auth("Corporate")},
	{name: "Anglers", isAllowed: and(normalAuth(), excludeCivic("Agrarian Idyll"))},
	{name: "Byzantine Bureaucracy", isAllowed: normalAuth()},
	{name: "Corvee System", isAllowed: and(normalAuth(), excludeCivic("Free Haven"))},
	{name: "Cutthroat Politics", isAllowed: normalAuth()},
	{name: "Diplomatic Corps", isAllowed: and(normalAuth(), excludeCivic("Fanatic Purifiers", "Inward Perfection"))},
	{name: "Efficient Bureaucracy", isAllowed: normalAuth()},
	{name: "Environmentalist", isAllowed: normalAuth()},
	{name: "Functional Architecture", isAllowed: normalAuth()},
	{name: "Masterful Crafters", isAllowed: normalAuth()},
	{name: "Memorialists", isAllowed: and(normalAuth(), excludeCivic("Fanatic Purifiers"))},
	{name: "Merchant Guilds", isAllowed: and(normalAuth(), excludeCivic("Exalted Priesthood", "Aristocratic Elite", "Technocracy"))},
	{name: "Mining Guilds", isAllowed: normalAuth()},
	{name: "Philosopher King", isAllowed: auth("Dictatorial", "Imperial")},
	{name: "Pleasure Seekers", isAllowed: and(normalAuth(), excludeCivic("Warrior Culture", "Shared Burdens", "Slaver Guilds"))},
	{name: "Police State", isAllowed: normalAuth()},
	{name: "Shadow Council", isAllowed: auth("Democratic", "Oligarchy", "Dictatorial")},
	{name: "Aristocratic Elite", isAllowed: and(auth("Oligarchy", "Dictatorial"), excludeCivic("Exalted Priesthood", "Merchant Guilds", "Technocracy"))},
	{name: "Beacon of Libery", isAllowed: auth("Democratic")},
	{name: "Citizen Service", isAllowed: and(auth("Democratic", "Oligarchy"), excludeCivic("Reanimators"))},
	{name: "Death Cult", isAllowed: and(normalAuth(), excludeCivic("Fanatic Purifiers", "Inward Perfection"))},
	{name: "Distinquished Admiralty", isAllowed: normalAuth()},
	{name: "Exalted Priesthood", isAllowed: and(auth("Oligarchy", "Dictatorial"), excludeCivic("Aristocratic Elite", "Merchant Guilds", "Technocracy"))},
	{name: "Feudal Society", isAllowed: auth("Imperial")},
	{name: "Free Haven", isAllowed: and(normalAuth(), excludeCivic("Corvee System"))},
	{name: "Idyllic Bloom", isAllowed: auth("Imperial")},
	{name: "Imperial Cult", isAllowed: normalAuth()},
	{name: "Inward Perfection", isAllowed: and(normalAuth(), excludeCivic("Pompous Purists"))},
	{name: "Meritocracy", isAllowed: auth("Democratic", "Oligarchy")},
	{name: "Nationalistic Zeal", isAllowed: normalAuth()},
	{name: "Parliamentary System", isAllowed: auth("Democratic")},
	{name: "Pompous Purists", isAllowed: and(normalAuth(), excludeCivic("Fanatic Purifiers", "Inward Perfection"))},
	{name: "Shared Burdens", isAllowed: and(normalAuth(), excludeCivic("Technocracy", "Pleasure Seekers"))},
	{name: "Slaver Guilds", isAllowed: and(normalAuth(), excludeCivic("Pleasure Seekers"))},
	{name: "Technocracy", isAllowed: and(normalAuth(), excludeCivic("Exalted Priesthood", "Merchant Guilds", "Aristocratic Elite", "Shared Burdens"))},
	{name: "Warrior Culture", isAllowed: and(normalAuth(), excludeCivic("Pleasure Seekers"))},
	//here we have the civics with a slight edit to their requirements, because their ethics are very strict
	{name: "Idealistic Foundation", isAllowed: auth("Democratic", "Oligarchy")},
	{name: "Reanimators", isAllowed: and(normalAuth(), excludeCivic("Citizen Service"))},
	{name: "Agrarian Idyll", isAllowed: and(normalAuth(), excludeCivic("Anglers"))},
	{name: "Barbaric Despoilers", isAllowed: and(normalAuth(), excludeCivic("Fanatic Purifiers"))},
	{name: "Fanatic Purifiers", isAllowed: and(normalAuth(), excludeCivic("Barbaric Despoilers", "Pompous Purists")), genocidal: true},
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
	{name: "Calamitous Birth", isAllowed: excludeCivic("Catalytic Processing", "Organic Reprocessing", "Catalytic Recyclers")},
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
