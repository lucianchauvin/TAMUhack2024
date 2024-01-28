package main


import(
	"fmt"	
	"os"
	list "container/list"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

const HALT = 0

type model struct {
	tape		*list.List
	head		*list.Element
    stateTable	[][]transState 
    state		int
	viewWidth	int //number of steps run
	
}

func initialModel() model {
	initTape := list.New()
	initTape.PushBack(1)
	return model{
		tape: initTape,
		head: initTape.Front(),
		state: 1,
		stateTable: [][]transState{{{1,1,true}, {1,0,false}},{{1,1,false},{0,0,true}}},
		viewWidth: 10,
	}
}

func (m model) step() model { 
	if m.state == HALT {
		return m
	}
    curValue := m.head.Value.(int)
    curState := m.state

    m.head.Value = m.stateTable[curState][curValue].write
    m.state = m.stateTable[curState][curValue].nextState
    if(m.stateTable[curState][curValue].direction){
        m.head = m.head.Next()
    }else{
        m.head = m.head.Prev()
    }
	return m
}


func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("T-soding")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n", "s":
			m = m.step()
		case "+":
			if m.viewWidth < 20 {
				m.viewWidth++
			}
		case "-":
			if m.viewWidth > 2 {
				m.viewWidth--
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "The current tape state:\n\n"

	viewHead := m.head;

	for i:=0; i < m.viewWidth; i++ {
		if viewHead.Prev() == nil {
			m.tape.InsertBefore(0,viewHead)
		}
		viewHead = viewHead.Prev();
	}
	for i:=0; i < 2*m.viewWidth+1; i++ {
		s += fmt.Sprintf(" %v", viewHead.Value)
		if viewHead.Next() == nil {
			m.tape.InsertAfter(0,viewHead)
		}
		viewHead = viewHead.Next();
	}
	s += " \n"

		//s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	s += "^\n"
	s += fmt.Sprintf("Current state: %v\n", toRune(m.state))

	s += "\nPress s to step, +/- to see more/less of the tape, q to quit.\n"

	return style.Render(s)
}

var style = lipgloss.NewStyle().
    //Foreground(lipgloss.Color("#FAFAFA")).
    //Background(lipgloss.Color("#7D56F4")).
    PaddingTop(2).
    PaddingLeft(4).
    Width(100).
	Align(lipgloss.Center)


func toRune(i int) string {
	if i == HALT {
		return "HALT"
	}
    return string(rune('A' + i))
}

type transState struct {
    nextState int
    write int
    direction bool // true -> move right
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
