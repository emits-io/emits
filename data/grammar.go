package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// GrammarFile structure is used to read/write the json file format.
type GrammarFile struct {
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Extension   []string  `json:"extension,omitempty"`
	Grammar     []Grammar `json:"grammar,omitempty"`
}

// Grammar structure
type Grammar struct {
	Pattern string           `json:"pattern,omitempty"`
	Match   map[string]Match `json:"match,omitempty"`
}

// Match struct
type Match struct {
	Grammar []Grammar         `json:"grammar,omitempty"`
	Set     map[string]string `json:"set,omitempty"`
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

func setPatternResults(g Grammar, pattern string, source string) map[string]string {
	r := regexp.MustCompile(pattern)
	m := r.MatchString(source)
	if m {
		data := r.FindStringSubmatch(source)
		name := r.SubexpNames()
		result := make(map[string]string)
		for match := range data {
			result[strings.ToLower(name[match])] = data[match]
		}
		return result
	}
	return nil
}

func setPatternMatch(g Grammar, node *Node, source string, pattern string, result map[string]string) *Node {
	for key, val := range g.Match {
		if len(val.Set) > 0 {
			for set, value := range val.Set {
				r := regexp.MustCompile("{{(\\w+)}}")
				matches := r.FindAllStringSubmatch(value, -1)
				for _, match := range matches {
					value = strings.ReplaceAll(value, match[0], result[match[1]])
				}
				switch set {
				case "flags":
					node.Flags = strings.Split(value, ",")
				case "index":
					index, err := strconv.ParseInt(value, 10, 0)
					if err == nil {
						node.Index = int(index)
					}
				case "keyword":
					node.Keyword = value
				case "line":
					line, err := strconv.ParseInt(value, 10, 0)
					if err == nil {
						node.Line = int(line)
					} else {
						node.Line = 0
					}
				case "parent":
					parent, err := strconv.ParseInt(value, 10, 0)
					if err == nil {
						node.Parent = int(parent)
					} else {
						node.Parent = 0
					}
				case "separator":
					node.Separator = strings.ToLower(value) == "true"
				case "value":
					node.Value = value
				}
			}
		} else if len(result[key]) > 0 {
			for _, grammar := range val.Grammar {
				results := setPatternResults(grammar, grammar.Pattern, result[key])
				if results != nil {
					setPatternMatch(grammar, node, result[key], grammar.Pattern, results)
				}
			}
		}
	}
	return node
}

func (g *GrammarFile) process(n Node) Node {
	for _, grammar := range g.Grammar {
		setPatternMatch(grammar, &n, n.Value, grammar.Pattern, setPatternResults(grammar, grammar.Pattern, n.Value))
	}
	return n
}
