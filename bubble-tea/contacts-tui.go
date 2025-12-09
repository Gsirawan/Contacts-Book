package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Screen/State
type screen int

const (
	menuScreen   screen = iota // Main menu
	listScreen                 // List contacts
	addScreen                  // Add contact
	searchScreen               // Search
	deleteScreen               // Delete
	editScreen                 // Edit
	filePath     = "contacts.txt"
)

type Contact struct {
	Name   string
	Email  string
	Mobile string
}

type model struct {
	currentScreen screen
	cursor        int
	contacts      []Contact
	table         table.Model
	focusIndex    int
	inputs        []textinput.Model
}

func initialModel() model {
	return model{
		currentScreen: menuScreen,
		cursor:        0,
		contacts:      []Contact{},
		inputs:        initialInputs(),
		focusIndex:    0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// table handling
	if m.currentScreen == listScreen {
		m.table, cmd = m.table.Update(msg)
	}
	if m.currentScreen == addScreen {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				// Go back to menu
				m.currentScreen = menuScreen
				m.cursor = 0
				return m, nil

			case "tab", "down":
				m.focusIndex++
				if m.focusIndex > len(m.inputs)-1 {
					m.focusIndex = 0
				}

			case "shift+tab", "up":
				m.focusIndex--
				if m.focusIndex < 0 {
					m.focusIndex = len(m.inputs) - 1
				}

			case "ctrl+s":
				// save contacts
				name := m.inputs[0].Value()
				email := m.inputs[1].Value()
				mobile := m.inputs[2].Value()
				if name == "" || email == "" || mobile == "" {
					return m, nil
				}

				// Save to file
				saveContact(Contact{
					Name:   name,
					Email:  email,
					Mobile: mobile,
				})

				// Reset inputs and go back to menu
				m.inputs = initialInputs()
				m.currentScreen = menuScreen
				m.cursor = 0
				return m, nil
			}
			// Update which input has focus
			for i := range m.inputs {
				if i == m.focusIndex {
					m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
		}
		m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
		return m, cmd
	}
	// main key handling
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
			if m.cursor < 5 { // 6 menu items (0-5)
				m.cursor++
			}
		case "esc":
			if m.currentScreen != menuScreen {
				m.currentScreen = menuScreen
				m.cursor = 0
			}
		case "enter":
			// Handle menu selection (we'll add this next)
			switch m.cursor {
			case 0: // List Contacts
				m.contacts = loadContacts()
				m.table = makeContactTable(m.contacts)
				m.currentScreen = listScreen
			case 1: // add contacts
				m.inputs = initialInputs()
				m.focusIndex = 0
				m.currentScreen = addScreen

			case 5: // Exit
				return m, tea.Quit
			}
		}
	}

	return m, cmd
}

func makeContactTable(contacts []Contact) table.Model {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Email", Width: 30},
		{Title: "Mobile", Width: 20},
	}

	rows := []table.Row{}
	for _, c := range contacts {
		rows = append(rows, table.Row{c.Name, c.Email, c.Mobile})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
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
	t.SetStyles(s)

	return t
}

func saveContact(contact Contact) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "%s,%s,%s\n", contact.Name, contact.Email, contact.Mobile)
}

func loadContacts() []Contact {
	var contacts []Contact
	var err error
	var file *os.File

	file, err = os.OpenFile(filePath, os.O_RDONLY, 0o644)
	if err != nil {
		log.Fatalf("Error Opening the file %v\n", err)
		return contacts
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatalf("Error closing file %v\n:", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) == 3 {
			contacts = append(contacts, Contact{
				Name:   parts[0],
				Email:  parts[1],
				Mobile: parts[2],
			})
		}
	}

	return contacts
}

func initialInputs() []textinput.Model {
	inputs := make([]textinput.Model, 3)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Name"
	inputs[0].Focus()
	inputs[0].CharLimit = 50
	inputs[0].Width = 30

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Email"
	inputs[1].CharLimit = 50
	inputs[1].Width = 30

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Mobile"
	inputs[2].CharLimit = 15
	inputs[2].Width = 30

	return inputs
}

func (m model) View() string {
	if m.currentScreen == menuScreen {
		s := "Contact Manager\n\n"

		choices := []string{"List Contacts", "Add Contact", "Search", "Delete", "Edit", "Exit"}

		for i, choice := range choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		s += "\nUse arrow keys to navigate, Enter to select, q to quit\n"
		return s
	}
	if m.currentScreen == listScreen {
		s := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 30).
			Render("Contact List")
		return s + "\n\n" + m.table.View() + "\n\nPress ESC to go back, and 'q' to quit\n"
	}
	if m.currentScreen == addScreen {
		s := "Add New Contact\n\n"

		s += "Name:\n"
		s += m.inputs[0].View() + "\n\n"

		s += "Email:\n"
		s += m.inputs[1].View() + "\n\n"

		s += "Mobile:\n"
		s += m.inputs[2].View() + "\n\n"
		s += "Tab to move between fields, Ctrl+S to save, ESC to cancel\n"
		s += "ESC to cancel\n"

		return s
	}
	return "Other screen (TODO)"
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
