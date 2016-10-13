package clix

import "github.com/gdamore/tcell"

// for drawing multiple menus (using m.AddSibling)
func (m *MenuBar) drawnextto(parent *MenuBar, sibnum int) {
	xmax, ymax := m.screen.Size()
	var ts = 0
	if parent.title != "" {
		ts = 1
	}
	if parent.title != "" || m.title != "" {
		Type(m.screen, 20*(sibnum+1), ymax-parent.mostitems-2, tcell.StyleDefault, m.title)
		ts++
	}
	var itemnum int
	for i, v := range m.Children {
		itemnum++
		runelabel := []rune(v.Label)
		for r, x, y := 0, 20*(sibnum+1), ymax-parent.mostitems-ts+itemnum+1; r < len(runelabel); r++ {
			if r >= len(runelabel) {
				break
			}
			x++
			if x > xmax {
				y++
				x = 20
			}
			if y > ymax {
				break
			}

			if i == m.Selection && parent.zindex-1 == sibnum {
				m.screen.SetCell(x, y, 2, rune(runelabel[r]))

			} else {
				m.screen.SetCell(x, y, tcell.StyleDefault, rune(runelabel[r]))

			}

		}

	}

}
