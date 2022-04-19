package main

type listModel struct {
	count  int
	first  int
	last   int
	height int
	cursor int
	marked int
}

func newListModel(count int) listModel {
	return listModel{
		count:  count,
		marked: -1,
		cursor: 0,
	}
}

func newCursorlessListModel(count int) listModel {
	return listModel{
		count:  count,
		marked: -1,
		cursor: -1,
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
	}

	m.first += m.height
	m.last += m.height
	if m.last >= m.count {
		m.last = m.count - 1
		m.first = m.last - m.height + 1
	}

	if m.cursor != -1 {
		if m.cursor > m.last {
			m.cursor = m.last
		}
	}
}

func (m *listModel) prevPage() {
	if m.cursor != -1 {
		m.cursor -= m.height
	}

	m.first -= m.height
	m.last -= m.height
	if m.first < 0 {
		m.first = 0
		m.last = m.first + m.height - 1
	}

	if m.cursor != -1 {
		if m.cursor < 0 {
			m.cursor = 0
		}
	}
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
