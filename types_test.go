package main

import (
	"sort"
	"testing"
)

func TestCropGroupValidate(t *testing.T) {
	tests := []struct {
		name    string
		group   CropGroup
		wantErr bool
	}{
		{"valid herb", CropGroupHerb, false},
		{"valid tree", CropGroupTree, false},
		{"valid fruit tree", CropGroupFruitTree, false},
		{"valid allotment", CropGroupAllotment, false},
		{"valid flower", CropGroupFlower, false},
		{"valid bush", CropGroupBush, false},
		{"valid hops", CropGroupHops, false},
		{"valid cactus", CropGroupCactus, false},
		{"valid seaweed", CropGroupSeaweed, false},
		{"valid mushroom", CropGroupMushroom, false},
		{"valid belladonna", CropGroupBelladonna, false},
		{"valid calquat", CropGroupCalquat, false},
		{"valid celastrus", CropGroupCelastrus, false},
		{"valid redwood", CropGroupRedwood, false},
		{"valid spirit tree", CropGroupSpiritTree, false},
		{"valid hespori", CropGroupHespori, false},
		{"valid birdhouse", CropGroupBirdhouse, false},
		{"invalid empty", CropGroup(""), true},
		{"invalid gibberish", CropGroup("not_a_crop"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.group.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestCropGroupDisplayNamePlural(t *testing.T) {
	tests := []struct {
		group CropGroup
		want  string
	}{
		{CropGroupHerb, "herbs"},
		{CropGroupTree, "trees"},
		{CropGroupFruitTree, "fruit trees"},
		{CropGroupAllotment, "allotments"},
		{CropGroupFlower, "flowers"},
		{CropGroupBush, "bushes"},
		{CropGroupHops, "hops"},
		{CropGroupCactus, "cactus"},
		{CropGroupSeaweed, "seaweeds"},
		{CropGroupMushroom, "mushrooms"},
		{CropGroupBelladonna, "belladonnas"},
		{CropGroupCalquat, "calquats"},
		{CropGroupCelastrus, "celastrus"},
		{CropGroupRedwood, "redwoods"},
		{CropGroupSpiritTree, "spirit trees"},
		{CropGroupHespori, "hesporis"},
		{CropGroupBirdhouse, "birdhouses"},
	}

	for _, tt := range tests {
		t.Run(string(tt.group), func(t *testing.T) {
			got := tt.group.DisplayNamePlural()
			if got != tt.want {
				t.Errorf("DisplayNamePlural() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCropGroupDisplayName(t *testing.T) {
	tests := []struct {
		group CropGroup
		want  string
	}{
		{CropGroupHerb, "Herb"},
		{CropGroupTree, "Tree"},
		{CropGroupFruitTree, "Fruit Tree"},
		{CropGroupAllotment, "Allotment"},
		{CropGroupBirdhouse, "Birdhouse"},
		{CropGroupSpiritTree, "Spirit Tree"},
		{CropGroupBelladonna, "Belladonna"},
	}

	for _, tt := range tests {
		t.Run(string(tt.group), func(t *testing.T) {
			got := tt.group.DisplayName()
			if got != tt.want {
				t.Errorf("DisplayName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSortedCropGroups_AllGroupsPresent(t *testing.T) {
	got := sortedCropGroups()

	if len(got) != len(validCropGroups) {
		t.Errorf("sortedCropGroups() returned %d groups, want %d", len(got), len(validCropGroups))
	}

	seen := make(map[CropGroup]bool)
	for _, g := range got {
		if _, ok := validCropGroups[g]; !ok {
			t.Errorf("unexpected group %q in sorted output", g)
		}
		if seen[g] {
			t.Errorf("duplicate group %q in sorted output", g)
		}
		seen[g] = true
	}
}

func TestSortedCropGroups_IsSorted(t *testing.T) {
	got := sortedCropGroups()
	if !sort.SliceIsSorted(got, func(i, j int) bool { return got[i] < got[j] }) {
		t.Errorf("sortedCropGroups() output is not sorted: %v", got)
	}
}

func TestAllowedCropGroups_AllGroupsPresent(t *testing.T) {
	got := allowedCropGroups()
	if len(got) != len(validCropGroups) {
		t.Errorf("allowedCropGroups() returned %d groups, want %d", len(got), len(validCropGroups))
	}

	for _, g := range got {
		if _, ok := validCropGroups[CropGroup(g)]; !ok {
			t.Errorf("unexpected group %q in allowed output", g)
		}
	}
}

func TestAllowedCropGroups_IsSorted(t *testing.T) {
	got := allowedCropGroups()
	if !sort.StringsAreSorted(got) {
		t.Errorf("allowedCropGroups() output is not sorted: %v", got)
	}
}

func TestPatchLocationValidate(t *testing.T) {
	if !PatchLocationFarmingGuild.Validate() {
		t.Error("expected Farming Guild to be valid")
	}
	if !PatchLocationAlKharid.Validate() {
		t.Error("expected Al Kharid to be valid")
	}
	if PatchLocation("").Validate() {
		t.Error("expected empty location to be invalid")
	}
	if PatchLocation("Nowhere").Validate() {
		t.Error("expected unknown location to be invalid")
	}
}

func TestValidPatchLocationsForGroup_Herb(t *testing.T) {
	locs := ValidPatchLocationsForGroup(CropGroupHerb)
	if len(locs) == 0 {
		t.Fatal("expected non-empty list")
	}
	seen := make(map[PatchLocation]bool)
	for _, loc := range locs {
		if seen[loc] {
			t.Errorf("duplicate location %q", loc)
		}
		seen[loc] = true
		if !loc.Validate() {
			t.Errorf("location %q should be valid", loc)
		}
	}

	if !containsLocation(locs, PatchLocationArdougne) {
		t.Error("expected Ardougne in herb locations")
	}
}

func TestValidPatchLocationsForGroup_Birdhouse(t *testing.T) {
	locs := ValidPatchLocationsForGroup(CropGroupBirdhouse)
	if locs != nil {
		t.Errorf("expected nil for birdhouse, got %v", locs)
	}
}

func containsLocation(locs []PatchLocation, loc PatchLocation) bool {
	for _, l := range locs {
		if l == loc {
			return true
		}
	}
	return false
}
