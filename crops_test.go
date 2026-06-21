package main

import (
	"testing"
	"time"
)

func TestCropsForGroup(t *testing.T) {
	t.Run("returns crops for valid group", func(t *testing.T) {
		crops := cropsForGroup(CropGroupSeaweed)
		if len(crops) == 0 {
			t.Fatal("expected at least one crop for seaweed")
		}
		if crops[0].Name != "Giant Seaweed" {
			t.Errorf("first seaweed crop = %q, want %q", crops[0].Name, "Giant Seaweed")
		}
	})

	t.Run("returns nil for unknown group", func(t *testing.T) {
		crops := cropsForGroup(CropGroup("invalid"))
		if crops != nil {
			t.Errorf("expected nil, got %v", crops)
		}
	})
}

func TestCropForGroup(t *testing.T) {
	t.Run("finds existing crop by Value", func(t *testing.T) {
		crop, ok := cropForGroup(CropGroupHerb, "ranarr")
		if !ok {
			t.Fatal("expected to find ranarr")
		}
		if crop.Name != "Ranarr" {
			t.Errorf("crop.Name = %q, want %q", crop.Name, "Ranarr")
		}
		if crop.Duration != 80*time.Minute {
			t.Errorf("crop.Duration = %v, want %v", crop.Duration, 80*time.Minute)
		}
		if crop.WikiTitle != "Ranarr" {
			t.Errorf("crop.WikiTitle = %q, want %q", crop.WikiTitle, "Ranarr")
		}
	})

	t.Run("finds existing crop by RLName", func(t *testing.T) {
		crop, ok := cropForGroup(CropGroupHerb, "Ranarr")
		if !ok {
			t.Fatal("expected to find Ranarr by RLName")
		}
		if crop.Value != "ranarr" {
			t.Errorf("crop.Value = %q, want %q", crop.Value, "ranarr")
		}
	})

	t.Run("finds tree by RLName (Oak)", func(t *testing.T) {
		crop, ok := cropForGroup(CropGroupTree, "Oak")
		if !ok {
			t.Fatal("expected to find Oak by RLName")
		}
		if crop.Value != "oak" {
			t.Errorf("crop.Value = %q, want %q", crop.Value, "oak")
		}
	})

	t.Run("not found for missing value", func(t *testing.T) {
		_, ok := cropForGroup(CropGroupHerb, "nonexistent")
		if ok {
			t.Error("expected not found")
		}
	})

	t.Run("not found for invalid group", func(t *testing.T) {
		_, ok := cropForGroup(CropGroup("invalid"), "ranarr")
		if ok {
			t.Error("expected not found for invalid group")
		}
	})
}

func TestDefaultCropForGroup(t *testing.T) {
	t.Run("returns first crop for herb", func(t *testing.T) {
		crop, ok := defaultCropForGroup(CropGroupHerb)
		if !ok {
			t.Fatal("expected to find default crop")
		}
		if crop.Name != "Guam" {
			t.Errorf("default herb = %q, want %q", crop.Name, "Guam")
		}
	})

	t.Run("returns first crop for seaweed (single-crop group)", func(t *testing.T) {
		crop, ok := defaultCropForGroup(CropGroupSeaweed)
		if !ok {
			t.Fatal("expected to find default crop")
		}
		if crop.Name != "Giant Seaweed" {
			t.Errorf("default seaweed = %q, want %q", crop.Name, "Giant Seaweed")
		}
	})

	t.Run("returns false for invalid group", func(t *testing.T) {
		_, ok := defaultCropForGroup(CropGroup("invalid"))
		if ok {
			t.Error("expected not found for invalid group")
		}
	})
}

func TestCropOptionRequired(t *testing.T) {
	tests := []struct {
		name  string
		group CropGroup
		want  bool
	}{
		{"seaweed (single crop)", CropGroupSeaweed, false},
		{"mushroom (single crop)", CropGroupMushroom, false},
		{"belladonna (single crop)", CropGroupBelladonna, false},
		{"calquat (single crop)", CropGroupCalquat, false},
		{"celastrus (single crop)", CropGroupCelastrus, false},
		{"redwood (single crop)", CropGroupRedwood, false},
		{"spirit tree (single crop)", CropGroupSpiritTree, false},
		{"hespori (single crop)", CropGroupHespori, false},
		{"flower (uniform duration)", CropGroupFlower, false},
		{"herb (uniform duration)", CropGroupHerb, false},
		{"tree (mixed durations)", CropGroupTree, true},
		{"fruit tree (uniform duration)", CropGroupFruitTree, false},
		{"allotment (mixed durations)", CropGroupAllotment, true},
		{"hops (mixed durations)", CropGroupHops, true},
		{"bush (mixed durations)", CropGroupBush, true},
		{"cactus (mixed durations)", CropGroupCactus, true},
		{"birdhouse (uniform duration)", CropGroupBirdhouse, false},
		{"grape (single crop)", CropGroupGrape, false},
		{"anima (uniform duration)", CropGroupAnima, false},
		{"hardwood (mixed durations)", CropGroupHardwood, true},
		{"crystal (single crop)", CropGroupCrystal, false},
		{"coral (uniform duration)", CropGroupCoral, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cropOptionRequired(tt.group)
			if got != tt.want {
				t.Errorf("cropOptionRequired(%q) = %v, want %v", tt.group, got, tt.want)
			}
		})
	}
}

func TestCropDurationPrecision(t *testing.T) {
	for group, crops := range cropsByGroup {
		for _, crop := range crops {
			if crop.Duration <= 0 {
				t.Errorf("%s/%s has non-positive duration %v", group, crop.Name, crop.Duration)
			}
			if crop.Name == "" {
				t.Errorf("%s has a crop with empty Name", group)
			}
			if crop.Value == "" {
				t.Errorf("%s/%s has empty Value", group, crop.Name)
			}
			if crop.RLName == "" {
				t.Errorf("%s/%s has empty RLName", group, crop.Name)
			}
		}
	}
}
