# yaml-go

go binding for [libyaml](https://github.com/yaml/libyaml), which is a pure ANSI C implementation of yaml generator and parser

in this binding, `YamlObject` is introduced, and that can be compared:

```go
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
```

`cgo` must be enabled