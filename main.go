package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var seed = time.Now().UnixNano()
var r = rand.New(rand.NewSource(seed))

func main() {
	fmt.Println("Seed is " + fmt.Sprint(seed))
	app.Route("/", &data{})
	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:        "Stellaris",
		Description: "A Stellaris Empire Generator",
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func (d *data) Render() app.UI {
	return app.Div().Body(
		app.Button().Text("Generate").OnClick(d.generateEmpire),
		app.Range(d.Empires).Slice(func(i int) app.UI {
			return app.Div().Body(
				app.Label().Text("Authority:").For("authority"),
				app.Span().ID("authority").Text(d.Empires[i].authority),
				app.Br(),
				app.Label().Text("Ethics:").For("ethics"),
				app.Span().ID("ethics").Text(d.Empires[i].ethics),
				app.Br(),
				app.Label().Text("Civics:").For("civics"),
				app.Span().ID("civics").Text(d.Empires[i].civics),
				app.Br(),
				app.Label().Text("Origin:").For("origin"),
				app.Span().ID("origin").Text(d.Empires[i].origin.name),
				app.Br(),
				app.Label().Text("Planet Class:").For("planet"),
				app.Span().ID("planet").Text(d.Empires[i].homeplanet),
			)
		}),
	)
}

type data struct {
	app.Compo
	Empires []Empire
}

func (d *data) generateEmpire(ctx app.Context, e app.Event) {
	d.Empires = []Empire{}
	for i := 0; i < 3; i++ {
		empire := Empire{}
		empire = chooseEthic(empire)
		empire = chooseAuthority(empire)
		empire = chooseCivic(empire)
		empire = chooseCivic(empire)
		empire = chooseOrigin(empire)
		empire = chooseHomeplanet(empire)
		d.Empires = append(d.Empires, empire)
	}
}

func (e Empire) String() string {
	res := e.authority + "\nEthics: "
	for _, ethic := range e.ethics {
		res += ethic.name + " "
	}
	res += "\nCivics: "
	for _, civic := range e.civics {
		res += civic.name + " "
	}
	res += "\nOrigin: " + e.origin.name
	res += "\nPlanet: " + e.homeplanet
	return res
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
outer:
	for _, civic := range allCivics {
		if civic.isAllowed(empire) {
			for _, existing := range empire.civics {
				if existing.name == civic.name {
					continue outer
				}
			}
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
outer:
	for _, ethic := range allEthics {
		if ethic.isAllowed(empire) {
			for _, existing := range empire.ethics {
				if existing.name == ethic.name || existing.name == "Fanatic "+ethic.name {
					continue outer
				}
			}
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

func chooseHomeplanet(empire Empire) Empire {
	planets := []string{"Desert", "Arid", "Savanna", "Ocean", "Continental", "Tropical", "Arctic", "Alpine", "Tundra"}
	empire.homeplanet = planets[r.Intn(len(planets))]
	return empire
}

func generateSpecies(empire Empire) Species {
	species := Species{}
	if empire.authority == "Machine Intelligence" {
		//generate machine species
		species.popType = "Machine"
		species.initialTraitPoints = 1
	} else if empire.origin.name == "Calamitous Birth" {
		//generate lithoid species
		species.popType = "Lithoid"
		species.initialTraitPoints = 2
	} else if empire.origin.name == "Ocean Paradise" {
		//aquatic species, forced aqautic trait
		species.popType = "Aquatic"
		species.traits = append(species.traits, Trait{name: "Aquatic", cost: 1})
		species.initialTraitPoints = 2
	} else {
		//standard species
		popTypes := []string{"Aquatic", "Mammalian", "Reptilian", "Avian", "Arthropoid", "Molluscoid", "Fungoid", "Plantoid", "Lithoid", "Necroid"}
		species.popType = popTypes[r.Intn(len(popTypes))]
		species.initialTraitPoints = 2
	}
	return fillSpecies(species)
}

func fillSpecies(s Species) Species {
	for len(s.traits) < 5 {
		filteredTraits := []Trait{}
		for _, trait := range allTraits {
			if trait.isAllowed(s) {
				filteredTraits = append(filteredTraits, trait)
			}
		}
		s.traits = append(s.traits, filteredTraits[r.Intn(len(filteredTraits))])
	}
	if s.availablePoints() > 0 {
		// try to replace a trait with one costing one more
	}
	return s
}

func (s Species) availablePoints() int {
	res := s.initialTraitPoints
	for _, trait := range s.traits {
		res += trait.cost
	}
	return res
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

func (e Civic) String() string {
	return e.name
}

func (e Ethic) String() string {
	return e.name
}

func (e Origin) String() string {
	return e.name
}

func (e Authority) String() string {
	return e.name
}

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
	popType            string
	initialTraitPoints int
	traits             []Trait
}

type speciesPredicate func(s Species) bool

type Trait struct {
	cost      int
	name      string
	isAllowed speciesPredicate
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
	{name: "Authoritarian", isAllowed: excludeEthic("Egalitarian", "Fanatic Egalitarian", "Gestalt Consciousness")},
	{name: "Spiritualist", isAllowed: excludeEthic("Materialist", "Fanatic Materialist", "Gestalt Consciousness")},
	{name: "Militarist", isAllowed: excludeEthic("Pacifist", "Fanatic Pacifist", "Gestalt Consciousness")},
	{name: "Xenophobe", isAllowed: excludeEthic("Xenophile", "Fanatic Xenophile", "Gestalt Consciousness")},
	{name: "Egalitarian", isAllowed: excludeEthic("Authoritarian", "Fanatic Authoritarian", "Gestalt Consciousness")},
	{name: "Materialist", isAllowed: excludeEthic("Spiritualist", "Fanatic Spiritualist", "Gestalt Consciousness")},
	{name: "Pacifist", isAllowed: excludeEthic("Militarist", "Fanatic Militarist", "Gestalt Consciousness")},
	{name: "Xenophile", isAllowed: excludeEthic("Xenophobe", "Fanatic Xenophobe", "Gestalt Consciousness")},
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
	{name: "Devouring Swarm", isAllowed: and(auth("Hive Mind"), excludeCivic("Terravore", "Empath")), genocidal: true},
	{name: "Terravore", isAllowed: and(auth("Hive Mind"), excludeCivic("Devouring Swarm", "Empath")), genocidal: true},
	{name: "Divided Attention", isAllowed: auth("Hive Mind")},
	{name: "Empath", isAllowed: and(auth("Hive Mind"), excludeCivic("Terravore", "Devouring Swarm"))},
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
	{name: "Permanent Employment", isAllowed: and(auth("Corporate"), excludeEthic("Egalitarian", "Fanatic Egalitarian"))},
	{name: "Private Prospectors", isAllowed: auth("Corporate")},
	{name: "Public Relations Specialists", isAllowed: auth("Corporate")},
	{name: "Ruthless Competition", isAllowed: auth("Corporate")},
	{name: "Trading Posts", isAllowed: auth("Corporate")},
	{name: "Corporate Death Cult", isAllowed: and(auth("Corporate"), includeEthic("Spiritualist", "Fanatic Spiritualist"))},
	{name: "Gospel of the Masses", isAllowed: and(auth("Corporate"), includeEthic("Spiritualist", "Fanatic Spiritualist"))},
	{name: "Indentured Assets", isAllowed: and(auth("Corporate"), excludeCivic("Corporate Hedonism"), includeEthic("Authoritarian", "Fanatic Authoritarian"))},
	{name: "Naval Contractors", isAllowed: and(auth("Corporate"), includeEthic("Militarist", "Fanatic Militarist"))},
	{name: "Private Military Companies", isAllowed: and(auth("Corporate"), includeEthic("Militarist", "Fanatic Militarist"))},
	{name: "Anglers", isAllowed: and(normalAuth(), excludeCivic("Agrarian Idyll"))},
	{name: "Byzantine Bureaucracy", isAllowed: and(normalAuth(), excludeEthic("Spiritualist", "Fanatic Spiritualist"))},
	{name: "Corvee System", isAllowed: and(normalAuth(), excludeCivic("Free Haven"), excludeEthic("Egalitarian", "Fanatic Egalitarian"))},
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
	{name: "Police State", isAllowed: and(normalAuth(), excludeEthic("Fanatic Egalitarian"))},
	{name: "Shadow Council", isAllowed: auth("Democratic", "Oligarchy", "Dictatorial")},
	{name: "Aristocratic Elite", isAllowed: and(auth("Oligarchy", "Dictatorial"), excludeCivic("Exalted Priesthood", "Merchant Guilds", "Technocracy"), excludeEthic("Egalitarian", "Fanatic Egalitarian"))},
	{name: "Beacon of Libery", isAllowed: and(auth("Democratic"), includeEthic("Egalitarian", "Fanatic Egalitarian"), excludeEthic("Xenophobe", "Fanatic Xenophobe"))},
	{name: "Citizen Service", isAllowed: and(auth("Democratic", "Oligarchy"), excludeCivic("Reanimators"), includeEthic("Militarist", "Fanatic Militarist"), excludeEthic("Fanatic Xenophile"))},
	{name: "Death Cult", isAllowed: and(normalAuth(), excludeCivic("Fanatic Purifiers", "Inward Perfection"), includeEthic("Spiritualist", "Fanatic Spiritualist"))},
	{name: "Distinquished Admiralty", isAllowed: and(normalAuth(), includeEthic("Militarist", "Fanatic Militarist"))},
	{name: "Exalted Priesthood", isAllowed: and(auth("Oligarchy", "Dictatorial"), excludeCivic("Aristocratic Elite", "Merchant Guilds", "Technocracy"), includeEthic("Spiritualist", "Fanatic Spiritualist"))},
	{name: "Feudal Society", isAllowed: auth("Imperial")},
	{name: "Free Haven", isAllowed: and(normalAuth(), excludeCivic("Corvee System"), includeEthic("Xenophile", "Fanatic Xenophile"))},
	{name: "Idyllic Bloom", isAllowed: normalAuth()},
	{name: "Imperial Cult", isAllowed: and(auth("Imperial"), includeEthic("Spiritualist", "Fanatic Spiritualist"), includeEthic("Authoritarian", "Fanatic Authoritarian"))},
	{name: "Inward Perfection", isAllowed: and(normalAuth(), excludeCivic("Pompous Purists"), includeEthic("Pacifist", "Fanatic Pacifist"), includeEthic("Xenophobe", "Fanatic Xenophobe"))},
	{name: "Meritocracy", isAllowed: auth("Democratic", "Oligarchy")},
	{name: "Nationalistic Zeal", isAllowed: and(normalAuth(), includeEthic("Militarist", "Fanatic Militarist"))},
	{name: "Parliamentary System", isAllowed: auth("Democratic")},
	{name: "Pompous Purists", isAllowed: and(normalAuth(), excludeCivic("Fanatic Purifiers", "Inward Perfection"), includeEthic("Xenophobe", "Fanatic Xenophobe"))},
	{name: "Shared Burdens", isAllowed: and(normalAuth(), excludeCivic("Technocracy", "Pleasure Seekers"), includeEthic("Fanatic Egalitarian"), excludeEthic("Xenophobe"))},
	{name: "Slaver Guilds", isAllowed: and(normalAuth(), excludeCivic("Pleasure Seekers"), includeEthic("Authoritarian", "Fanatic Authoritarian"))},
	{name: "Technocracy", isAllowed: and(normalAuth(), excludeCivic("Exalted Priesthood", "Merchant Guilds", "Aristocratic Elite", "Shared Burdens"), includeEthic("Materialist", "Fanatic Materialist"))},
	{name: "Warrior Culture", isAllowed: and(normalAuth(), excludeCivic("Pleasure Seekers"), includeEthic("Militarist", "Fanatic Militarist"))},
	//here we have the civics with a slight edit to their requirements, because their ethics are very strict
	{name: "Idealistic Foundation", isAllowed: and(normalAuth(), includeEthic("Egalitarian", "Fanatic Egalitarian"))},
	{name: "Reanimators", isAllowed: and(normalAuth(), excludeCivic("Citizen Service"), excludeEthic("Pacifist", "Fanatic Pacifist"))},
	{name: "Agrarian Idyll", isAllowed: and(normalAuth(), excludeCivic("Anglers"), includeEthic("Pacifist", "Fanatic Pacifist"))},
	{name: "Barbaric Despoilers", isAllowed: and(normalAuth(), excludeCivic("Fanatic Purifiers"), includeEthic("Militarist", "Fanatic Militarist"), includeEthic("Authoritarian", "Fanatic Authoritarian", "Xenophobe", "Fanatic Xenophobe"), excludeEthic("Xenophile", "Fanatic Xenophile"))},
	{name: "Fanatic Purifiers", isAllowed: and(normalAuth(), excludeCivic("Barbaric Despoilers", "Pompous Purists"), includeEthic("Fanatic Xenophobe"), includeEthic("Militarist", "Spiritualist")), genocidal: true},
}

var allOrigins = []Origin{
	{name: "Prosperous Unification", isAllowed: always},
	{name: "Mechanist", isAllowed: and(includeEthic("Materialist", "Fanatic Materialist"), excludeCivic("Permanent Employment"))},
	{name: "Syncretic Evolution", isAllowed: and(excludeEthic("Gestalt Consciousness"), excludeCivic("Fanatic Purifiers"))},
	{name: "Life-Seeded", isAllowed: notAuth("Machine Intelligence")},
	{name: "Post-Apocalyptic", isAllowed: and(notAuth("Machine Intelligence"), excludeCivic("Agrarian Idyll", "Anglers"))},
	{name: "Remnants", isAllowed: excludeCivic("Agrarian Idyll")},
	{name: "Shattered Ring", isAllowed: excludeCivic("Agrarian Idyll", "Anglers")},
	{name: "Void Dwellers", isAllowed: and(notAuth("Machine Intelligence"), excludeCivic("Idyllic Bloom", "Agrarian Idyll", "Anglers"))},
	{name: "Scion", isAllowed: and(excludeEthic("Gestalt Consciousness", "Fanatic Xenophobe"), excludeCivic("Pompous Purists"))},
	{name: "Galactic Doorstep", isAllowed: always},
	{name: "Tree of Life", isAllowed: and(auth("Hive Mind"), excludeCivic("Devouring Swarm", "Terravore"))},
	{name: "On the Shoulders of Giants", isAllowed: excludeEthic("Gestalt Consciousness")},
	{name: "Calamitous Birth", isAllowed: and(excludeCivic("Catalytic Processing", "Organic Reprocessing", "Catalytic Recyclers"), notAuth("Machine Intelligence"))},
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

var allTraits = []Trait{
	{name: "Aquatic", cost: 1},
	{name: "Agrarian", cost: 2},
	{name: "Ingenious", cost: 2},
	{name: "Industrious", cost: 2},
	{name: ""},
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
