package ui

// SelectChoice for a SelectPopUp box.
type SelectChoice struct {
	Message string
	Value   interface{}
}

// SelectPopUp is a scrollable box that floats over other boxes, showing
// selection choices.
type SelectPopUp struct {
	Choices []SelectChoice
}

// Selection returns the selection a user made.
func (s *SelectPopUp) Selection() (*SelectChoice, bool) {
	return nil, false
}

// Abort removes the popup box from the screen.
func (s *SelectPopUp) Abort() {

}
