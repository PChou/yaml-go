package yaml

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var _ = fmt.Println

func TestCompare(t *testing.T) {
	y1, _ := ParseYaml(`key: value
other-key: other-value`)
	y2, _ := ParseYaml(`key: value
other-key: other-value  `)

	assert.True(t, y1.Compare(y2))

	y3, _ := ParseYaml(`
app.modules:
- module: system
  hosts: [1,2,3]
  enabled: true
  set:
    - a
    - b
    - c
- module: system2
- module: system3
`)

	y4, _ := ParseYaml(`
app.modules:
- module: system
  hosts: [1,3,2]
  enabled: true
  set:
    - c
    - b
    - a
- module: system3
- module: system2
`)
	assert.True(t, y3.Compare(y4))
}
