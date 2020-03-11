package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
)

// GrammarFile structure is used to read/write the json file format.
type GrammarFile struct {
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Extension   []string  `json:"extension,omitempty"`
	Grammar     []Grammar `json:"grammar,omitempty"`
}

// CacheGrammarFile func
func CacheGrammarFile(name string) (grammar GrammarFile, err error) {
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

// hasExtension func
func (g *GrammarFile) hasExtension(extension string) bool {
	for _, e := range g.Extension {
		if e == extension {
			return true
		}
	}
	return false
}

// Grammar structure
type Grammar struct {
	Pattern string           `json:"pattern,omitempty"`
	Match   map[string]Match `json:"match,omitempty"`
}

func setPatternResults(g Grammar, pattern string, source string) map[string]string {
	r := regexp.MustCompile(pattern)
	data := r.FindStringSubmatch(source)
	name := r.SubexpNames()
	result := make(map[string]string)
	for match := range data {
		result[name[match]] = data[match]
	}
	return result
}

func setPatternMatch(g Grammar, node *Node, source string, pattern string, result map[string]string) *Node {
	for key, val := range g.Match {
		if len(val.Set) > 0 {
			switch val.Set {
			case "value":
				node.Value = result[key]
			case "keyword":
				node.Keyword = result[key]
			}
		} else if len(result[key]) > 0 {
			setPatternMatch(val.Grammar, node, result[key], val.Grammar.Pattern, setPatternResults(val.Grammar, val.Grammar.Pattern, result[key]))
		}
	}
	return node
}

// Match struct
type Match struct {
	Grammar Grammar `json:"grammar,omitempty"`
	Set     string  `json:"set,omitempty"`
}

func (g *GrammarFile) process(n Node) Node {
	for _, grammar := range g.Grammar {
		setPatternMatch(grammar, &n, n.Value, grammar.Pattern, setPatternResults(grammar, grammar.Pattern, n.Value))
	}
	return n
}

func (g *Grammar) setPeek()         {}
func (g *Grammar) getPeek()         {}
func (g *Grammar) setSource()       {}
func (g *Grammar) mergeNodeValues() {}
