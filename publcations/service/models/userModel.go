package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Groups []Group
}

func (u User) IsInGroup(group_id uint) bool {
	for _, group := range u.Groups {
		if group.ID == group_id {
			return true
		}
	}
	return false
}

// returns map[group_id]it_is_a_parent; keys contains only groups upper at hierarchy
func (u User) GetAllParentGroups() (visited map[uint]bool) {
	visited = make(map[uint]bool)
	for _, g := range u.Groups {
		recoursiveGroupBypas(g, visited)
	}
	return visited
}
func recoursiveGroupBypas(group Group, visited map[uint]bool) {
	if visited[group.ID] {
		return
	}
	for _, g := range group.Parents {
		if visited[g.ID] {
			continue
		}
		visited[g.ID] = true
		recoursiveGroupBypas(g, visited)
	}
}
