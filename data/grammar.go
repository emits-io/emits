package data

// grammar structure is used to read/write the json file format.
type grammar struct {
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Extension   []string  `json:"extension,omitempty"`
	Grammar     []Grammar `json:"grammar,omitempty"`
}

// Grammar structure
type Grammar struct {
	Option Option  `json:"option,omitempty"`
	Group  []Group `json:"group,omitempty"`
}

// Option structure
type Option struct {
	GroupKeyword    string `json:"groupKeyword"`
	GroupByKeyword  bool   `json:"groupByKeyword"`
	AllowEmptyValue bool   `json:"allowEmptValue"`
}

// Group structure
type Group struct {
	Source  string  `json:"-"`
	Peek    []Peek  `json:"-"`
	Pattern string  `json:"pattern,omitempty"`
	Matches []Match `json:"match,omitempty"`
}

// Match structure
type Match struct {
	Grammar Grammar `json:"grammar,omitempty"`
	Node    Node    `json:"node,omitempty"`
}

// Peek structure stores lines ahead and behind the current (source) line.
type Peek struct {
	Ahead  string `json:"-"`
	Behind string `json:"-"`
}

func (g *Grammar) setPeek()   {}
func (g *Grammar) getPeek()   {}
func (g *Grammar) setSource() {}
func (g *Grammar) mergeNodeValues() {
	// for _, group := range g.Group {
	// 	// group.Source
	// }
}
