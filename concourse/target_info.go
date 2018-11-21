package concourse

import "time"

type TargetInfo struct {
	Name   string
	URL    string
	Team   string
	Expiry time.Time
}
