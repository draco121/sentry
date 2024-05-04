package core

import (
	"github.com/draco121/common/constants"
	"slices"
)

func authorizationEngine(allowedActions []constants.Action, requiredActions []constants.Action) bool {
	if len(requiredActions) == 0 || requiredActions == nil {
		return false
	}
	if slices.Contains(allowedActions, constants.All) {
		return true
	} else {
		for i := range requiredActions {
			if !slices.Contains(allowedActions, requiredActions[i]) {
				return false
			}
		}
		return true
	}
}
