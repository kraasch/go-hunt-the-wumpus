package main

import (
	"fmt"
	"math/rand"
	"time"
	"os"
  _ "strconv" // for debug prints.
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
  top_msg = "HUNT THE WUMPUS\n Quit (q), move (hjkl), shoot (HJKL)."
  bot_msg = "..."
  you_died = false
  game_over = false
  arrow_count = 0
  outerBox = lipgloss.NewStyle().
    BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color("56"))
  cursorStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4"))
  fogStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FFFFFF")).
    Background(lipgloss.Color("#FFFFFF"))
  noFogStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#000000")).
    Background(lipgloss.Color("#000000"))
)

type model struct {
	width    int
	height   int
  arr      [5][5]string
  cursor_x int
  cursor_y int
}

func min(a, b int) int {
  if a < b {
    return a
  }
  return b
}

func max(a, b int) int {
  if a > b {
    return a
  }
  return b
}

func (m model) Init() tea.Cmd { return nil }

func update_positional_messages(m *model) {
  bot_msg = ""
  ////////////////////////////////////////////////
  c := m.arr[m.cursor_y][m.cursor_x]
  if c == "a" { // Stepped on an arrow.
    // 2. a = arrow
    m.arr[m.cursor_y][m.cursor_x] = "N" // Remove the arrow.
    arrow_count += 1
    // bot_msg += "You found an arrow! (" + strconv.Itoa(arrow_count) + ") " // NOTE: debug print.
    bot_msg += "You found an arrow! " // TODO: uncomment.
  } else if c == "X" { // Stepped into fog.
    m.arr[m.cursor_y][m.cursor_x] = "N" // Remove the fog.
  } else if c == "b" || c == "o" || c == "w" {
    // 3. b = bat, 4. o = hole, 5. w = wumpus
    bot_msg += "You Died! "
    switch c {
    case "b":
      bot_msg += "The bat eats you. "
    case "o":
      bot_msg += "You fall into a hole. "
    case "w":
      bot_msg += "The wumpus got you. "
    }
    you_died = true
    return
  } else {
    // 1. h = hunter
    {} // golang's no-op.
  }
  ////////////////////////////////////////////////
  // Loop over the 3x3 subarray around hunter's position at (cursor_x, cursor_y).
  // Make sure that (I, J) stays within bounds of the 5x5 array.
  for I := max(0, m.cursor_y-1); I < m.cursor_y+2 && I < 5 && I >= 0; I++ {
    for J := max(0, m.cursor_x-1); J < m.cursor_x+2 && J < 5 && J >= 0; J++ {
      switch m.arr[I][J] {
      case "b":
        bot_msg += "You hear flapping. "
      case "o":
        bot_msg += "You feel a draft. "
      case "w":
        bot_msg += "You smell wumpus. "
      }
    }
  }
  ////////////////////////////////////////////////
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd
  switch msg := msg.(type) {
  case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
  case tea.KeyMsg:
    switch msg.String() {
      case "H", "a": // shoot left.
      if you_died || game_over {
        return m, tea.Quit
      }
      if arrow_count > 0 {
        bot_msg = "You shoot left! "
        arrow_count--
        left := max(m.cursor_x - 1, 0)
        if m.arr[m.cursor_y][left] == "w" {
          bot_msg += "You shoot the wumpus!"
          game_over = true
        }
      }
      case "L", "d": // shoot right.
      if you_died || game_over {
        return m, tea.Quit
      }
      if arrow_count > 0 {
        bot_msg = "You shoot right! "
        arrow_count--
        right := min(m.cursor_x + 1, 4)
        if m.arr[m.cursor_y][right] == "w" {
          bot_msg += "You shoot the wumpus!"
          game_over = true
        }
      }
      case "K", "w": // shoot up.
      if you_died || game_over {
        return m, tea.Quit
      }
      if arrow_count > 0 {
        bot_msg = "You shoot up! "
        arrow_count--
        up := max(m.cursor_y - 1, 0)
        if m.arr[up][m.cursor_x] == "w" {
          bot_msg += "You shoot the wumpus!"
          game_over = true
        }
      }
      case "J", "s": // shoot down.
      if you_died || game_over {
        return m, tea.Quit
      }
      if arrow_count > 0 {
        bot_msg = "You shoot down! "
        arrow_count--
        down := min(m.cursor_y + 1, 4)
        if m.arr[down][m.cursor_x] == "w" {
          bot_msg += "You shoot the wumpus!"
          game_over = true
        }
      }
      case "h", "left": // move left.
      if you_died || game_over {
        return m, tea.Quit
      }
      if m.cursor_x > 0 {
        m.cursor_x--
        update_positional_messages(&m)
      }
      case "l", "right": // move right.
      if you_died || game_over {
        return m, tea.Quit
      }
      if m.cursor_x < 4 {
        m.cursor_x++
        update_positional_messages(&m)
      }
      case "k", "up": // move up.
      if you_died || game_over {
        return m, tea.Quit
      }
      if m.cursor_y > 0 {
        m.cursor_y--
        update_positional_messages(&m)
      }
      case "j", "down": // move down.
      if you_died || game_over {
        return m, tea.Quit
      }
      if m.cursor_y < 4 {
        m.cursor_y++
        update_positional_messages(&m)
      }
    case "q":
      return m, tea.Quit
    }
  }
  return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return ""
	}
  r   := top_msg + "\n"
  r   += outerBox.Render(pack(m.arr, m))
  r   += "\n" + bot_msg + "\n"
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, r)
}

func pack(in [5][5]string, m model) string {
  s := ""
  for i := 0; i < len(in); i++ {
    s += " "
    for j := 0; j < len(in[i]); j++ {
      // character := in[i][j] // NOTE: for debugging.
      if m.cursor_x == j && m.cursor_y == i {
        // Render the hunter.
        // cursor := cursorStyle.Render(character) // NOTE: for debugging.
        cursor := cursorStyle.Render(" ")
        s += fmt.Sprintf("%s ", cursor)
      } else if in[i][j] == "N" {
        // Render empty tile.
        // cursor := noFogStyle.Render(character) // NOTE: for debugging.
        cursor := noFogStyle.Render(" ")
        s += fmt.Sprintf("%s ", cursor)
      } else if in[i][j] == "X" {
        // Render foggy tile.
        // cursor := fogStyle.Render(character) // NOTE: for debugging.
        cursor := fogStyle.Render(" ")
        s += fmt.Sprintf("%s ", cursor)
      } else {
        // Render object.
        // s += fmt.Sprintf("%s ", character) // NOTE: for debugging.
        {} // golang's no-op.
        cursor := fogStyle.Render(" ")
        s += fmt.Sprintf("%s ", cursor)
      }
    }
    if i != len(in) -1 { // All but last line append newline.
      s += fmt.Sprintln()
    }
  }
  return s
}

func main() {
  //////////////////////////////////////////////////////////
  rand.Seed(time.Now().UnixNano()) // Seed the random number generator to get different results each time.
  selectedIndices := make(map[int]bool) // Create a slice to hold the selected indices.
  for len(selectedIndices) < 5 { // Loop until we have 5 unique indices.
    index := rand.Intn(25) // Generate a random index between 0 and 24.
    selectedIndices[index] = true // Add the index to the map (map keys are unique, so no duplicates).
  }
  var indices []int // Convert the map keys to a slice for easy access.
  for index := range selectedIndices {
    indices = append(indices, index)
  }
  // fmt.Printf("Random indices: %d\n", indices) // NOTE: debug print.
  //////////////////////////////////////////////////////////
  // 1. h = hunter
  // 2. a = arrow
  // 3. b = bat
  // 4. o = hole
  // 5. w = wumpus
  letters := []rune{'h', 'a', 'b', 'o', 'w'} // The slice to shuffle.
  rand.Shuffle(len(letters), func(i, j int) { // Shuffle the slice in place.
    letters[i], letters[j] = letters[j], letters[i]
  })
  // fmt.Println("Shuffled letters:", string(letters)) // NOTE: debug print.
  //////////////////////////////////////////////////////////
  var newArr [5][5]string
  m := model{0, 0, newArr, 0, 0}
  for i := 0; i < len(m.arr); i++ { // Filling the array with values.
    for j := 0; j < len(m.arr[i]); j++ {
      // m.arr[i][j] = fmt.Sprintf("%c", 'A' + i*5 + j) // NOTE: for debug.
      m.arr[i][j] = "X"
    }
  }
  //////////////////////////////////////////////////////////
  // MERGE.
  for i, I := range indices {
    m.arr[I%5][I/5] = string(letters[i])
    // set inital hunter position.
    if letters[i] == 'h' {
      m.cursor_x = I/5
      m.cursor_y = I%5
      m.arr[I%5][I/5] = "N" // Uncover the start location.
    }
  }
  //////////////////////////////////////////////////////////

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
