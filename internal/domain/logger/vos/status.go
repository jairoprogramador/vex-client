package vos

import (
	"fmt"
)

type Status int

const (
	Running Status = iota
	Success
	Failure
	Warning
)

func (s Status) String() string {
	switch s {
	case Running:
		return "Running"
	case Success:
		return "Success"
	case Failure:
		return "Failure"
	case Warning:
		return "Warning"
	default:
		return "Unknown"
	}
}

func NewStatusFromString(status string) (Status, error) {
	switch status {
	case "Running":
		return Running, nil
	case "Success":
		return Success, nil
	case "Failure":
		return Failure, nil
	case "Warning":
		return Warning, nil
	default:
		return 0, fmt.Errorf("invalid status: %s", status)
	}
}
