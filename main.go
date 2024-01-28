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
	//"github.com/muesli/termenv"
)

const ACCEPT = 24 //Y
const REJECT = 13 //N

var tableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240")).
	AlignVertical(lipgloss.Center)

type model struct {
	tape		*list.List
	head		*list.Element
	stateTable	[][]transState 
	table		table.Model
	state		int // A = 0
	viewWidth	int //number of steps run
	started		bool // you can't edit the table if you've started!

	input []textinput.Model
	editMode bool
	editFailed bool
	focusIndex int
}

//run messages on startup
func (m model) Init() tea.Cmd {

	return tea.Batch(
		tea.SetWindowTitle("T-soding"), 
		tea.EnterAltScreen,
	 )
}

func initialModel() model {
	columns := []table.Column{
		{Title: "State/Symbol", Width: 12},
		{Title: "0", Width: 8},
		{Title: "1", Width: 8},
		{Title: "2", Width: 8},
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
		table.WithHeight(6),
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
	initTape.PushBack(2)
	initTape.PushBack(2)
	initTape.PushBack(2)
	initTape.PushBack(1)

	inputs := make([]textinput.Model, 3)
	
	var ti textinput.Model
	for i := range inputs {
		ti = textinput.New()
		ti.CharLimit = 1
		switch i {
		case 0:
			ti.Prompt = "New State   (A-Z): "
			ti.Placeholder = "A"
		case 1:
			ti.Prompt = "New Symbol  (0-9): "
			ti.Placeholder = "0"
		case 2:
			ti.Prompt = "Direction   (R/L): "
			ti.Placeholder = "L"
		}
		ti.PromptStyle.AlignHorizontal(lipgloss.Left)
		ti.TextStyle.AlignHorizontal(lipgloss.Left)
		ti.PlaceholderStyle.AlignHorizontal(lipgloss.Left)
		inputs[i] = ti
	}	
	A := ACCEPT
	R := REJECT
	
	return model{
		tape: initTape,
		head: initTape.Front(),
		//rows are states, columns are symbols
		stateTable: [][]transState{
			{{A,0,true }, {4,0,false},{R,2,false}},
			{{A,0,false}, {R,1,false},{5,0,false}},
			{{0,0,false}, {2,1,true },{2,2,true }},
			{{1,0,false}, {3,1,true },{3,2,true }},
			{{5,0,true }, {4,1,false},{4,2,false}},
			{{A,0,true }, {2,0,true },{3,0,true }},
		},
		state:		4,
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
	renderTable(r.stateTable, &r.table)
	r.table.UpdateViewport()
	return r
}

func (m model) step() model { 
	if m.state == ACCEPT || m.state == REJECT {
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
	m.table.SetCursor(curState, curValue+1)
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd	
	
	renderTable(m.stateTable, &m.table)
	m.table.UpdateViewport()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.editMode {
			// we only want the table to take inputs when not in edit mode
			m.table, cmd = m.table.Update(msg)
			switch msg.String() {
			case "n", "s":
				m = m.step()
			case "r":
				m = m.resetModel()
			case "e":
				if(!m.started) { //don't edit if we've already started
					m.editMode = true
					m = m.ResetInputs()
					m.input[m.focusIndex].Blur()				
					m.focusIndex = 0
					m.input[m.focusIndex].Focus()				
				}
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
		//	case "enter":
		//		return m, tea.Batch(
		//			tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
		//		)

			}
		} else {
			switch msg.String() {
			case "e":
				m = m.ResetInputs()
				m.editMode = false
			case "down", "tab":
				m.input[m.focusIndex].Blur()
				m.focusIndex++
				if m.focusIndex > 2 {
					m.focusIndex = 0
				}
				m.input[m.focusIndex].Focus()				

			case "up":
				m.input[m.focusIndex].Blur()
				m.focusIndex--
				if m.focusIndex < 0 {
					m.focusIndex = 2
				}
				m.input[m.focusIndex].Focus()				
				
			case "-":
				if m.viewWidth > 2 {
					m.viewWidth--
				}
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				InTrans, ok := m.getInputTrans()
				if ok {
					m.editFailed = false
					m.stateTable[m.table.GetCursorY()][m.table.GetCursorX() -1] = InTrans
					m.editMode = false
					m = m.ResetInputs()
				} else {
					m.editFailed = true
				}
			}
			cmd = m.updateInputs(msg)
		}
		
	}
	return m, cmd
}

func (m model) ResetInputs() model {
		for i := range m.input {
			m.input[i].Reset()
		}
		return m
}
func (m model) getInputTrans() (transState, bool) {
	for i := range m.input {
		if len(m.input[i].Value()) == 0 {
			return transState{0,0,false}, false
		}
	}
	t := transState{int(m.input[0].Value()[0] - 'A'), 
			int(m.input[1].Value()[0] - '0'),
			m.input[2].Value() == "R"} 
	isdir := m.input[2].Value() == "R" || m.input[2].Value() == "L" 

	return t, t.nextState >= 0 && t.nextState <= 26 &&
		t.write >= 0 && t.write < 10 &&
		isdir
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
		//s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		s += "^\n"
		s += fmt.Sprintf("Current state: %v\n", toRune(m.state))
	} else {

		for i := range m.input {
			s += m.input[i].View()
			if i < len(m.input)-1 {
				s += "\n"
			}
		}
		s += "\n"
	}

	s += "\n" + tableStyle.Render(m.table.View()) + "\n"
	if m.editFailed { s += "\nthat's not a valid table entry!\n"
	} else if m.editMode {
		s += "\n↑/↓: change field • e: stop editing • ↵: save to table • q: quit.\nUse 'Y' as the yes/accept state and 'N' as the no/fail state."
	} else if !m.started {
		s += "\ns: step • e: edit  • +/-: view more/less tape • q: quit.\n"
	} else {
	s += "\ns: step • r: reset • +/-: view more/less tape • q: quit.\n"
	}

	return style.
		PaddingTop(5).
		PaddingLeft(20).
		Render(s)
}

var style = lipgloss.NewStyle().
	//Foreground(lipgloss.Color("#FAFAFA")).
	//Background(lipgloss.Color("#7D56F4")).
	Align(lipgloss.Center)
var InputStyle = lipgloss.NewStyle().Align(lipgloss.Left)


func toRune(i int) string {
	if i == ACCEPT {
		return "✓"
	} else if i == REJECT {
		return "⨯"
	}
	return string(rune('A' + i))
}

type TestCase struct {
	tape *list.List
	head *list.Element
	pass bool
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
