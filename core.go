package main

type Datum struct {
	Value float64 `json:"y"`
	T     int64   `json:"x"`
}

type StatRecord struct {
	Name      string
	IsCounter bool
	Values    []Datum

	Capacity int
	w        int
}

type Core struct {
	Stats map[string]*StatRecord
}

type StatUpdate struct {
	Key       string
	Value     float64
	Timestamp int64
}

func NewCore() *Core {
	return &Core{make(map[string]*StatRecord)}
}

func (c *Core) Loop(updateChan chan StatUpdate, funcChan chan func(c *Core)) {
	for {
		select {
		case up, ok := <-updateChan:
			if !ok {
				return
			}
			c.Update(up)
		case f, ok := <-funcChan:
			if !ok {
				return
			}
			f(c)
		}

	}
}

func (c *Core) Update(up StatUpdate) {

	rec, ok := c.Stats[up.Key]
	if !ok {
		rec = &StatRecord{}
		c.Stats[up.Key] = rec
	}
	rec.AppendValue(up.Value, up.Timestamp)
}

func (sr *StatRecord) AppendValue(f float64, t int64) {
	if sr.Capacity == 0 {
		sr.Capacity = 1000
	}

	if len(sr.Values) < sr.Capacity {
		sr.w++
		if len(sr.Values) == 0 {
			sr.Values = append(sr.Values, Datum{f, t})
			return
		}
		at := len(sr.Values) - 1
		if sr.Values[at].T == t {
			if sr.IsCounter {
				sr.Values[at].Value += f
			} else {
				sr.Values[at].Value = f
			}
		} else {
			sr.Values = append(sr.Values, Datum{f, t})
		}
	} else {
		if sr.Values[sr.w].T == t {
			if sr.IsCounter {
				sr.Values[sr.w].Value += f
			} else {
				sr.Values[sr.w].Value = f
			}
		} else {
			sr.Values[sr.w].Value = f
			sr.Values[sr.w].T = t
			sr.w++
			sr.w = sr.w % len(sr.Values)
		}
	}
}