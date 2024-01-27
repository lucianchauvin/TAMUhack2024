package main


import(
	"fmt"	
	"os"
	"strings"
	list "container/list"
	tea "github.com/charmbracelet/bubbletea"
)

func toRune(i int) rune {
    return rune('A' - 1 + i)
}

type transState struct {
    nextState int
    write int
    direction bool // true -> move right
}

type model struct {
	tape	*list.List
	head	*list.Element
    state   int
    stateTable [][]transState 
}

func initialModel() model {
	initTape := list.New()
	initTape.PushBack(1)
	return model{
		tape: initTape,
		head: initTape.Front(),
		state: 2,
		stateTable: [][]transState{},
	}
}

func (m model) step() model { 
	if m.head.Next() == nil {
		m.tape.InsertAfter(0,m.head)
	}
	m.head = m.head.Next()
    //curValue := m.head.Value.(int)
    //curState := m.state

    //m.head.Value = m.stateTable[curState][curValue].write
    //m.state = m.stateTable[curState][curValue].nextState
    //if(m.stateTable[curState][curValue].direction){
    //    m.head = m.head.Next()
    //}else{
    //    m.head = m.head.Prev()
    //}
	return m
}


func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("T-soding")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			m = m.step()
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "The current tape state:\n\n"

	viewWidth := 4;
	viewHead := m.head;

	for i:=0; i < viewWidth; i++ {
		if viewHead.Prev() == nil {
			m.tape.InsertBefore(0,viewHead)
		}
		viewHead = viewHead.Prev();
	}
	for i:=0; i < 2*viewWidth+1; i++ {
		s += fmt.Sprintf(" %v", viewHead.Value)
		if viewHead.Next() == nil {
			m.tape.InsertAfter(0,viewHead)
		}
		viewHead = viewHead.Next();
	}
	s += "\n"

		//s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	s += strings.Repeat(" ",2*viewWidth+1)
	s += "^\n"
	s += fmt.Sprintf("Current state: %q\n", toRune(m.state))

	s += "\nPress s to step, q to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
