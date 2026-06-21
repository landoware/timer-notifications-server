package main

import "time"

type Crop struct {
	Name      string
	Value     string
	RLName    string
	Duration  time.Duration
	WikiTitle string
}

var cropsByGroup = map[CropGroup][]Crop{
	CropGroupAllotment: {
		{Name: "Potato", Value: "potato", RLName: "Potato", Duration: 40 * time.Minute, WikiTitle: "Potato"},
		{Name: "Onion", Value: "onion", RLName: "Onion", Duration: 40 * time.Minute, WikiTitle: "Onion"},
		{Name: "Cabbage", Value: "cabbage", RLName: "Cabbage", Duration: 40 * time.Minute, WikiTitle: "Cabbage"},
		{Name: "Tomato", Value: "tomato", RLName: "Tomato", Duration: 40 * time.Minute, WikiTitle: "Tomato"},
		{Name: "Sweetcorn", Value: "sweetcorn", RLName: "Sweetcorn", Duration: 1 * time.Hour, WikiTitle: "Sweetcorn"},
		{Name: "Strawberry", Value: "strawberry", RLName: "Strawberry", Duration: 1 * time.Hour, WikiTitle: "Strawberry"},
		{Name: "Watermelon", Value: "watermelon", RLName: "Watermelon", Duration: 80 * time.Minute, WikiTitle: "Watermelon"},
		{Name: "Snape Grass", Value: "snape_grass", RLName: "Snape grass", Duration: 70 * time.Minute, WikiTitle: "Snape grass"},
	},
	CropGroupFlower: {
		{Name: "Marigold", Value: "marigold", RLName: "Marigold", Duration: 20 * time.Minute, WikiTitle: "Marigold"},
		{Name: "Rosemary", Value: "rosemary", RLName: "Rosemary", Duration: 20 * time.Minute, WikiTitle: "Rosemary"},
		{Name: "Nasturtium", Value: "nasturtium", RLName: "Nasturtium", Duration: 20 * time.Minute, WikiTitle: "Nasturtium"},
		{Name: "Woad Leaves", Value: "woad", RLName: "Woad leaf", Duration: 20 * time.Minute, WikiTitle: "Woad leaf"},
		{Name: "Limpwurt Root", Value: "limpwurt", RLName: "Limpwurt", Duration: 20 * time.Minute, WikiTitle: "Limpwurt root"},
		{Name: "White Lily", Value: "white_lily", RLName: "White lily", Duration: 20 * time.Minute, WikiTitle: "White lily"},
	},
	CropGroupHerb: {
		{Name: "Guam", Value: "guam", RLName: "Guam", Duration: 80 * time.Minute, WikiTitle: "Guam"},
		{Name: "Marrentill", Value: "marrentill", RLName: "Marrentill", Duration: 80 * time.Minute, WikiTitle: "Marrentill"},
		{Name: "Tarromin", Value: "tarromin", RLName: "Tarromin", Duration: 80 * time.Minute, WikiTitle: "Tarromin"},
		{Name: "Harralander", Value: "harralander", RLName: "Harralander", Duration: 80 * time.Minute, WikiTitle: "Harralander"},
		{Name: "Goutweed", Value: "goutweed", RLName: "Goutweed", Duration: 80 * time.Minute, WikiTitle: "Goutweed"},
		{Name: "Ranarr", Value: "ranarr", RLName: "Ranarr", Duration: 80 * time.Minute, WikiTitle: "Ranarr"},
		{Name: "Toadflax", Value: "toadflax", RLName: "Toadflax", Duration: 80 * time.Minute, WikiTitle: "Toadflax"},
		{Name: "Irit", Value: "irit", RLName: "Irit", Duration: 80 * time.Minute, WikiTitle: "Irit"},
		{Name: "Avantoe", Value: "avantoe", RLName: "Avantoe", Duration: 80 * time.Minute, WikiTitle: "Avantoe"},
		{Name: "Kwuarm", Value: "kwuarm", RLName: "Kwuarm", Duration: 80 * time.Minute, WikiTitle: "Kwuarm"},
		{Name: "Snapdragon", Value: "snapdragon", RLName: "Snapdragon", Duration: 80 * time.Minute, WikiTitle: "Snapdragon"},
		{Name: "Huasca", Value: "huasca", RLName: "Huasca", Duration: 80 * time.Minute, WikiTitle: "Huasca"},
		{Name: "Cadantine", Value: "cadantine", RLName: "Cadantine", Duration: 80 * time.Minute, WikiTitle: "Cadantine"},
		{Name: "Lantadyme", Value: "lantadyme", RLName: "Lantadyme", Duration: 80 * time.Minute, WikiTitle: "Lantadyme"},
		{Name: "Dwarf Weed", Value: "dwarf_weed", RLName: "Dwarf weed", Duration: 80 * time.Minute, WikiTitle: "Dwarf weed"},
		{Name: "Torstol", Value: "torstol", RLName: "Torstol", Duration: 80 * time.Minute, WikiTitle: "Torstol"},
	},
	CropGroupHops: {
		{Name: "Barley", Value: "barley", RLName: "Barley", Duration: 40 * time.Minute, WikiTitle: "Barley"},
		{Name: "Hammerstone", Value: "hammerstone", RLName: "Hammerstone", Duration: 40 * time.Minute, WikiTitle: "Hammerstone hops"},
		{Name: "Asgarnian", Value: "asgarnian", RLName: "Asgarnian", Duration: 50 * time.Minute, WikiTitle: "Asgarnian hops"},
		{Name: "Jute", Value: "jute", RLName: "Jute", Duration: 50 * time.Minute, WikiTitle: "Jute fibre"},
		{Name: "Yanillian", Value: "yanillian", RLName: "Yanillian", Duration: 1 * time.Hour, WikiTitle: "Yanillian hops"},
		{Name: "Flax", Value: "flax", RLName: "Flax", Duration: 1 * time.Hour, WikiTitle: "Flax"},
		{Name: "Krandorian", Value: "krandorian", RLName: "Krandorian", Duration: 70 * time.Minute, WikiTitle: "Krandorian hops"},
		{Name: "Wildblood", Value: "wildblood", RLName: "Wildblood", Duration: 80 * time.Minute, WikiTitle: "Wildblood hops"},
		{Name: "Hemp", Value: "hemp", RLName: "Hemp", Duration: 80 * time.Minute, WikiTitle: "Hemp"},
		{Name: "Cotton", Value: "cotton", RLName: "Cotton boll", Duration: 100 * time.Minute, WikiTitle: "Cotton boll"},
	},
	CropGroupBush: {
		{Name: "Redberry", Value: "redberry", RLName: "Redberry", Duration: 100 * time.Minute, WikiTitle: "Redberries"},
		{Name: "Cadavaberry", Value: "cadavaberry", RLName: "Cadavaberry", Duration: 2 * time.Hour, WikiTitle: "Cadava berries"},
		{Name: "Dwellberry", Value: "dwellberry", RLName: "Dwellberry", Duration: 140 * time.Minute, WikiTitle: "Dwellberries"},
		{Name: "Jangerberry", Value: "jangerberry", RLName: "Jangerberry", Duration: 160 * time.Minute, WikiTitle: "Jangerberries"},
		{Name: "Whiteberry", Value: "whiteberry", RLName: "Whiteberry", Duration: 160 * time.Minute, WikiTitle: "White berries"},
		{Name: "Poison Ivy", Value: "poison_ivy", RLName: "Poison ivy", Duration: 160 * time.Minute, WikiTitle: "Poison ivy berries"},
	},
	CropGroupTree: {
		{Name: "Oak Tree", Value: "oak", RLName: "Oak", Duration: 160 * time.Minute, WikiTitle: "Oak tree (Farming)"},
		{Name: "Willow Tree", Value: "willow", RLName: "Willow", Duration: 4 * time.Hour, WikiTitle: "Willow tree (Farming)"},
		{Name: "Maple Tree", Value: "maple", RLName: "Maple", Duration: 320 * time.Minute, WikiTitle: "Maple tree (Farming)"},
		{Name: "Yew Tree", Value: "yew", RLName: "Yew", Duration: 400 * time.Minute, WikiTitle: "Yew tree (Farming)"},
		{Name: "Magic", Value: "magic", RLName: "Magic tree", Duration: 8 * time.Hour, WikiTitle: "Magic tree (Farming)"},
	},
	CropGroupFruitTree: {
		{Name: "Apple Tree", Value: "apple", RLName: "Apple", Duration: 16 * time.Hour, WikiTitle: "Apple tree"},
		{Name: "Banana Tree", Value: "banana", RLName: "Banana", Duration: 16 * time.Hour, WikiTitle: "Banana tree (Farming)"},
		{Name: "Orange Tree", Value: "orange", RLName: "Orange", Duration: 16 * time.Hour, WikiTitle: "Orange tree"},
		{Name: "Curry Tree", Value: "curry", RLName: "Curry", Duration: 16 * time.Hour, WikiTitle: "Curry tree"},
		{Name: "Pineapple Tree", Value: "pineapple", RLName: "Pineapple", Duration: 16 * time.Hour, WikiTitle: "Pineapple plant"},
		{Name: "Papaya Tree", Value: "papaya", RLName: "Papaya", Duration: 16 * time.Hour, WikiTitle: "Papaya tree"},
		{Name: "Palm Tree", Value: "palm", RLName: "Palm", Duration: 16 * time.Hour, WikiTitle: "Palm tree"},
		{Name: "Dragonfruit Tree", Value: "dragonfruit", RLName: "Dragonfruit", Duration: 16 * time.Hour, WikiTitle: "Dragonfruit tree"},
	},
	CropGroupCactus: {
		{Name: "Cactus", Value: "cactus", RLName: "Cactus", Duration: 560 * time.Minute, WikiTitle: "Cactus spine"},
		{Name: "Potato Cactus", Value: "potato_cactus", RLName: "Potato cactus", Duration: 70 * time.Minute, WikiTitle: "Potato cactus"},
	},
	CropGroupSeaweed: {
		{Name: "Giant Seaweed", Value: "giant_seaweed", RLName: "Giant seaweed", Duration: 40 * time.Minute, WikiTitle: "Giant seaweed"},
	},
	CropGroupMushroom: {
		{Name: "Mushroom", Value: "mushroom", RLName: "Mushroom", Duration: 4 * time.Hour, WikiTitle: "Mushroom"},
	},
	CropGroupBelladonna: {
		{Name: "Belladonna", Value: "belladonna", RLName: "Belladonna", Duration: 320 * time.Minute, WikiTitle: "Belladonna"},
	},
	CropGroupCalquat: {
		{Name: "Calquat", Value: "calquat", RLName: "Calquat", Duration: 1280 * time.Minute, WikiTitle: "Calquat tree"},
	},
	CropGroupCelastrus: {
		{Name: "Celastrus Tree", Value: "celastrus", RLName: "Celastrus", Duration: 800 * time.Minute, WikiTitle: "Celastrus tree"},
	},
	CropGroupRedwood: {
		{Name: "Redwood Tree", Value: "redwood", RLName: "Redwood", Duration: 6400 * time.Minute, WikiTitle: "Redwood tree (Farming)"},
	},
	CropGroupSpiritTree: {
		{Name: "Spirit Tree", Value: "spirit_tree", RLName: "Spirit tree (Farming)", Duration: 3840 * time.Minute, WikiTitle: "Spirit Tree (Farming)"},
	},
	CropGroupHespori: {
		{Name: "Hespori", Value: "hespori", RLName: "Hespori", Duration: 1920 * time.Minute, WikiTitle: "Hespori"},
	},
	CropGroupBirdhouse: {
		{Name: "Bird House", Value: "bird_house", RLName: "Bird house", Duration: 50 * time.Minute, WikiTitle: "Bird house"},
		{Name: "Oak Bird House", Value: "oak_bird_house", RLName: "Oak bird house", Duration: 50 * time.Minute, WikiTitle: "Oak bird house"},
		{Name: "Willow Bird House", Value: "willow_bird_house", RLName: "Willow bird house", Duration: 50 * time.Minute, WikiTitle: "Willow bird house"},
		{Name: "Teak Bird House", Value: "teak_bird_house", RLName: "Teak bird house", Duration: 50 * time.Minute, WikiTitle: "Teak bird house"},
		{Name: "Maple Bird House", Value: "maple_bird_house", RLName: "Maple bird house", Duration: 50 * time.Minute, WikiTitle: "Maple bird house"},
		{Name: "Mahogany Bird House", Value: "mahogany_bird_house", RLName: "Mahogany bird house", Duration: 50 * time.Minute, WikiTitle: "Mahogany bird house"},
		{Name: "Magic Bird House", Value: "magic_bird_house", RLName: "Magic bird house", Duration: 50 * time.Minute, WikiTitle: "Magic bird house"},
		{Name: "Redwood Bird House", Value: "redwood_bird_house", RLName: "Redwood bird house", Duration: 50 * time.Minute, WikiTitle: "Redwood bird house"},
	},
	CropGroupGrape: {
		{Name: "Grapes", Value: "grapes", RLName: "Grapes", Duration: 35 * time.Minute, WikiTitle: "Grapes"},
	},
	CropGroupAnima: {
		{Name: "Kronos", Value: "kronos", RLName: "Kronos", Duration: 5120 * time.Minute, WikiTitle: "Kronos plant"},
		{Name: "Iasor", Value: "iasor", RLName: "Iasor", Duration: 5120 * time.Minute, WikiTitle: "Iasor plant"},
		{Name: "Attas", Value: "attas", RLName: "Attas", Duration: 5120 * time.Minute, WikiTitle: "Attas plant"},
	},
	CropGroupHardwood: {
		{Name: "Teak tree", Value: "teak", RLName: "Teak", Duration: 4480 * time.Minute, WikiTitle: "Teak tree"},
		{Name: "Mahogany tree", Value: "mahogany", RLName: "Mahogany", Duration: 5120 * time.Minute, WikiTitle: "Mahogany tree (Farming)"},
		{Name: "Camphor tree", Value: "camphor", RLName: "Camphor", Duration: 5120 * time.Minute, WikiTitle: "Camphor tree"},
		{Name: "Ironwood tree", Value: "ironwood", RLName: "Ironwood", Duration: 5120 * time.Minute, WikiTitle: "Ironwood tree"},
		{Name: "Rosewood tree", Value: "rosewood", RLName: "Rosewood", Duration: 6400 * time.Minute, WikiTitle: "Rosewood tree"},
	},
	CropGroupCrystal: {
		{Name: "Crystal Tree", Value: "crystal", RLName: "Crystal tree", Duration: 480 * time.Minute, WikiTitle: "Crystal tree"},
	},
	CropGroupCoral: {
		{Name: "Elkhorn", Value: "elkhorn", RLName: "Elkhorn", Duration: 160 * time.Minute, WikiTitle: "Elkhorn coral"},
		{Name: "Pillar", Value: "pillar", RLName: "Pillar", Duration: 160 * time.Minute, WikiTitle: "Pillar coral"},
		{Name: "Umbral", Value: "umbral", RLName: "Umbral", Duration: 160 * time.Minute, WikiTitle: "Umbral coral"},
	},
	CropGroupFarmingContract: {
		// Allotment
		{Name: "Potato", Value: "potato", RLName: "Potato", Duration: 40 * time.Minute, WikiTitle: "Potato"},
		{Name: "Onion", Value: "onion", RLName: "Onion", Duration: 40 * time.Minute, WikiTitle: "Onion"},
		{Name: "Cabbage", Value: "cabbage", RLName: "Cabbage", Duration: 40 * time.Minute, WikiTitle: "Cabbage"},
		{Name: "Tomato", Value: "tomato", RLName: "Tomato", Duration: 40 * time.Minute, WikiTitle: "Tomato"},
		{Name: "Sweetcorn", Value: "sweetcorn", RLName: "Sweetcorn", Duration: 1 * time.Hour, WikiTitle: "Sweetcorn"},
		{Name: "Strawberry", Value: "strawberry", RLName: "Strawberry", Duration: 1 * time.Hour, WikiTitle: "Strawberry"},
		{Name: "Watermelon", Value: "watermelon", RLName: "Watermelon", Duration: 80 * time.Minute, WikiTitle: "Watermelon"},
		{Name: "Snape Grass", Value: "snape_grass", RLName: "Snape grass", Duration: 70 * time.Minute, WikiTitle: "Snape grass"},
		// Flower
		{Name: "Marigold", Value: "marigold", RLName: "Marigold", Duration: 20 * time.Minute, WikiTitle: "Marigold"},
		{Name: "Rosemary", Value: "rosemary", RLName: "Rosemary", Duration: 20 * time.Minute, WikiTitle: "Rosemary"},
		{Name: "Nasturtium", Value: "nasturtium", RLName: "Nasturtium", Duration: 20 * time.Minute, WikiTitle: "Nasturtium"},
		{Name: "Woad Leaves", Value: "woad", RLName: "Woad leaf", Duration: 20 * time.Minute, WikiTitle: "Woad leaf"},
		{Name: "Limpwurt Root", Value: "limpwurt", RLName: "Limpwurt", Duration: 20 * time.Minute, WikiTitle: "Limpwurt root"},
		{Name: "White Lily", Value: "white_lily", RLName: "White lily", Duration: 20 * time.Minute, WikiTitle: "White lily"},
		// Bush
		{Name: "Redberry", Value: "redberry", RLName: "Redberry", Duration: 100 * time.Minute, WikiTitle: "Redberries"},
		{Name: "Cadavaberry", Value: "cadavaberry", RLName: "Cadavaberry", Duration: 2 * time.Hour, WikiTitle: "Cadava berries"},
		{Name: "Dwellberry", Value: "dwellberry", RLName: "Dwellberry", Duration: 140 * time.Minute, WikiTitle: "Dwellberries"},
		{Name: "Jangerberry", Value: "jangerberry", RLName: "Jangerberry", Duration: 160 * time.Minute, WikiTitle: "Jangerberries"},
		{Name: "Whiteberry", Value: "whiteberry", RLName: "Whiteberry", Duration: 160 * time.Minute, WikiTitle: "White berries"},
		{Name: "Poison Ivy", Value: "poison_ivy", RLName: "Poison ivy", Duration: 160 * time.Minute, WikiTitle: "Poison ivy berries"},
		// Cactus
		{Name: "Cactus", Value: "cactus", RLName: "Cactus", Duration: 560 * time.Minute, WikiTitle: "Cactus spine"},
		{Name: "Potato Cactus", Value: "potato_cactus", RLName: "Potato cactus", Duration: 70 * time.Minute, WikiTitle: "Potato cactus"},
		// Herb
		{Name: "Guam", Value: "guam", RLName: "Guam", Duration: 80 * time.Minute, WikiTitle: "Guam"},
		{Name: "Marrentill", Value: "marrentill", RLName: "Marrentill", Duration: 80 * time.Minute, WikiTitle: "Marrentill"},
		{Name: "Tarromin", Value: "tarromin", RLName: "Tarromin", Duration: 80 * time.Minute, WikiTitle: "Tarromin"},
		{Name: "Harralander", Value: "harralander", RLName: "Harralander", Duration: 80 * time.Minute, WikiTitle: "Harralander"},
		{Name: "Ranarr", Value: "ranarr", RLName: "Ranarr", Duration: 80 * time.Minute, WikiTitle: "Ranarr"},
		{Name: "Toadflax", Value: "toadflax", RLName: "Toadflax", Duration: 80 * time.Minute, WikiTitle: "Toadflax"},
		{Name: "Irit", Value: "irit", RLName: "Irit", Duration: 80 * time.Minute, WikiTitle: "Irit"},
		{Name: "Avantoe", Value: "avantoe", RLName: "Avantoe", Duration: 80 * time.Minute, WikiTitle: "Avantoe"},
		{Name: "Kwuarm", Value: "kwuarm", RLName: "Kwuarm", Duration: 80 * time.Minute, WikiTitle: "Kwuarm"},
		{Name: "Snapdragon", Value: "snapdragon", RLName: "Snapdragon", Duration: 80 * time.Minute, WikiTitle: "Snapdragon"},
		{Name: "Cadantine", Value: "cadantine", RLName: "Cadantine", Duration: 80 * time.Minute, WikiTitle: "Cadantine"},
		{Name: "Lantadyme", Value: "lantadyme", RLName: "Lantadyme", Duration: 80 * time.Minute, WikiTitle: "Lantadyme"},
		{Name: "Dwarf Weed", Value: "dwarf_weed", RLName: "Dwarf weed", Duration: 80 * time.Minute, WikiTitle: "Dwarf weed"},
		{Name: "Torstol", Value: "torstol", RLName: "Torstol", Duration: 80 * time.Minute, WikiTitle: "Torstol"},
		// Tree
		{Name: "Oak Tree", Value: "oak", RLName: "Oak", Duration: 160 * time.Minute, WikiTitle: "Oak tree (Farming)"},
		{Name: "Willow Tree", Value: "willow", RLName: "Willow", Duration: 4 * time.Hour, WikiTitle: "Willow tree (Farming)"},
		{Name: "Maple Tree", Value: "maple", RLName: "Maple", Duration: 320 * time.Minute, WikiTitle: "Maple tree (Farming)"},
		{Name: "Yew Tree", Value: "yew", RLName: "Yew", Duration: 400 * time.Minute, WikiTitle: "Yew tree (Farming)"},
		{Name: "Magic Tree", Value: "magic", RLName: "Magic tree", Duration: 8 * time.Hour, WikiTitle: "Magic tree (Farming)"},
		// Fruit Tree
		{Name: "Apple Tree", Value: "apple", RLName: "Apple", Duration: 16 * time.Hour, WikiTitle: "Apple tree"},
		{Name: "Banana Tree", Value: "banana", RLName: "Banana", Duration: 16 * time.Hour, WikiTitle: "Banana tree (Farming)"},
		{Name: "Orange Tree", Value: "orange", RLName: "Orange", Duration: 16 * time.Hour, WikiTitle: "Orange tree"},
		{Name: "Curry Tree", Value: "curry", RLName: "Curry", Duration: 16 * time.Hour, WikiTitle: "Curry tree"},
		{Name: "Pineapple Tree", Value: "pineapple", RLName: "Pineapple", Duration: 16 * time.Hour, WikiTitle: "Pineapple plant"},
		{Name: "Papaya Tree", Value: "papaya", RLName: "Papaya", Duration: 16 * time.Hour, WikiTitle: "Papaya tree"},
		{Name: "Palm Tree", Value: "palm", RLName: "Palm", Duration: 16 * time.Hour, WikiTitle: "Palm tree"},
		{Name: "Dragonfruit Tree", Value: "dragonfruit", RLName: "Dragonfruit", Duration: 16 * time.Hour, WikiTitle: "Dragonfruit tree"},
		// Special
		{Name: "Celastrus Tree", Value: "celastrus", RLName: "Celastrus", Duration: 800 * time.Minute, WikiTitle: "Celastrus tree"},
		{Name: "Redwood Tree", Value: "redwood", RLName: "Redwood", Duration: 6400 * time.Minute, WikiTitle: "Redwood tree (Farming)"},
	},
}

func cropsForGroup(group CropGroup) []Crop {
	return cropsByGroup[group]
}

func cropForGroup(group CropGroup, value string) (Crop, bool) {
	for _, crop := range cropsByGroup[group] {
		if crop.Value == value || crop.RLName == value {
			return crop, true
		}
	}

	return Crop{}, false
}

func defaultCropForGroup(group CropGroup) (Crop, bool) {
	crops := cropsForGroup(group)
	if len(crops) == 0 {
		return Crop{}, false
	}

	return crops[0], true
}

func cropOptionRequired(group CropGroup) bool {
	crops := cropsForGroup(group)
	if len(crops) <= 1 {
		return false
	}

	firstDuration := crops[0].Duration
	for _, crop := range crops[1:] {
		if crop.Duration != firstDuration {
			return true
		}
	}

	return false
}

func gameModeDuration(d time.Duration, mode GameMode) time.Duration {
	switch mode {
	case GameModeLeagues, GameModeDeadman:
		return d / 5
	default:
		return d
	}
}
