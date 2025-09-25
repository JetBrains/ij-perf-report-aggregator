package tc_properties

import (
	"fmt"
)

// LoadBytes load reads a buffer into a Properties struct.
func LoadBytes(buf []byte, isExcludedKey func(string) bool) (*Properties, error) {
	if isExcludedKey == nil {
		isExcludedKey = func(string) bool { return false }
	}
	return parse(buf, isExcludedKey)
}

type parser struct {
	lex *lexer
}

func parse(input []byte, isExcludedKey func(string) bool) (*Properties, error) {
	p := &parser{lex: lex(input)}

	properties := NewProperties()
	key := ""

	for {
		token, err := p.expectOneOf(itemComment, itemKey, itemEOF)
		if err != nil {
			return nil, err
		}
		switch token.typ {
		case itemEOF:
			goto done
		case itemComment:
			continue
		case itemKey:
			key = token.val
		}

		token, err = p.expectOneOf(itemValue, itemEOF)
		if err != nil {
			return nil, err
		}
		switch token.typ {
		case itemEOF:
			properties.m[key] = ""
			goto done
		case itemValue:
			if !isExcludedKey(key) {
				properties.m[key] = token.val
			}
		}
	}

done:
	return properties, nil
}

func (p *parser) expectOneOf(expected ...itemType) (item, error) {
	token := p.lex.nextItem()
	for _, v := range expected {
		if token.typ == v {
			return token, nil
		}
	}
	return item{}, fmt.Errorf("properties: Line %d: %s: %s", p.lex.lineNumber(), "unexpected token", token)
}
