package ui

import (
	"fmt"
	"os"
	"container/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	tape	list
	head	*element
}

func initialModel() model {
	return model{
		tape: list.New(),
		head: list.front(),
	}
}


func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Grocery List")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "The current tape state:\n\n"

	viewWidth = 4;
	viewHead = m.head;

	for i:=0; i < viewWidth; i++ {
		if viewHead.Prev() == nil {
			m.list.insertBefore(0,viewHead)
		}
		viewHead = viewHead.Prev();
	}
	for i:=0; i < 2*viewWidth+1; i++ {
		s += fmt.Sprintf(" %v", viewHead.value)
		if viewHead.next == nil {
			m.list.insertAfter(0,viewHead)
		}
		viewHead = viewHead.next();
	}

		//s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	s += " "*2*viewWidth;
	s += "^\n"

	s += "\nPress q to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
