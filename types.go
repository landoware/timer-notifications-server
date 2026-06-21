package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type PatchInfo struct {
	Crop     string        `json:"crop"`
	Location PatchLocation `json:"location"`
}

type PatchLocation string

const (
	PatchLocationAlKharid         PatchLocation = "Al Kharid"
	PatchLocationAldarin          PatchLocation = "Aldarin"
	PatchLocationAnglersRetreat   PatchLocation = "Anglers' Retreat"
	PatchLocationArdougne         PatchLocation = "Ardougne"
	PatchLocationArdougneSouth    PatchLocation = "Ardougne (South)"
	PatchLocationAuburnvale       PatchLocation = "Auburnvale"
	PatchLocationBrimhaven        PatchLocation = "Brimhaven"
	PatchLocationCanifis          PatchLocation = "Canifis"
	PatchLocationCoralNursery     PatchLocation = "Coral Nurseries"
	PatchLocationCatherby         PatchLocation = "Catherby"
	PatchLocationCatherbyEast     PatchLocation = "Catherby (East)"
	PatchLocationChampionsGuild   PatchLocation = "Champions' Guild"
	PatchLocationDraynorManor     PatchLocation = "Draynor Manor"
	PatchLocationEntrana          PatchLocation = "Entrana"
	PatchLocationEtceteria        PatchLocation = "Etceteria"
	PatchLocationFalador          PatchLocation = "Falador"
	PatchLocationFaladorPark      PatchLocation = "Falador Park"
	PatchLocationFarmingGuild     PatchLocation = "Farming Guild"
	PatchLocationFossilIsland     PatchLocation = "Fossil Island"
	PatchLocationGnomeStronghold  PatchLocation = "Gnome Stronghold"
	PatchLocationHarmonyIsland    PatchLocation = "Harmony Island"
	PatchLocationHosidius         PatchLocation = "Hosidius"
	PatchLocationKastori          PatchLocation = "Kastori"
	PatchLocationLletya           PatchLocation = "Lletya"
	PatchLocationLocusOasis       PatchLocation = "Locus Oasis"
	PatchLocationLumbridge        PatchLocation = "Lumbridge"
	PatchLocationMcGruborsWood    PatchLocation = "McGrubor's Wood"
	PatchLocationNemusRetreat     PatchLocation = "Nemus Retreat"
	PatchLocationOrtusFarm        PatchLocation = "Ortus Farm"
	PatchLocationPortPhasmatys    PatchLocation = "Port Phasmatys"
	PatchLocationPortSarim        PatchLocation = "Port Sarim"
	PatchLocationPrifddinas       PatchLocation = "Prifddinas"
	PatchLocationRimmington       PatchLocation = "Rimmington"
	PatchLocationSummerShore      PatchLocation = "Summer Shore"
	PatchLocationTaiBwoWannai     PatchLocation = "Tai Bwo Wannai"
	PatchLocationTaverley         PatchLocation = "Taverley"
	PatchLocationTreeGnomeVillage PatchLocation = "Tree Gnome Village"
	PatchLocationTrollStronghold  PatchLocation = "Troll Stronghold"
	PatchLocationVarrock          PatchLocation = "Varrock"
	PatchLocationWeiss            PatchLocation = "Weiss"
	PatchLocationYanille          PatchLocation = "Yanille"
)

func (p PatchLocation) Validate() bool {
	_, ok := validPatchLocations[p]
	return ok
}

var validPatchLocations = map[PatchLocation]struct{}{
	PatchLocationAlKharid:         {},
	PatchLocationAldarin:          {},
	PatchLocationAnglersRetreat:   {},
	PatchLocationArdougne:         {},
	PatchLocationArdougneSouth:    {},
	PatchLocationAuburnvale:       {},
	PatchLocationBrimhaven:        {},
	PatchLocationCanifis:          {},
	PatchLocationCatherby:         {},
	PatchLocationCatherbyEast:     {},
	PatchLocationCoralNursery:     {},
	PatchLocationChampionsGuild:   {},
	PatchLocationDraynorManor:     {},
	PatchLocationEntrana:          {},
	PatchLocationEtceteria:        {},
	PatchLocationFalador:          {},
	PatchLocationFaladorPark:      {},
	PatchLocationFarmingGuild:     {},
	PatchLocationFossilIsland:     {},
	PatchLocationGnomeStronghold:  {},
	PatchLocationHarmonyIsland:    {},
	PatchLocationHosidius:         {},
	PatchLocationKastori:          {},
	PatchLocationLletya:           {},
	PatchLocationLocusOasis:       {},
	PatchLocationLumbridge:        {},
	PatchLocationMcGruborsWood:    {},
	PatchLocationNemusRetreat:     {},
	PatchLocationOrtusFarm:        {},
	PatchLocationPortPhasmatys:    {},
	PatchLocationPortSarim:        {},
	PatchLocationPrifddinas:       {},
	PatchLocationRimmington:       {},
	PatchLocationSummerShore:      {},
	PatchLocationTaiBwoWannai:     {},
	PatchLocationTaverley:         {},
	PatchLocationTreeGnomeVillage: {},
	PatchLocationTrollStronghold:  {},
	PatchLocationVarrock:          {},
	PatchLocationWeiss:            {},
	PatchLocationYanille:          {},
}

// ValidPatchLocationsForGroup returns the known patch locations for a given crop group.
func ValidPatchLocationsForGroup(group CropGroup) []PatchLocation {
	switch group {
	case CropGroupHerb:
		return []PatchLocation{
			PatchLocationArdougne,
			PatchLocationCatherby,
			PatchLocationFalador,
			PatchLocationFarmingGuild,
			PatchLocationHarmonyIsland,
			PatchLocationHosidius,
			PatchLocationOrtusFarm,
			PatchLocationPortPhasmatys,
			PatchLocationTrollStronghold,
			PatchLocationWeiss,
		}
	case CropGroupAllotment:
		return []PatchLocation{
			PatchLocationArdougne,
			PatchLocationCatherby,
			PatchLocationFalador,
			PatchLocationFarmingGuild,
			PatchLocationHarmonyIsland,
			PatchLocationHosidius,
			PatchLocationOrtusFarm,
			PatchLocationPortPhasmatys,
			PatchLocationPrifddinas,
		}
	case CropGroupFlower:
		return []PatchLocation{
			PatchLocationArdougne,
			PatchLocationCatherby,
			PatchLocationFalador,
			PatchLocationFarmingGuild,
			PatchLocationHosidius,
			PatchLocationKastori,
			PatchLocationOrtusFarm,
			PatchLocationPortPhasmatys,
			PatchLocationPrifddinas,
		}
	case CropGroupHops:
		return []PatchLocation{
			PatchLocationAldarin,
			PatchLocationEntrana,
			PatchLocationLumbridge,
			PatchLocationMcGruborsWood,
			PatchLocationYanille,
		}
	case CropGroupBush:
		return []PatchLocation{
			PatchLocationArdougneSouth,
			PatchLocationChampionsGuild,
			PatchLocationEtceteria,
			PatchLocationFarmingGuild,
			PatchLocationRimmington,
		}
	case CropGroupTree:
		return []PatchLocation{
			PatchLocationFaladorPark,
			PatchLocationFarmingGuild,
			PatchLocationGnomeStronghold,
			PatchLocationLumbridge,
			PatchLocationNemusRetreat,
			PatchLocationTaverley,
			PatchLocationVarrock,
		}
	case CropGroupFruitTree:
		return []PatchLocation{
			PatchLocationBrimhaven,
			PatchLocationCatherbyEast,
			PatchLocationFarmingGuild,
			PatchLocationGnomeStronghold,
			PatchLocationKastori,
			PatchLocationLletya,
			PatchLocationTreeGnomeVillage,
		}
	case CropGroupSpiritTree:
		return []PatchLocation{
			PatchLocationBrimhaven,
			PatchLocationEtceteria,
			PatchLocationFarmingGuild,
			PatchLocationHosidius,
			PatchLocationPortSarim,
		}
	case CropGroupCactus:
		return []PatchLocation{
			PatchLocationAlKharid,
			PatchLocationFarmingGuild,
		}
	case CropGroupSeaweed:
		return []PatchLocation{
			PatchLocationFossilIsland,
		}
	case CropGroupMushroom:
		return []PatchLocation{
			PatchLocationCanifis,
		}
	case CropGroupBelladonna:
		return []PatchLocation{
			PatchLocationAuburnvale,
			PatchLocationDraynorManor,
		}
	case CropGroupHespori:
		return []PatchLocation{
			PatchLocationFarmingGuild,
		}
	case CropGroupCalquat:
		return []PatchLocation{
			PatchLocationKastori,
			PatchLocationSummerShore,
			PatchLocationTaiBwoWannai,
		}
	case CropGroupCelastrus:
		return []PatchLocation{
			PatchLocationFarmingGuild,
		}
	case CropGroupRedwood:
		return []PatchLocation{
			PatchLocationFarmingGuild,
		}
	case CropGroupGrape:
		return []PatchLocation{
			PatchLocationHosidius,
		}
	case CropGroupAnima:
		return []PatchLocation{
			PatchLocationFarmingGuild,
		}
	case CropGroupHardwood:
		return []PatchLocation{
			PatchLocationAnglersRetreat,
			PatchLocationFossilIsland,
			PatchLocationLocusOasis,
		}
	case CropGroupCrystal:
		return []PatchLocation{
			PatchLocationPrifddinas,
		}
	case CropGroupCoral:
		return []PatchLocation{
			PatchLocationCoralNursery,
		}
	case CropGroupFarmingContract:
		return []PatchLocation{
			PatchLocationFarmingGuild,
		}
	default:
		return nil
	}
}

type CropGroup string

const (
	CropGroupHerb            CropGroup = "herb"
	CropGroupTree            CropGroup = "tree"
	CropGroupFruitTree       CropGroup = "fruit_tree"
	CropGroupAllotment       CropGroup = "allotment"
	CropGroupFlower          CropGroup = "flower"
	CropGroupBush            CropGroup = "bush"
	CropGroupHops            CropGroup = "hops"
	CropGroupCactus          CropGroup = "cactus"
	CropGroupSeaweed         CropGroup = "seaweed"
	CropGroupMushroom        CropGroup = "mushroom"
	CropGroupBelladonna      CropGroup = "belladonna"
	CropGroupCalquat         CropGroup = "calquat"
	CropGroupCelastrus       CropGroup = "celastrus"
	CropGroupRedwood         CropGroup = "redwood"
	CropGroupSpiritTree      CropGroup = "spirit_tree"
	CropGroupHespori         CropGroup = "hespori"
	CropGroupBirdhouse       CropGroup = "birdhouse"
	CropGroupGrape           CropGroup = "grape"
	CropGroupAnima           CropGroup = "anima"
	CropGroupHardwood        CropGroup = "hardwood"
	CropGroupCrystal         CropGroup = "crystal"
	CropGroupCoral           CropGroup = "coral"
	CropGroupFarmingContract CropGroup = "farming_contract"
)

var validCropGroups = map[CropGroup]struct{}{
	CropGroupHerb:            {},
	CropGroupTree:            {},
	CropGroupFruitTree:       {},
	CropGroupAllotment:       {},
	CropGroupFlower:          {},
	CropGroupBush:            {},
	CropGroupHops:            {},
	CropGroupCactus:          {},
	CropGroupSeaweed:         {},
	CropGroupMushroom:        {},
	CropGroupBelladonna:      {},
	CropGroupCalquat:         {},
	CropGroupCelastrus:       {},
	CropGroupRedwood:         {},
	CropGroupSpiritTree:      {},
	CropGroupHespori:         {},
	CropGroupBirdhouse:       {},
	CropGroupGrape:           {},
	CropGroupAnima:           {},
	CropGroupHardwood:        {},
	CropGroupCrystal:         {},
	CropGroupCoral:           {},
	CropGroupFarmingContract: {},
}

type NotifyMode string

const (
	NotifyModeFirstReady NotifyMode = "first_ready"
	NotifyModeAllReady   NotifyMode = "all_ready"
)

var validNotifyModes = map[NotifyMode]struct{}{
	NotifyModeFirstReady: {},
	NotifyModeAllReady:   {},
}

func (n NotifyMode) Validate() error {
	if _, ok := validNotifyModes[n]; !ok {
		return fmt.Errorf("invalid notifyMode %q", n)
	}

	return nil
}

type GameMode string

const (
	GameModeStandard GameMode = "standard"
	GameModeLeagues  GameMode = "leagues"
	GameModeDeadman  GameMode = "deadman"
)

var validGameModes = map[GameMode]struct{}{
	GameModeStandard: {},
	GameModeLeagues:  {},
	GameModeDeadman:  {},
}

func (g GameMode) Validate() error {
	if _, ok := validGameModes[g]; !ok {
		return fmt.Errorf("invalid gameMode %q", g)
	}

	return nil
}

type NotificationRequest struct {
	UserID          string      `json:"userId"`
	CropGroup       CropGroup   `json:"cropGroup,omitempty"`
	NotifyInMinutes int         `json:"notifyInMinutes"`
	CropName        string      `json:"-"`
	CropValue       string      `json:"crop,omitempty"`
	GameMode        GameMode    `json:"gameMode,omitempty"`
	NotifyMode      NotifyMode  `json:"notifyMode,omitempty"`
	Patches         []PatchInfo `json:"patches,omitempty"`
}

type NotificationResponse struct {
	UserID       string      `json:"userId"`
	CropGroup    CropGroup   `json:"cropGroup"`
	ScheduledFor time.Time   `json:"scheduledFor"`
	Status       string      `json:"status"`
	GameMode     GameMode    `json:"gameMode,omitempty"`
	NotifyMode   NotifyMode  `json:"notifyMode,omitempty"`
	Patches      []PatchInfo `json:"patches,omitempty"`
}

func (c CropGroup) Validate() error {
	if _, ok := validCropGroups[c]; !ok {
		return fmt.Errorf("invalid cropGroup %q", c)
	}

	return nil
}

func allowedCropGroups() []string {
	groups := make([]string, 0, len(validCropGroups))
	for group := range validCropGroups {
		groups = append(groups, string(group))
	}
	sort.Strings(groups)

	return groups
}

func sortedCropGroups() []CropGroup {
	groups := make([]CropGroup, 0, len(validCropGroups))
	for group := range validCropGroups {
		groups = append(groups, group)
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i] < groups[j]
	})

	return groups
}

func (c CropGroup) DisplayName() string {
	parts := strings.Split(string(c), "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}

	return strings.Join(parts, " ")
}

func (c CropGroup) DisplayNamePlural() string {
	name := c.DisplayName()
	lower := strings.ToLower(name)

	if strings.HasSuffix(lower, "s") {
		return lower
	}

	if strings.HasSuffix(lower, "sh") || strings.HasSuffix(lower, "ch") || strings.HasSuffix(lower, "x") || strings.HasSuffix(lower, "z") {
		return lower + "es"
	}

	return lower + "s"
}

func (c CropGroup) DisplayNamePluralTitle() string {
	name := c.DisplayName()
	lower := strings.ToLower(name)
	if strings.HasSuffix(lower, "s") {
		return name
	}
	if strings.HasSuffix(lower, "sh") || strings.HasSuffix(lower, "ch") || strings.HasSuffix(lower, "x") || strings.HasSuffix(lower, "z") {
		return name + "es"
	}
	return name + "s"
}
