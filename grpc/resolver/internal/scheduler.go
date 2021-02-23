package internal

// Scheduler scheduler.
type Scheduler struct {
	Clients map[string]*ZoneStrategy `json:"clients"`
}

// ZoneStrategy is the scheduling strategy of all zones
type ZoneStrategy struct {
	Zones map[string]*Strategy `json:"zones"`
}

// Strategy is zone scheduling strategy.
type Strategy struct {
	Weight int64 `json:"weight"`
}

