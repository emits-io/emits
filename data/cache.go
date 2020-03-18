package data

// Cache struct
type Cache struct {
	GrammarFile []GrammarFile
}

// HasGrammar func
func (c *Cache) HasGrammar() bool {
	return len(c.GrammarFile) > 0
}

// ProcessGrammar func
func (c *Cache) ProcessGrammar(n Node, extension string) Node {
	for _, grammarFile := range c.GrammarFile {
		if grammarFile.hasExtension(extension) {
			n = grammarFile.process(n)
		}
	}
	return n
}
