package main

import (
	"fmt"

	"github.com/fossil/fossil-delta-go/delta"
)

func main() {
	orgin := "{\"glossary\":{\"title\":\"example glossary\",\"GlossDiv\":{\"title\":\"S\",\"GlossList\":{\"GlossEntry\":{\"ID\":\"SGML\",\"SortAs\":\"SGML\",\"GlossTerm\":\"Standard Generalized Markup Language\",\"Acronym\":\"SGML\",\"Abbrev\":\"ISO 8879:1986\",\"GlossDef\":{\"para\":\"A meta-markup language, used to create markup languages such as DocBook.\",\"GlossSeeAlso\":[\"GML\",\"XML\"]},\"GlossSee\":\"markup\"}}}}}"
	target := "{\"glossary\":{\"title\":\"changed glossary\",\"GlossDiv\":{\"title\":\"S++\",\"GlossList\":{\"GlossEntry\":{\"ID\":\"LMGS\",\"SortAs\":\"LMGS\",\"GlossTerm\":\"Changed Standard Generalized Markup Language next version\",\"Acronym\":\"LMGS\",\"Abbrev\":\"ISO 8889:2006\",\"GlossDef\":{\"para\":\"A meta-markup language, used to create markup languages such as DocBook. Changed version\",\"GlossSeeAlso\":[\"GML\",\"XML\"]},\"GlossSee\":\"markup changed version\"}}}}}"
	fossilDelta := delta.Create(orgin, target)
	fmt.Println(fossilDelta)
	fmt.Println(string(fossilDelta))

	afterApply, error := delta.Apply(orgin, fossilDelta, false)
	fmt.Println(afterApply)
	fmt.Println(string(afterApply))
	fmt.Println(error)
}
