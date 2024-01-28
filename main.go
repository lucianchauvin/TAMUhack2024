package main

import(
	"fmt"	
	"os"
	list "container/list"
    "tsoding/table"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	//"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
)

const ACCEPT = 0
const REJECT = -1

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	tape		*list.List
	head		*list.Element
    stateTable	[][]transState 
    table		table.Model
    state		int
	viewWidth	int //number of steps run
	started		bool // you can't edit the table if you've started!
}


func initialModel() model {
    columns := []table.Column{
		{Title: "State/Symbol", Width: 12},
		{Title: "1", Width: 8},
		{Title: "2", Width: 8},
		{Title: "3", Width: 8},
	}

	rows := []table.Row{
		{"A", "aaa", "bbb", "ccc"},
		{"B", "Tokyo", "Japan", "ddd"},
		{"C", "Tokyo", "Japan", "eee"},
		{"D", "Tokyo", "Japan", "fff"},
		{"E", "Tokyo", "Japan", "37,274,000"},
		{"F", "Tokyo", "Japan", "37,274,000"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
    s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	s.Cell = s.Cell.
		BorderForeground(lipgloss.Color("240"))
	t.SetStyles(s)

	initTape := list.New()
	initTape.PushBack(1)
	return model{
		tape: initTape,
		head: initTape.Front(),
		stateTable: [][]transState{{{1,1,true}, {1,0,false}},{{1,1,false},{0,0,true}}},
		state:		1,
		viewWidth:	10,
		started:	false,
        table: t,
	}
}
func (m model) resetModel() model {
	r := initialModel()
	r.stateTable = m.stateTable
	r.viewWidth = m.viewWidth
	return r
}

func (m model) step() model { 
	if m.state <= 0 {
		return m
	}
	m.started = true
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

func (m model) Edit() model {
	var t textinput.Model
	in := make([]textinput.Model, 5)
	for i := range in {
		t = textinput.New()
		t.CharLimit = 1
		switch i {
		case 0:
			t.Placeholder = "New State (A-Z)"
		case 1:
			t.Placeholder = "New Symbol (0-9)"
		case 2:
			t.Placeholder = "Move (L/R)"
		}
		in[i] = t
	}
	//textinput example: https://github.com/charmbracelet/bubbletea/blob/master/examples/textinputs/main.go
	//m.stateTable[m.table.GetCursorY()][m.table.GetCursorX()].
	//	nextState = in[0].value[0]
	//m.stateTable[m.table.GetCursorY()][m.table.GetCursorX()].write = in[1]
	//m.stateTable[m.table.GetCursorY()][m.table.GetCursorX()].direction = (in[2] == 'R')
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n", "s":
			m = m.step()
		case "r":
			m = m.resetModel()
		case "e":
			m = m.Edit()
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
        case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)

		}
	}
    m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	s := "The current tape state:\n\n..."

	viewHead := m.head;

	for i:=0; i < m.viewWidth; i++ {
		if viewHead.Prev() == nil {
			m.tape.InsertBefore(0,viewHead)
		}
		viewHead = viewHead.Prev();
	}
	for i:=0; i < 2*m.viewWidth+1; i++ {
		s += fmt.Sprintf("%v ", viewHead.Value)
		if viewHead.Next() == nil {
			m.tape.InsertAfter(0,viewHead)
		}
		viewHead = viewHead.Next();
	}
	s = s[:len(s) - 1]
	s += "...\n"

		//s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	s += "^\n"
	s += fmt.Sprintf("Current state: %v\n", toRune(m.state))
    s += baseStyle.Render(m.table.View()) + "\n"
	s += "\nPress s to step, r to reset, +/- to see more/less of the tape, q to quit.\n"

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
	if i == ACCEPT {
		return "ACCEPT"
	} else if i == REJECT {
		return "REJECT"
	}
    return string(rune('A' + i))
}

type transState struct {
    nextState int
    write int
    direction bool // true -> move right
}

func TransStateToString(t transState) string {
	s := fmt.Sprintf("%v, %v, %v", toRune(t.nextState), t.write,
	func (dir bool) string { if dir { return "R" } 
		return "L" }(t.direction))
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
