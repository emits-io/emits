package data

import "fmt"

// Cache struct
type Cache struct {
	Grammar []Grammar
}

// HasGrammar func
func (c *Cache) HasGrammar() bool {
	return len(c.Grammar) > 0
}

// ProcessGrammar func
func (c *Cache) ProcessGrammar(n Node) Node {
	for _, grammar := range c.Grammar {
		n.Value = fmt.Sprintf("%v :: %v", grammar.Name, n.Value)
	}
	return n
}
