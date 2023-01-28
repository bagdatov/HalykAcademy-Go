package office

import (
	"errors"
	"fmt"
)

type Location interface {
	CheckMoveToArea(Location) bool
	GetLocationTitle() string
}

type BaseLocationParams struct {
	title    string
	nearWith []string
}

var ErrUnknownLocationType = errors.New("unknown location type")

func NewLocationFactory(title string) (Location, error) {
	switch title {
	case "office":
		return newOffice(), nil
	case "workArea":
		return newWorkArea(), nil
	case "servers":
		return newServers(), nil
	}
	return nil, fmt.Errorf("newLocation: %w", ErrUnknownLocationType)
}

// TODO:: impl Location
type office struct {
	BaseLocationParams
}

func newOffice() Location {
	nearWith := []string{"workArea"}

	return &office{
		BaseLocationParams{
			title:    "office",
			nearWith: nearWith,
		},
	}
}

func (o *office) CheckMoveToArea(l Location) bool {
	locationName := l.GetLocationTitle()
	for _, neighbor := range o.nearWith {
		if locationName == neighbor {
			return true
		}
	}
	return false
}

func (o *office) GetLocationTitle() string {
	return o.title
}

// TODO:: impl Location
type workArea struct {
	BaseLocationParams
}

func newWorkArea() Location {
	nearWith := []string{"office", "servers"}

	return &workArea{
		BaseLocationParams{
			title:    "workArea",
			nearWith: nearWith,
		},
	}
}

func (w *workArea) CheckMoveToArea(l Location) bool {
	locationName := l.GetLocationTitle()
	for _, neighbor := range w.nearWith {
		if locationName == neighbor {
			return true
		}
	}
	return false
}

func (w *workArea) GetLocationTitle() string {
	return w.title
}

// TODO:: impl Location
type servers struct {
	BaseLocationParams
}

func newServers() Location {
	nearWith := []string{"workArea"}

	return &servers{
		BaseLocationParams{
			title:    "servers",
			nearWith: nearWith,
		},
	}
}

func (s *servers) CheckMoveToArea(l Location) bool {
	locationName := l.GetLocationTitle()
	for _, neighbor := range s.nearWith {
		if locationName == neighbor {
			return true
		}
	}
	return false
}

func (s *servers) GetLocationTitle() string {
	return s.title
}
