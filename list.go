package main

type listModel struct {
	viewModel
	count  int
	first  int
	last   int
	cursor int
	marked int
}

type listView interface {
	nextPage()
	prevPage()
	nextItem()
	prevItem()
	cursor()
	setCursor(int)
	marker()
	setMarker(int)
}

func (m *listModel) init(count int, cursor bool) {
	m.count = count
	m.marked = -1
	if cursor {
		m.cursor = 0
	} else {
		m.cursor = -1
	}
}

func (m *listModel) setHeight(height int) {
	if height > m.count {
		m.height = m.count + 1
	} else {
		m.height = height
	}
	m.last = m.first + m.height - 1
	if m.cursor > m.last-1 {
		m.last = m.cursor + 1
		m.first = m.last - m.height + 1
	}
}

func (m *listModel) setSize(width, height int) {
	m.setWidth(width)
	m.setHeight(height)
}

func (m *listModel) setCount(count int) {
	m.count = count
	if count < m.height {
		m.height = m.count
	}
	m.last = m.first + m.height - 1
	if m.last > m.count-1 {
		m.last = m.count - 1
	}
	if m.cursor > m.last-1 {
		m.last = m.cursor + 1
		m.first = m.last - m.height + 1
	}
}

func (m *listModel) nextPage() {
	if m.cursor != -1 {
		m.cursor += m.height
		if m.cursor >= m.count {
			m.cursor = m.count - 1
		}
	}

	if m.last == m.count {
		return
	}

	delta := m.height
	if m.last+m.height >= m.count {
		delta = m.count - m.last
	}

	m.last += delta
	m.first += delta
}

func (m *listModel) prevPage() {
	if m.cursor != -1 {
		m.cursor -= m.height
		if m.cursor < 0 {
			m.cursor = 0
		}
	}

	if m.first == 0 {
		return
	}

	delta := m.height
	if m.first-m.height < 0 {
		delta = m.first
	}

	m.first -= delta
	m.last -= delta
}

func (m *listModel) nextItem() {
	if m.cursor == -1 {
		// Not using cursor
		if m.last < m.count-1 {
			m.first += 1
			m.last += 1
		}
	} else if m.cursor < m.count-1 {
		m.cursor += 1
		if m.cursor > m.last-1 {
			m.first += 1
			m.last += 1
		}
	}
}

func (m *listModel) prevItem() {
	if m.cursor == -1 {
		// Not using cursor
		if m.first > 0 {
			m.first -= 1
			m.last -= 1
		}
	} else if m.cursor > 0 {
		m.cursor -= 1
		if m.cursor < m.first {
			m.first -= 1
			m.last -= 1
		}
	}
}

func (m *listModel) mark() {
	if m.marked == m.cursor {
		m.marked = -1
	} else {
		m.marked = m.cursor
	}
}