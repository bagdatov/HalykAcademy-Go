[![Open in Visual Studio Code](https://classroom.github.com/assets/open-in-vscode-f059dc9a6f8d3a56e377f745f24479a46679e63a5d9fe6f495e02850cd0d8118.svg)](https://classroom.github.com/online_ide?assignment_repo_id=5765752&assignment_repo_type=AssignmentRepo)

## What was done

To suit contract of Employee interface, it is required to implement two methods.

```go 
type Employee interface {
	GetCurrentLocation() Location
	MoveToLocation(Location) error
}
```

Let's see it on hr struct example.

```go 
type hr struct {
	baseEmployeeParams
}

type baseEmployeeParams struct {
	location Location
	accesses []string
}
```
First method is quite easy. We are taking required information from field location.

```go
func (h *hr) GetCurrentLocation() Location {
	return h.location
}
```

With the second method, we first see if location has been assigned at all. <br>
If not, it must be changed to office as per task requirements. <br>
Then, we look if we can move from current location to target room. <br>
Afterwards, we validate if employee has an authority to enter target room. <br>
Finally, location field of our hr is updated. <br>

```go
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

```

And with Location interface, it was similar. Except, some changes in task code.

```go
type Location interface {
	CheckMoveToArea(Location) bool 
	GetLocationTitle() string // this method was added by myself
}

func newOffice() Location {
	nearWith := []string{"workArea"}

	return &office{
		BaseLocationParams{
			title:    "office", // this field was also added
			nearWith: nearWith,
		},
	}
}
```
Let's see on office example.

```go
type office struct {
	BaseLocationParams
}

type BaseLocationParams struct {
	title    string
	nearWith []string
}
```
Getting location name through private field of struct. <br>
Then, we iterate through neighbors of the room to validate move.

```go
func (o *office) GetLocationTitle() string {
	return o.title
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
```