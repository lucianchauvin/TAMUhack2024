package main


import (
	"fmt"	
	"os"
	list "container/list"
	tea "github.com/charmbracelet/bubbletea"
)

type transState struct {
    nextState int
    direction bool // right is true
}

type model struct {
	tape	list.List
	head	*list.Element
    state   int
    stateTable [][]transState 
}

func initialModel() model {
	initTape = list.New()
	return model{
		tape: initTape,
		head: initTape.Front(),
	}
}

func (m model) step(){ 
//    curValue := m.head.Value
//    curValue := m.head.Value
//    curState := m.state
//    m.head.Value = 
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
		case "s":
			m.step()
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
	s += fmt.Sprintf("Current state: %v\n", m.state)

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
