package office

import (
	"errors"
	"fmt"
)

type Employee interface {
	GetCurrentLocation() Location
	MoveToLocation(Location) error
}

type baseEmployeeParams struct {
	location Location
	accesses []string
}

var ErrUnknownEmplType = errors.New("unknown employee type")

func NewEmployeeFactory(title string) (Employee, error) {
	switch title {
	case "hr":
		return newHr(), nil
	case "itSecurity":
		return newItSecurity(), nil
	}
	return nil, fmt.Errorf("newEmployee: %w", ErrUnknownEmplType)
}

// TODO:: must impl Employee
type hr struct {
	baseEmployeeParams
}

func newHr() Employee {
	accesses := []string{"office", "workArea"}

	return &hr{
		baseEmployeeParams{
			accesses: accesses,
		},
	}
}

func (h *hr) GetCurrentLocation() Location {
	return h.location
}

func (h *hr) MoveToLocation(l Location) error {
	if h.location == nil {
		newLoc, err := NewLocationFactory("office")

		if err != nil {
			return err
		}

		h.location = newLoc
		return nil
	}

	locationName := l.GetLocationTitle()

	if !h.location.CheckMoveToArea(l) {
		return fmt.Errorf("%T cannot enter %s from %s", h, locationName, h.location.GetLocationTitle())
	}

	var hasAccess bool
	for _, access := range h.accesses {
		if access == locationName {
			hasAccess = true
		}
	}

	if !hasAccess {
		return fmt.Errorf("%T does not have access to %s", h, locationName)
	}

	newLoc, err := NewLocationFactory(locationName)
	if err != nil {
		return err
	}

	h.location = newLoc
	return nil
}

// TODO:: must impl Employee
type itSecurity struct {
	baseEmployeeParams
}

func newItSecurity() Employee {
	accesses := []string{"office", "workArea", "servers"}

	return &itSecurity{
		baseEmployeeParams{
			accesses: accesses,
		},
	}
}

func (i *itSecurity) GetCurrentLocation() Location {
	return i.location
}

func (i *itSecurity) MoveToLocation(l Location) error {
	if i.location == nil {
		newLoc, err := NewLocationFactory("office")

		if err != nil {
			return err
		}

		i.location = newLoc
		return nil
	}

	locationName := l.GetLocationTitle()

	if !i.location.CheckMoveToArea(l) {
		return fmt.Errorf("%T cannot enter %s from %s", i, locationName, i.location.GetLocationTitle())
	}

	var hasAccess bool
	for _, access := range i.accesses {
		if access == locationName {
			hasAccess = true
		}
	}

	if !hasAccess {
		return fmt.Errorf("%T does not have access to %s", i, locationName)
	}

	newLoc, err := NewLocationFactory(locationName)
	if err != nil {
		return err
	}

	i.location = newLoc
	return nil
}
