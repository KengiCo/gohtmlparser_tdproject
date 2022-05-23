package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// <parse><div ciao="ciao" ></></parse>
// <!-- <parse ciao="nonono" /> -->

// declaration
const htmlWithoutBuffer = `<parse first="dsf" baz="uella" />
`

// <parse thisno="dsf" baz="uella" />
// bacground-image: url('<parse thisno="foo" />k');
// "</styl>"
// <parse thisno="scripnono" ciao="nono" />
// </style>
// <script>
// <parse thisno="stylenono" ciao="nono" />
// </script>
// '<parse property="si" ciao="sisi" />'
// <parse>"<parse property="dsggfs" ciao="nononono" />"</parse>
// </html>
// .style {

// }

var startingBuffer string = "12345678"

var finishingBuffer string = "87654321"

var html string = startingBuffer + htmlWithoutBuffer + finishingBuffer

//DECLARATION VARIABLE
type Stack []string

var tempRaw string = ""
var attributes []string
var from int
var amIrecording bool = false
var opening bool = true

var stack Stack
var stackLastElement string

var validityStack Stack
var validityStackLastElement string

var inQuotes = false

// IsEmpty: check if stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func handleStackLastElement() {
	if len(stack) > 0 {
		stackLastElement = stack[len(stack)-1]
	} else {
		stackLastElement = ""
	}
}

func handleValidityStackLastElement() {
	if len(validityStack) > 0 {
		validityStackLastElement = validityStack[len(validityStack)-1]
	} else {
		validityStackLastElement = ""
	}
}

func handleQuotes(char string, i int) {

	_ = stackLastElement

	handleStackLastElement()

	// management quotes
	if char == `"` && stackLastElement == `"` { // code to be executed if char is " and last element of stack is "
		stack = stack[:len(stack)-1]
		if len(stack) == 0 {
			inQuotes = false
		}
	} else if char == `"` && stackLastElement != `"` { // code to be executed if char is " and last element of stack is not "
		stack = append(stack, char)
		inQuotes = true
	} else if char == `'` && stackLastElement == `'` { // code to be executed if char is ' and last element of stack is '
		stack = stack[0 : len(stack)-1]
		if len(stack) == 0 {
			inQuotes = false
		}
	} else if char == `'` && stackLastElement != `'` { // code to be executed if char is ' and last element of stack is not '
		stack = append(stack, char)
		inQuotes = true
	}

	// management tag script & tag style
	if html[i-8:i+1] == `</script>` && len(stack) == 0 {
		stack = stack[0 : len(stack)-1]
		inQuotes = false
	} else if html[i:i+8] == `<script>` {
		stack = append(stack, "</script>")
		inQuotes = true
	} else if html[i-7:i+1] == `</style>` && len(stack) == 0 {
		stack = stack[0 : len(stack)-1]
		inQuotes = false
	} else if html[i:i+7] == `<style>` {
		stack = append(stack, "</style>")
		inQuotes = true
	}
}

type FinalResult struct {
	Result []Raw `json:"array,omitempty"`
}

// type Property map[string]string `json:"array_text,omitempty"`

type Raw struct {
	RawString  string   `json:"raw,omitempty"`
	From       int      `json:"from,omitempty"`
	To         int      `json:"to,omitempty"`
	Properties []string `json:"properties,omitempty"`
}

func generateResult(tempRaw string, from int, pos int) string {
	// fmt.Println("the row is:", tempRaw, "/n", "from:", from, " to:", pos, "attributes", attrs)
	// attrs = append(attrs, `{property: "foo"}`, `{pro: "faao"}`)

	FinalResult := Raw{
		RawString:  tempRaw,
		From:       from - 8,
		To:         pos - 8,
		Properties: attributes,
	}
	btResult, _ := json.MarshalIndent(&FinalResult, "", "  ")
	fmt.Println(string(btResult))
	return string(btResult)
}

func getAttributes(str string) {
	var attrs = []string{""}
	var props []string

	attrs = strings.Split(str, " ")
	for i := 1; i < len(attrs)-1; i++ {
		if strings.Contains(attrs[i], "=") {
			props = append(props, attrs[i])
		}
	}

	var prop string = "{"
	for i := 0; i < len(props); i++ {
		temp := strings.Split(props[i], "=")
		_ = temp
		prop += temp[0] + ":" + temp[1] + "}"
		attributes = append(attributes, prop)
		prop = "{"
	}
}

func main() {

	fmt.Println("begin")
	// loop string html
	// for pos, char := range html {

	_ = validityStackLastElement

	for pos := 8; pos < len(html)-8; pos++ {
		char := string(html[pos])
		handleQuotes(char, pos)

		if len(validityStack) > 0 {
			validityStackLastElement = validityStack[len(validityStack)-1]
		} else {
			validityStackLastElement = ""
		}
		if inQuotes == false {
			if opening == true {
				if char == "<" && html[pos:pos+7] == "<parse>" {
					validityStack = append(validityStack, "</parse>")
					opening = false
					amIrecording = true
					from = pos
				} else if char == "<" && html[pos:pos+7] == "<parse " {
					validityStack = append(validityStack, ">")
					opening = false
					amIrecording = true
					from = pos
				} else {
					continue
				}
			} else if opening == false {
				if html[pos-7:pos+1] == validityStackLastElement {
					opening = true
					validityStack = validityStack[:len(validityStack)-1]
					tempRaw += char
					getAttributes(tempRaw)
					generateResult(tempRaw, from, pos)
					attributes = attributes[0:0]
					tempRaw = ""
					amIrecording = false
					continue
				} else if char == validityStackLastElement {
					validityStack = validityStack[:len(validityStack)-1]
					tempRaw += char
					getAttributes(tempRaw)
					if lastChar := string(html[pos-1]); lastChar == "/" {
						generateResult(tempRaw, from, pos)
						attributes = attributes[0:0]
						opening = true
						tempRaw = ""
						amIrecording = false
					} else {
						validityStack = append(validityStack, "</parse>")
					}
					continue
				}
			}
		}
		if amIrecording == true {
			tempRaw += char
		}

	}
	fmt.Println("end")
}
