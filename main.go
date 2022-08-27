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
		Styles: []string{
			"/web/app.css", // Loads hello.css file.
		},
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func (d *data) Render() app.UI {
	return app.Div().Body(
		app.Button().Text("Generate").OnClick(d.generateEmpire),
		app.Div().Class("horizontal").Body(
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
					app.Br(),
					app.Div().Body(
						app.Span().Text("Main Species:"),
						app.Br(),
						app.Label().Text("Type").For("MainType"),
						app.Span().ID("MainType").Text(d.Empires[i].mainSpecies.popType),
						app.Ul().Body(app.Range(d.Empires[i].mainSpecies.traits).Slice(func(j int) app.UI {
							trait := d.Empires[i].mainSpecies.traits[j]
							return app.Li().Text(trait.name)
						})),
						app.If(len(d.Empires[i].subSpecies.traits) > 0, app.Div().Body(
							app.Span().Text("Sub Species:"),
							app.Br(),
							app.Label().Text("Type").For("SubType"),
							app.Span().ID("SubType").Text(d.Empires[i].subSpecies.popType),
							app.Ul().Body(app.Range(d.Empires[i].subSpecies.traits).Slice(func(j int) app.UI {
								trait := d.Empires[i].subSpecies.traits[j]
								return app.Li().Text(trait.name)
							})),
						)),
					),
					app.Br(),
					app.Br(),
				)
			}),
		))
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
		empire = generateSpecies(empire)
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

func generateSpecies(empire Empire) Empire {
	species := Species{}
	generateSubSpecies := false
	popTypes := []string{"Aquatic", "Mammalian", "Reptilian", "Avian", "Arthropoid", "Molluscoid", "Fungoid", "Plantoid", "Lithoid", "Necroid"}
	subspecies := Species{
		initialTraitPoints: 2,
	}
	for _, civic := range empire.civics {
		if civic.name == "Anglers" {
			species.traits = append(species.traits, originTraits["Aquatic"])
			subspecies.traits = append(subspecies.traits, originTraits["Aquatic"])
		}
		if civic.name == "Idyllic Bloom" {
			popTypes = []string{"Fungoid", "Plantoid"}
		}
		if civic.name == "Terravore" {
			popTypes = []string{"Lithoid"}
		}
	}
	if empire.authority == "Machine Intelligence" {
		//generate machine species
		species.popType = "Machine"
		species.initialTraitPoints = 1
	} else if empire.origin.name == "Calamitous Birth" {
		//generate lithoid species
		species.popType = "Lithoid"
		species.traits = append(species.traits, originTraits["Lithoid"])
		species.initialTraitPoints = 2
	} else if empire.origin.name == "Ocean Paradise" {
		//aquatic species, forced aqautic trait
		species.popType = "Aquatic"
		species.traits = append(species.traits, originTraits["Aquatic"])
		species.initialTraitPoints = 2
	} else {
		switch empire.origin.name {
		case "Clone Soldier":
			species.traits = append(species.traits, originTraits["Clone Soldier"])
		case "Post-Apocalyptic":
			species.traits = append(species.traits, originTraits["Survivor"])
		case "Void Dwellers":
			species.traits = append(species.traits, originTraits["Void Dweller"])
		case "Necrophage":
			species.traits = append(species.traits, originTraits["Necrophage"])
			generateSubSpecies = true
		case "Subterrenean":
			species.traits = append(species.traits, originTraits["Cave Dweller"])
		case "Syncretic Evolution":
			generateSubSpecies = true
			subspecies.traits = append(subspecies.traits, originTraits["Serviles"])
		}
		//standard species
		species.popType = popTypes[r.Intn(len(popTypes))]
		species.initialTraitPoints = 2
	}
	empire.mainSpecies = fillSpecies(species, empire.authority == "Hive Mind")
	if generateSubSpecies {
		subspecies.popType = popTypes[r.Intn(len(popTypes))]
		empire.subSpecies = fillSpecies(subspecies, empire.authority == "Hive Mind")
	}
	return empire
}

func fillSpecies(s Species, gestalt bool) Species {
	for {
		result, ok := singleSpeciesTry(s, gestalt)
		if ok {
			return result
		}
	}
}

func singleSpeciesTry(s Species, gestalt bool) (Species, bool) {
	traitCountOptions := []int{1, 2, 3, 3, 4, 4, 5, 5, 5}
	traitsToGenerate := traitCountOptions[r.Intn(len(traitCountOptions))]
	for i := 0; i < traitsToGenerate; i++ {
		traits := availableTraits(s, gestalt)
		s.traits = append(s.traits, traits[r.Intn(len(traits))])
	}
	res := s.initialTraitPoints
	for _, trait := range s.traits {
		res -= trait.cost
	}
	if res == 0 {
		return s, true
	}
	return Species{}, false
}

func availableTraits(s Species, gestalt bool) []Trait {
	result := []Trait{}
outer:
	for _, trait := range allTraits {
		if !trait.isAllowed(s) || (trait.nonGestalt && gestalt) {
			continue
		}
		for _, sTrait := range s.traits {
			if trait.name == sTrait.name {
				continue outer
			}
		}
		result = append(result, trait)
	}
	return result
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
	cost       int
	name       string
	nonGestalt bool
	isAllowed  speciesPredicate
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
	// {name: "Gestalt Consciousness", isAllowed: onlyGestalt},
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
	{name: "Terravore", isAllowed: and(auth("Hive Mind"), excludeCivic("Devouring Swarm", "Empath", "Idyllic Bloom")), genocidal: true},
	{name: "Divided Attention", isAllowed: auth("Hive Mind")},
	{name: "Empath", isAllowed: and(auth("Hive Mind"), excludeCivic("Terravore", "Devouring Swarm"))},
	{name: "Idyllic Bloom", isAllowed: and(auth("Hive Mind"), excludeCivic("Terravore"))},
	{name: "Memorialist", isAllowed: auth("Hive Mind")},
	{name: "Natural Neural Network", isAllowed: auth("Hive Mind")},
	{name: "One Mind", isAllowed: auth("Hive Mind")},
	{name: "Organic Reprocessing", isAllowed: auth("Hive Mind")},
	{name: "Pooled Knowledge", isAllowed: auth("Hive Mind")},
	{name: "Strength of Legions", isAllowed: auth("Hive Mind")},
	{name: "Subspace Ephase", isAllowed: auth("Hive Mind")},
	{name: "Subsumed Will", isAllowed: auth("Hive Mind")},
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
	{name: "Life-Seeded", isAllowed: and(notAuth("Machine Intelligence"), excludeCivic("Anglers"))},
	{name: "Post-Apocalyptic", isAllowed: and(notAuth("Machine Intelligence"), excludeCivic("Agrarian Idyll", "Anglers"))},
	{name: "Remnants", isAllowed: excludeCivic("Agrarian Idyll")},
	{name: "Shattered Ring", isAllowed: excludeCivic("Agrarian Idyll", "Anglers")},
	{name: "Void Dwellers", isAllowed: and(excludeEthic("Gestalt Consciousness"), excludeCivic("Idyllic Bloom", "Agrarian Idyll", "Anglers"))},
	{name: "Scion", isAllowed: and(excludeEthic("Gestalt Consciousness", "Fanatic Xenophobe"), excludeCivic("Pompous Purists"))},
	{name: "Galactic Doorstep", isAllowed: always},
	{name: "Tree of Life", isAllowed: and(auth("Hive Mind"), excludeCivic("Devouring Swarm", "Terravore"))},
	{name: "On the Shoulders of Giants", isAllowed: excludeEthic("Gestalt Consciousness")},
	{name: "Calamitous Birth", isAllowed: and(excludeCivic("Catalytic Processing", "Organic Reprocessing", "Catalytic Recyclers", "Devouring Swarm", "Idyllic Bloom"), notAuth("Machine Intelligence"))},
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
	{name: "Subterrenean", isAllowed: and(notAuth("Machine Intelligence"), excludeCivic("Anglers"))},
	{name: "Slingshot to the Stars", isAllowed: always},
	{name: "Teachers of the Shroud", isAllowed: and(includeEthic("Spiritualist", "Fanatic Spiritualist"), excludeCivic("Fanatic Purifiers"))},
	{name: "Imperial Fiefdom", isAllowed: excludeCivic("Inward Perfection", "Fanatic Purifiers", "Devouring Swarm", "Terravore", "Driven Assimilator", "Determined Exterminator")},
}

var allTraits = []Trait{
	{name: "Adaptive", cost: 2, isAllowed: andS(excludeType("Machine"), excludeTrait("Extremely Adaptive", "Nonadaptive", "Lithoid"))},
	{name: "Extremely Adaptive", cost: 4, isAllowed: andS(excludeType("Machine"), excludeTrait("Adaptive", "Nonadaptive", "Lithoid"))},
	{name: "Agrarian", cost: 2, isAllowed: andS(excludeType("Machine"), excludeTrait("Lithoid"))},
	{name: "Aquatic", cost: 1, isAllowed: andS(excludeType("Machine"), excludeTrait("Cave Dweller"))},
	{name: "Charismatic", cost: 2, isAllowed: andS(excludeType("Machine"), excludeTrait("Repugnant"))},
	{name: "Communal", cost: 1, isAllowed: andS(excludeType("Machine"), excludeTrait("Solitary"))},
	{name: "Conformists", cost: 2, isAllowed: andS(excludeType("Machine"), excludeTrait("Deviants")), nonGestalt: true},
	{name: "Conservationist", cost: 1, isAllowed: andS(excludeType("Machine"), excludeTrait("Wasteful"))},
	{name: "Docile", cost: 2, isAllowed: andS(excludeType("Machine"), excludeTrait("Unruly"))},
	{name: "Enduring", cost: 1, isAllowed: andS(excludeType("Machine"), excludeTrait("Fleeting", "Venerable"))},
	{name: "Venerable", cost: 4, isAllowed: andS(excludeType("Machine"), excludeTrait("Fleeting", "Enduring"))},
	{name: "Industrious", cost: 2, isAllowed: andS(excludeType("Machine"))},
	{name: "Ingenious", cost: 2, isAllowed: andS(excludeType("Machine"))},
	{name: "Intelligent", cost: 2, isAllowed: andS(excludeType("Machine"), excludeTrait("Serviles"))},
	{name: "Natural Engineers", cost: 1, isAllowed: andS(excludeType("Machine"), excludeTrait("Natural Physicists", "Natural Sociologists", "Serviles"))},
	{name: "Natural Physicists", cost: 1, isAllowed: andS(excludeType("Machine"), excludeTrait("Natural Engineers", "Natural Sociologists", "Serviles"))},
	{name: "Natural Sociologists", cost: 1, isAllowed: andS(excludeType("Machine"), excludeTrait("Natural Engineers", "Natural Physicists", "Serviles"))},
	{name: "Nomadic", cost: 1, isAllowed: andS(excludeType("Machine"), excludeTrait("Sedentary"))},
	{name: "Quick Learners", cost: 1, isAllowed: andS(excludeType("Machine"), excludeTrait("Slow Learners"))},
	{name: "Rapid Breeders", cost: 2, isAllowed: andS(excludeType("Machine"), excludeTrait("Slow Breeders", "Clone Soldier", "Lithoid"))},
	{name: "Resilient", cost: 1, isAllowed: andS(excludeType("Machine"))},
	{name: "Strong", cost: 1, isAllowed: andS(excludeType("Machine"), excludeTrait("Very Strong", "Weak"))},
	{name: "Very Strong", cost: 3, isAllowed: andS(excludeType("Machine"), excludeTrait("Strong", "Weak"))},
	{name: "Talented", cost: 1, isAllowed: andS(excludeType("Machine"))},
	{name: "Thrifty", cost: 2, isAllowed: andS(excludeType("Machine")), nonGestalt: true},
	{name: "Traditional", cost: 1, isAllowed: andS(excludeType("Machine"), excludeTrait("Quarrelsome"))},
	{name: "Nonadaptive", cost: -2, isAllowed: andS(excludeType("Machine"), excludeTrait("Adaptive", "Extremely Adaptive", "Lithoid"))},
	{name: "Repugnant", cost: -2, isAllowed: andS(excludeType("Machine"), excludeTrait("Charismatic"))},
	{name: "Solitary", cost: -1, isAllowed: andS(excludeType("Machine"), excludeTrait("Communal"))},
	{name: "Deviants", cost: -1, isAllowed: andS(excludeType("Machine"), excludeTrait("Conformists")), nonGestalt: true},
	{name: "Wasteful", cost: -1, isAllowed: andS(excludeType("Machine"), excludeTrait("Conservationist"))},
	{name: "Unruly", cost: -2, isAllowed: andS(excludeType("Machine"), excludeTrait("Docile"))},
	{name: "Fleeting", cost: -1, isAllowed: andS(excludeType("Machine"), excludeTrait("Enduring", "Venerable"))},
	{name: "Sedentary", cost: -1, isAllowed: andS(excludeType("Machine"), excludeTrait("Nomadic"))},
	{name: "Slow Learners", cost: -1, isAllowed: andS(excludeType("Machine"), excludeTrait("Quick Learners"))},
	{name: "Slow Breeders", cost: -2, isAllowed: andS(excludeType("Machine"), excludeTrait("Rapid Breeders", "Lithoid", "Clone Soldier"))},
	{name: "Weak", cost: -1, isAllowed: andS(excludeType("Machine"), excludeTrait("Strong", "Very Strong"))},
	{name: "Quarrelsome", cost: -1, isAllowed: andS(excludeType("Machine"), excludeTrait("Traditional"))},
	{name: "Decadent", cost: -1, isAllowed: andS(excludeType("Machine")), nonGestalt: true},
	{name: "Phototropic", cost: 1, isAllowed: andS(includeType("Plantoid", "Fungoid"), excludeTrait("Radiotropic", "Cave Dweller"))},
	{name: "Radiotropic", cost: 2, isAllowed: andS(includeType("Plantoid", "Fungoid"), excludeTrait("Phototropic"))},
	{name: "Budding", cost: 2, isAllowed: andS(includeType("Plantoid", "Fungoid"), excludeTrait("Slow Breeders", "Rapid Breeders", "Clone Soldier", "Necrophage"))},
	{name: "Gaseous Byproducts", cost: 2, isAllowed: andS(includeType("Lithoid"), excludeTrait("Scintillating Skin", "Volatile Excretions"))},
	{name: "Scintillating Skin", cost: 2, isAllowed: andS(includeType("Lithoid"), excludeTrait("Gaseous Byproducts", "Volatile Excretions"))},
	{name: "Volatile Excretions", cost: 2, isAllowed: andS(includeType("Lithoid"), excludeTrait("Gaseous Byproducts", "Scintillating Skin"))},
	{name: "Double Jointed", cost: 1, isAllowed: andS(includeType("Machine"), excludeTrait("Bulky"))},
	{name: "Durable", cost: 1, isAllowed: andS(includeType("Machine"), excludeTrait("High Maintenance"))},
	{name: "Efficient Processors", cost: 3, isAllowed: andS(includeType("Machine"))},
	{name: "Emotion Emulators", cost: 1, isAllowed: andS(includeType("Machine"), excludeTrait("Uncanny"))},
	{name: "Enhanced Memory", cost: 2, isAllowed: andS(includeType("Machine"))},
	{name: "Logic Engines", cost: 2, isAllowed: andS(includeType("Machine"))},
	{name: "Mass-Produced", cost: 1, isAllowed: andS(includeType("Machine"), excludeTrait("Custom-Made"))},
	{name: "Power Drills", cost: 2, isAllowed: andS(includeType("Machine"))},
	{name: "Recycled", cost: 2, isAllowed: andS(includeType("Machine"), excludeTrait("Luxurious"))},
	{name: "Streamlined Protocols", cost: 2, isAllowed: andS(includeType("Machine"), excludeTrait("High Bandwidth"))},
	{name: "Superconducive", cost: 2, isAllowed: andS(includeType("Machine"))},
	{name: "Bulky", cost: -1, isAllowed: andS(includeType("Machine"), excludeTrait("Double Jointed"))},
	{name: "High Maintenance", cost: -1, isAllowed: andS(includeType("Machine"), excludeTrait("Durable"))},
	{name: "Uncanny", cost: -1, isAllowed: andS(includeType("Machine"), excludeTrait("Emotion Emulators"))},
	{name: "Custom-Made", cost: -1, isAllowed: andS(includeType("Machine"), excludeTrait("Mass-Produced"))},
	{name: "Luxurious", cost: -2, isAllowed: andS(includeType("Machine"), excludeTrait("Recycled"))},
	{name: "High Bandwidth", cost: -2, isAllowed: andS(includeType("Machine"), excludeTrait("Streamlined Protocols"))},
	{name: "Learning Algorithms", cost: 1, isAllowed: andS(includeType("Machine"), excludeTrait("Repurposed Hardware"))},
	{name: "Repurposed Hardware", cost: -1, isAllowed: andS(includeType("Machine"), excludeTrait("Learning Algorithms"))},
}

var originTraits = make(map[string]Trait)

func init() {
	originTraits["Lithoid"] = Trait{name: "Lithoid", cost: 0, isAllowed: never}
	originTraits["Serviles"] = Trait{name: "Serviles", cost: 1, isAllowed: never}
	originTraits["Clone Soldier"] = Trait{name: "Clone Soldier", cost: 0, isAllowed: never}
	originTraits["Survivor"] = Trait{name: "Survivor", cost: 0, isAllowed: never}
	originTraits["Void Dweller"] = Trait{name: "Void Dweller", cost: 0, isAllowed: never}
	originTraits["Necrophage"] = Trait{name: "Necrophage", cost: 0, isAllowed: never}
	originTraits["Cave Dweller"] = Trait{name: "Cave Dweller", cost: 0, isAllowed: never}
	originTraits["Aquatic"] = Trait{name: "Aquatic", cost: 1, isAllowed: andS(excludeType("Machine"), excludeTrait("Cave Dweller"))}
}

func never(s Species) bool {
	return false
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

func excludeTrait(s ...string) speciesPredicate {
	return func(species Species) bool {
		for _, trait := range s {
			for _, sTrait := range species.traits {
				if trait == sTrait.name {
					return false
				}
			}
		}
		return true
	}
}

func includeType(s ...string) speciesPredicate {
	return func(species Species) bool {
		for _, popType := range s {
			if species.popType == popType {
				return true
			}
		}
		return false
	}
}

func excludeType(s ...string) speciesPredicate {
	return func(species Species) bool {
		for _, popType := range s {
			if species.popType == popType {
				return false
			}
		}
		return true
	}
}

func andS(s ...speciesPredicate) speciesPredicate {
	return func(species Species) bool {
		for _, pred := range s {
			if !pred(species) {
				return false
			}
		}
		return true
	}
}
