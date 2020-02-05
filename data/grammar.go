package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Grammar structure is used to read/write the json file format.
type Grammar struct {
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Extension   []string      `json:"extension,omitempty"`
	Grammar     []GrammarData `json:"grammar,omitempty"`
}

// GrammarData structure
type GrammarData struct {
	Option GrammarOption  `json:"option,omitempty"`
	Group  []GrammarGroup `json:"group,omitempty"`
}

// GrammarOption structure
type GrammarOption struct {
	GroupKeyword    string `json:"groupKeyword"`
	GroupByKeyword  bool   `json:"groupByKeyword"`
	AllowEmptyValue bool   `json:"allowEmptValue"`
}

// GrammarGroup structure
type GrammarGroup struct {
	Source  string         `json:"-"`
	Peek    []GrammarPeek  `json:"-"`
	Pattern string         `json:"pattern,omitempty"`
	Matches []GrammarMatch `json:"match,omitempty"`
}

// GrammarMatch structure
type GrammarMatch struct {
	Grammar Grammar `json:"grammar,omitempty"`
	Node    Node    `json:"node,omitempty"`
}

// GrammarPeek structure stores lines ahead and behind the current (source) line.
type GrammarPeek struct {
	Ahead  string `json:"-"`
	Behind string `json:"-"`
}

// CacheGrammar func
func CacheGrammar(name string) (grammar Grammar, err error) {
	file, err := os.Open(grammarFilePrefix + name + fileExtension)
	if err != nil {
		return grammar, err
	}
	read, err := ioutil.ReadAll(file)
	if err != nil {
		return grammar, err
	}
	err = json.Unmarshal(read, &grammar)
	if err != nil {
		return grammar, err
	}
	return grammar, nil
}

func (g *Grammar) setPeek()   {}
func (g *Grammar) getPeek()   {}
func (g *Grammar) setSource() {}
func (g *Grammar) mergeNodeValues() {
	// for _, group := range g.Group {
	// 	// group.Source
	// }
}
