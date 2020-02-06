package data

// Cache struct
type Cache struct {
	Grammar []Grammar
}

// HasGrammar func
func (c *Cache) HasGrammar() bool {
	return len(c.Grammar) > 0
}

// ProcessGrammar func
func (c *Cache) ProcessGrammar(n Node, extension string) Node {
	for _, grammar := range c.Grammar {
		if grammar.hasExtension(extension) {
			n = grammar.process(n)
		}
	}
	return n
}
