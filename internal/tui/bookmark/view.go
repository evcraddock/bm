package bookmarktui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	if m.windowSize != nil {
		height := m.windowSize.Height - marginHeight
		docStyle = docStyle.Height(height)
	}

	return docStyle.Render(b.String())
}
