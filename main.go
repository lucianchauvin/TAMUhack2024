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

	input []textinput.Model
	editMode bool
	focusIndex int
}


func initialModel() model {
	columns := []table.Column{
		{Title: "State/Symbol", Width: 12},
		{Title: "1", Width: 8},
		{Title: "2", Width: 8},
		{Title: "3", Width: 8},
	}

	rows := []table.Row{
		{"A", "", "", ""},
		{"B", "", "", ""},
		{"C", "", "", ""},
		{"D", "", "", ""},
		{"E", "", "", ""},
		{"F", "", "", ""},
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

	inputs := make([]textinput.Model, 3)
	
	var ti textinput.Model
	for i := range inputs {
		ti = textinput.New()
		ti.CharLimit = 1
		switch i {
		case 0:
			ti.Placeholder = "New State (A-Z)"
		case 1:
			ti.Placeholder = "New Symbol (0-9)"
		case 2:
			ti.Placeholder = "Move (L/R)"
		}
		inputs[i] = ti
	}	
	
	return model{
		tape: initTape,
		head: initTape.Front(),
		stateTable: [][]transState{{{1,1,true}, {1,0,false}}, {{1,1,false}, {0,0,true}}},
		state:		1,
		viewWidth:	10,
		started:	false,
		table: t,
		input: inputs,
		editMode: false,
		focusIndex: 0,
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
	// var t textinput.Model
	// in := make([]textinput.Model, 5)
	// for i := range in {
	// 	t = textinput.New()
	// 	t.CharLimit = 1
	// 	switch i {
	// 	case 0:
	// 		t.Placeholder = "New State (A-Z)"
	// 	case 1:
	// 		t.Placeholder = "New Symbol (0-9)"
	// 	case 2:
	// 		t.Placeholder = "Move (L/R)"
	// 	}
	// 	in[i] = t
	// }
	
	//textinput example: https://github.com/charmbracelet/bubbletea/blob/master/examples/textinputs/main.go
	//m.stateTable[m.table.GetCursorY()][m.table.GetCursorX()].
	//	nextState = in[0].value[0]
	//m.stateTable[m.table.GetCursorY()][m.table.GetCursorX()].write = in[1]
	//m.stateTable[m.table.GetCursorY()][m.table.GetCursorX()].direction = (in[2] == 'R')
	//return in

	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd	
	renderTable(m.stateTable, &m.table)
	m.table, cmd = m.table.Update(msg)
	m.table.UpdateViewport()
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.editMode {
			switch msg.String() {
			case "n", "s":
				m = m.step()
			case "r":
				m = m.resetModel()
			case "e":
				m.editMode = !m.editMode
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
		} else {
			switch msg.String() {
			case "e":
				m.editMode = !m.editMode
			case "tab":
				m.input[m.focusIndex].Blur()
				m.focusIndex++
				if m.focusIndex > 2 {
					m.focusIndex = 0
				}
				m.input[m.focusIndex].Focus()				
				
			case "-":
				if m.viewWidth > 2 {
					m.viewWidth--
				}
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				a := m.getCurrentTransState()
				if a != -1 {
					m.stateTable[m.table.GetCursorY()][m.table.GetCursorX() -1] = transState {a, a, false}
				} 
			}
			cmd = m.updateInputs(msg)
		}
		
	}
	return m, cmd
}

func (m model) getCurrentTransState() int {
	var ns int
	nsc := int(m.input[0].Value()[0])
	if (nsc > 97 && nsc < 122) {
		ns = nsc - 97 
	} else if (nsc > 65 && nsc < 90) {
		ns = nsc - 65 
        } else {
		ns = int(nsc-int('A')+1)		
	}
	return ns
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.input))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.input {
		m.input[i], cmds[i] = m.input[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var s string
	if !m.editMode {
		s = "The current tape state:\n\n..."

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
	} else {

		for i := range m.input {
			s += m.input[i].View()
			if i < len(m.input)-1 {
				s += "\n"
			}
		}
	}

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
		return "✓"
	} else if i == REJECT {
		return "⨯"
	}
	return string(rune('A' + i))
}

type transState struct {
	nextState int
	write int
	direction bool // true -> move right
}

func (t transState) render() string {
	s := fmt.Sprintf("%v, %v, %v", toRune(t.nextState), t.write,
		func (dir bool) string { if dir { return "R" } 
			return "L" }(t.direction))
	return s	
}

func renderTable(ts [][]transState, table *table.Model) {
	for y := 0; y < len(ts); y++ {
		for x := 0; x < len(ts[y]); x++ {
			table.Rows()[y][x+1] = ts[y][x].render()
		}
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
