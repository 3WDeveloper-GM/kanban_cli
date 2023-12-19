package status

type Status int

const (
	Todo Status = iota
	InProgress
	Done
)

func (s Status) String() string {
	switch s {
	case Todo:
		return "todo"
	case InProgress:
		return "in Progress"
	case Done:
		return "done"
	default:
		return ""
	}
}

func (s *Status) Next() {
	if *s == Done {
		*s = Todo
	} else {
		*s++
	}
}

func (s *Status) Prev() {
	if *s == Todo {
		*s = Done
	} else {
		*s--
	}
}

func (s Status) GetNext() Status {
	if s == Done {
		return Todo
	}
	return s + 1
}

func (s Status) GetPrev() Status {
	if s == Todo {
		return Done
	}
	return s - 1
}
