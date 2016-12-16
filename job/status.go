package job

type Status uint8

const (
	Idle Status = iota
	Running
	Disabled
)

var status = [...]string{"Idle", "Running", "Disabled"}

func (s Status) String() string {
	return status[s]
}

// implementation of json.Marshaler interface
func (s Status) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

func AllStatus() []Status {
	return []Status{Idle, Running, Disabled}
}
