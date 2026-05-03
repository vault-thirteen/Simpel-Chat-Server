package helper

import (
	"slices"
)

func ArrayWithoutItemAt[S ~[]E, E any](s S, idx int) S {
	return slices.Delete(s, idx, idx+1)
}
