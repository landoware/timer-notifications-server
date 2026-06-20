package main

import "time"

type Crop struct {
	Name      string
	Value     string
	Duration  time.Duration
	WikiTitle string
}

var cropsByGroup = map[CropGroup][]Crop{
	CropGroupAllotment: {
		{Name: "Potato", Value: "potato", Duration: 40 * time.Minute, WikiTitle: "Potato"},
		{Name: "Onion", Value: "onion", Duration: 40 * time.Minute, WikiTitle: "Onion"},
		{Name: "Cabbage", Value: "cabbage", Duration: 40 * time.Minute, WikiTitle: "Cabbage"},
		{Name: "Tomato", Value: "tomato", Duration: 40 * time.Minute, WikiTitle: "Tomato"},
		{Name: "Sweetcorn", Value: "sweetcorn", Duration: 1 * time.Hour, WikiTitle: "Sweetcorn"},
		{Name: "Strawberry", Value: "strawberry", Duration: 1 * time.Hour, WikiTitle: "Strawberry"},
		{Name: "Watermelon", Value: "watermelon", Duration: 80 * time.Minute, WikiTitle: "Watermelon"},
		{Name: "Snape Grass", Value: "snape_grass", Duration: 70 * time.Minute, WikiTitle: "Snape grass"},
	},
	CropGroupFlower: {
		{Name: "Marigold", Value: "marigold", Duration: 20 * time.Minute, WikiTitle: "Marigold"},
		{Name: "Rosemary", Value: "rosemary", Duration: 20 * time.Minute, WikiTitle: "Rosemary"},
		{Name: "Nasturtium", Value: "nasturtium", Duration: 20 * time.Minute, WikiTitle: "Nasturtium"},
		{Name: "Woad", Value: "woad", Duration: 20 * time.Minute, WikiTitle: "Woad leaf"},
		{Name: "Limpwurt", Value: "limpwurt", Duration: 20 * time.Minute, WikiTitle: "Limpwurt root"},
		{Name: "White Lily", Value: "white_lily", Duration: 20 * time.Minute, WikiTitle: "White lily"},
	},
	CropGroupHerb: {
		{Name: "Guam", Value: "guam", Duration: 80 * time.Minute, WikiTitle: "Guam"},
		{Name: "Marrentill", Value: "marrentill", Duration: 80 * time.Minute, WikiTitle: "Marrentill"},
		{Name: "Tarromin", Value: "tarromin", Duration: 80 * time.Minute, WikiTitle: "Tarromin"},
		{Name: "Harralander", Value: "harralander", Duration: 80 * time.Minute, WikiTitle: "Harralander"},
		{Name: "Goutweed", Value: "goutweed", Duration: 80 * time.Minute, WikiTitle: "Goutweed"},
		{Name: "Ranarr", Value: "ranarr", Duration: 80 * time.Minute, WikiTitle: "Ranarr"},
		{Name: "Toadflax", Value: "toadflax", Duration: 80 * time.Minute, WikiTitle: "Toadflax"},
		{Name: "Irit", Value: "irit", Duration: 80 * time.Minute, WikiTitle: "Irit"},
		{Name: "Avantoe", Value: "avantoe", Duration: 80 * time.Minute, WikiTitle: "Avantoe"},
		{Name: "Kwuarm", Value: "kwuarm", Duration: 80 * time.Minute, WikiTitle: "Kwuarm"},
		{Name: "Snapdragon", Value: "snapdragon", Duration: 80 * time.Minute, WikiTitle: "Snapdragon"},
		{Name: "Huasca", Value: "huasca", Duration: 80 * time.Minute, WikiTitle: "Huasca"},
		{Name: "Cadantine", Value: "cadantine", Duration: 80 * time.Minute, WikiTitle: "Cadantine"},
		{Name: "Lantadyme", Value: "lantadyme", Duration: 80 * time.Minute, WikiTitle: "Lantadyme"},
		{Name: "Dwarf Weed", Value: "dwarf_weed", Duration: 80 * time.Minute, WikiTitle: "Dwarf weed"},
		{Name: "Torstol", Value: "torstol", Duration: 80 * time.Minute, WikiTitle: "Torstol"},
	},
	CropGroupHops: {
		{Name: "Barley", Value: "barley", Duration: 40 * time.Minute, WikiTitle: "Barley"},
		{Name: "Hammerstone", Value: "hammerstone", Duration: 40 * time.Minute, WikiTitle: "Hammerstone hops"},
		{Name: "Asgarnian", Value: "asgarnian", Duration: 50 * time.Minute, WikiTitle: "Asgarnian hops"},
		{Name: "Jute", Value: "jute", Duration: 50 * time.Minute, WikiTitle: "Jute fibre"},
		{Name: "Yanillian", Value: "yanillian", Duration: 1 * time.Hour, WikiTitle: "Yanillian hops"},
		{Name: "Flax", Value: "flax", Duration: 1 * time.Hour, WikiTitle: "Flax"},
		{Name: "Krandorian", Value: "krandorian", Duration: 70 * time.Minute, WikiTitle: "Krandorian hops"},
		{Name: "Wildblood", Value: "wildblood", Duration: 80 * time.Minute, WikiTitle: "Wildblood hops"},
		{Name: "Hemp", Value: "hemp", Duration: 80 * time.Minute, WikiTitle: "Hemp"},
		{Name: "Cotton", Value: "cotton", Duration: 100 * time.Minute, WikiTitle: "Cotton boll"},
	},
	CropGroupBush: {
		{Name: "Redberry", Value: "redberry", Duration: 100 * time.Minute, WikiTitle: "Redberries"},
		{Name: "Cadavaberry", Value: "cadavaberry", Duration: 2 * time.Hour, WikiTitle: "Cadava berries"},
		{Name: "Dwellberry", Value: "dwellberry", Duration: 140 * time.Minute, WikiTitle: "Dwellberries"},
		{Name: "Jangerberry", Value: "jangerberry", Duration: 160 * time.Minute, WikiTitle: "Jangerberries"},
		{Name: "Whiteberry", Value: "whiteberry", Duration: 160 * time.Minute, WikiTitle: "White berries"},
		{Name: "Poison Ivy", Value: "poison_ivy", Duration: 160 * time.Minute, WikiTitle: "Poison ivy berries"},
	},
	CropGroupTree: {
		{Name: "Oak", Value: "oak", Duration: 160 * time.Minute, WikiTitle: "Oak tree (Farming)"},
		{Name: "Willow", Value: "willow", Duration: 4 * time.Hour, WikiTitle: "Willow tree (Farming)"},
		{Name: "Maple", Value: "maple", Duration: 320 * time.Minute, WikiTitle: "Maple tree (Farming)"},
		{Name: "Yew", Value: "yew", Duration: 400 * time.Minute, WikiTitle: "Yew tree (Farming)"},
		{Name: "Magic", Value: "magic", Duration: 8 * time.Hour, WikiTitle: "Magic tree (Farming)"},
	},
	CropGroupFruitTree: {
		{Name: "Apple", Value: "apple", Duration: 16 * time.Hour, WikiTitle: "Apple tree"},
		{Name: "Banana", Value: "banana", Duration: 16 * time.Hour, WikiTitle: "Banana tree (Farming)"},
		{Name: "Orange", Value: "orange", Duration: 16 * time.Hour, WikiTitle: "Orange tree"},
		{Name: "Curry", Value: "curry", Duration: 16 * time.Hour, WikiTitle: "Curry tree"},
		{Name: "Pineapple", Value: "pineapple", Duration: 16 * time.Hour, WikiTitle: "Pineapple plant"},
		{Name: "Papaya", Value: "papaya", Duration: 16 * time.Hour, WikiTitle: "Papaya tree"},
		{Name: "Palm", Value: "palm", Duration: 16 * time.Hour, WikiTitle: "Palm tree"},
		{Name: "Dragonfruit", Value: "dragonfruit", Duration: 16 * time.Hour, WikiTitle: "Dragonfruit tree"},
	},
	CropGroupCactus: {
		{Name: "Cactus", Value: "cactus", Duration: 560 * time.Minute, WikiTitle: "Cactus spine"},
		{Name: "Potato Cactus", Value: "potato_cactus", Duration: 70 * time.Minute, WikiTitle: "Potato cactus"},
	},
	CropGroupSeaweed: {
		{Name: "Giant Seaweed", Value: "giant_seaweed", Duration: 40 * time.Minute, WikiTitle: "Giant seaweed"},
	},
	CropGroupMushroom: {
		{Name: "Mushroom", Value: "mushroom", Duration: 4 * time.Hour, WikiTitle: "Mushroom"},
	},
	CropGroupBelladonna: {
		{Name: "Belladonna", Value: "belladonna", Duration: 320 * time.Minute, WikiTitle: "Belladonna"},
	},
	CropGroupCalquat: {
		{Name: "Calquat", Value: "calquat", Duration: 1280 * time.Minute, WikiTitle: "Calquat tree"},
	},
	CropGroupCelastrus: {
		{Name: "Celastrus", Value: "celastrus", Duration: 800 * time.Minute, WikiTitle: "Celastrus tree"},
	},
	CropGroupRedwood: {
		{Name: "Redwood", Value: "redwood", Duration: 6400 * time.Minute, WikiTitle: "Redwood tree (Farming)"},
	},
	CropGroupSpiritTree: {
		{Name: "Spirit Tree", Value: "spirit_tree", Duration: 3840 * time.Minute, WikiTitle: "Spirit Tree (Farming)"},
	},
}

func cropsForGroup(group CropGroup) []Crop {
	return cropsByGroup[group]
}

func cropForGroup(group CropGroup, value string) (Crop, bool) {
	for _, crop := range cropsByGroup[group] {
		if crop.Value == value {
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
