package tracker

type PathStats struct {
	Methods *ConcurrentCounter[string]
	Status  *ConcurrentCounter[int]
	Total   *SingleCounter
}

func newPathStats() *PathStats {
	return &PathStats{
		Methods: NewConcurrentCounter[string](),
		Status:  NewConcurrentCounter[int](),
		Total:   &SingleCounter{},
	}
}

func (p *PathStats) Snapshot() map[string]interface{} {
	return map[string]interface{}{
		"Methods": p.Methods.Snapspot(),
		"Status":  p.Status.Snapspot(),
		"Total":   p.Total.Get(),
	}
}

// func (p *PathStats) ToJson() map[string]interface{} {
// 	return map[string]interface{}{
// 		"Methods": p.Methods.Snapspot(),
// 		"Status":  p.Status.Snapspot(),
// 		"Total":   p.Total.Get(),
// 	}
// }
