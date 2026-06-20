package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type CropGroup string

const (
	CropGroupHerb       CropGroup = "herb"
	CropGroupTree       CropGroup = "tree"
	CropGroupFruitTree  CropGroup = "fruit_tree"
	CropGroupAllotment  CropGroup = "allotment"
	CropGroupFlower     CropGroup = "flower"
	CropGroupBush       CropGroup = "bush"
	CropGroupHops       CropGroup = "hops"
	CropGroupCactus     CropGroup = "cactus"
	CropGroupSeaweed    CropGroup = "seaweed"
	CropGroupMushroom   CropGroup = "mushroom"
	CropGroupBelladonna CropGroup = "belladonna"
	CropGroupCalquat    CropGroup = "calquat"
	CropGroupCelastrus  CropGroup = "celastrus"
	CropGroupRedwood    CropGroup = "redwood"
	CropGroupSpiritTree CropGroup = "spirit_tree"
)

var validCropGroups = map[CropGroup]struct{}{
	CropGroupHerb:       {},
	CropGroupTree:       {},
	CropGroupFruitTree:  {},
	CropGroupAllotment:  {},
	CropGroupFlower:     {},
	CropGroupBush:       {},
	CropGroupHops:       {},
	CropGroupCactus:     {},
	CropGroupSeaweed:    {},
	CropGroupMushroom:   {},
	CropGroupBelladonna: {},
	CropGroupCalquat:    {},
	CropGroupCelastrus:  {},
	CropGroupRedwood:    {},
	CropGroupSpiritTree: {},
}

type NotificationRequest struct {
	UserID          string    `json:"userId"`
	CropGroup       CropGroup `json:"cropGroup,omitempty"`
	NotifyInMinutes int       `json:"notifyInMinutes"`
}

type NotificationResponse struct {
	UserID       string    `json:"userId"`
	CropGroup    CropGroup `json:"cropGroup"`
	ScheduledFor time.Time `json:"scheduledFor"`
	Status       string    `json:"status"`
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
