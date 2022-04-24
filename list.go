package main

type listModel struct {
	viewModel
	count      int
	start      int
	end        int
	marked     int
	cursor     int
	scrollLock bool
}

type listView interface {
	view
	nextPage()
	prevPage()
	nextItem()
	prevItem()
	mark()
	setCursor(int)
	scrollToTop()
	scrollToBottom()
	findNext(string)
	findPrev(string)
	getCount() int
	getEnd() int
	getCursor() int
}

func (m *listModel) init(count int, scrollLock bool) {
	m.count = count
	m.marked = -1
	m.scrollLock = scrollLock
	if count > 0 {
		m.cursor = 0
	} else {
		m.cursor = -1
	}
}

func (m *listModel) setHeight(height int) {
	m.height = height
	m.updateLayout()
}

func (m *listModel) setSize(width, height int) {
	m.setWidth(width)
	m.setHeight(height)
}

func (m *listModel) setCount(count int) {
	m.count = count
	m.updateLayout()
}

func (m *listModel) updateLayout() {
	m.end = min(m.start+m.height, m.count)
	if m.end - m.start < m.height {
		m.start = max(m.end - m.height, 0)
	}

	if !m.scrollLock {
		m.cursor = min(m.cursor, m.count - 1)
		if m.cursor >= m.end {
			m.end = m.cursor + 1
			m.start = max(m.end - m.height, 0)
		} else if m.cursor < m.start {
			m.start = m.cursor
			m.end = min(m.start + m.height, m.count)
		}
	}
}

func (m *listModel) nextPage() {
	if m.scrollLock {
		m.scrollBy(m.height)
	} else {
		m.setCursor(min(m.cursor+m.height, m.count-1))
	}
}

func (m *listModel) prevPage() {
	if m.scrollLock {
		m.scrollBy(-m.height)
	} else {
		m.setCursor(max(m.cursor-m.height, 0))
	}
}

func (m *listModel) nextItem() {
	if m.scrollLock {
		m.scrollBy(1)
	} else if m.cursor < m.count-1 {
		m.setCursor(m.cursor + 1)
	}
}

func (m *listModel) prevItem() {
	if m.scrollLock {
		m.scrollBy(-1)
	} else if m.cursor > 0 {
		m.setCursor(m.cursor - 1)
	}
}

func (m *listModel) mark() {
	if m.marked == m.cursor {
		m.marked = -1
	} else {
		m.marked = m.cursor
	}
}

func (m *listModel) setCursor(index int) {
	if m.count == 0 {
		m.cursor = -1
	} else {
		m.cursor = index
		if m.cursor > m.end-1 {
			m.end = m.cursor + 1
			m.start = max(m.end-m.height, 0)
		} else if m.cursor < m.start {
			m.start = m.cursor
			m.end = min(m.start+m.height, m.count)
		}
	}
}

func (m *listModel) scrollToBottom() {
	if m.scrollLock {
		m.end = m.count
		m.start = max(m.end-m.height, 0)
	} else {
		m.setCursor(m.count - 1)
	}
}

func (m *listModel) scrollToTop() {
	if m.scrollLock {
		m.start = 0
		m.end = min(m.start+m.height, m.count)
	} else {
		m.setCursor(0)
	}
}

func (m *listModel) scrollBy(amount int) {
	if amount == 0 {
		return
	}
	if amount > 0 {
		m.end = min(m.end + amount, m.count)
		m.start = max(m.end - m.height, 0)
	} else {
		m.start = max(m.start + amount, 0)
		m.end = min(m.start + m.height, m.count)
	}
}

func (m listModel) getCount() int {
	return m.count
}

func (m listModel) getEnd() int {
	return m.end
}

func (m listModel) getCursor() int {
	if m.scrollLock {
		return m.end - 1
	}
	return m.cursor
}