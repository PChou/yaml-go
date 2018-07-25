package yaml

// #include <stdlib.h>
// #include "yaml.h"
// #include "shim-go.h"
import "C"
import (
	"errors"
	//"fmt"
	"unsafe"
)

type YamlObject interface {
	Compare(other YamlObject) bool
}

type YamlMap struct {
	Inner map[string]YamlObject
}

type YamlArray struct {
	Inner []YamlObject
}

type YamlString struct {
	Inner string
}

func (yString *YamlString) Compare(other YamlObject) bool {
	if other == nil {
		return false
	}

	if oString, ok := other.(*YamlString); ok {
		return oString.Inner == yString.Inner
	} else {
		return false
	}
}

func (yMap *YamlMap) Compare(other YamlObject) bool {
	if other == nil {
		return false
	}

	if oMap, ok := other.(*YamlMap); ok {
		if len(oMap.Inner) != len(yMap.Inner) {
			return false
		}

		for k, v := range yMap.Inner {
			if v.Compare(oMap.Inner[k]) == false {
				return false
			}
		}

		return true
	} else {
		return false
	}
}

func (yArray *YamlArray) Compare(other YamlObject) bool {
	if other == nil {
		return false
	}

	if oArray, ok := other.(*YamlArray); ok {
		if len(oArray.Inner) != len(yArray.Inner) {
			return false
		}

		coArray := make([]YamlObject, len(oArray.Inner))
		copy(coArray, oArray.Inner)

		for _, sObject := range yArray.Inner {
			found := false
			for j, oObject := range coArray {
				if sObject.Compare(oObject) {
					coArray = append(coArray[:j], coArray[j+1:]...)
					found = true
					break
				}
			}

			if !found {
				return false
			}
		}

		return true
	} else {
		return false
	}
}

func ParseYaml(source string) (YamlObject, error) {

	parser := C.yaml_parser_create()
	defer C.yaml_parser_destroy(parser)

	cdoc := C.CString(source)
	defer C.free(unsafe.Pointer(cdoc))
	C.yaml_parser_set_input_string(parser, (*C.uchar)(unsafe.Pointer(cdoc)), C.size_t(len(source)))

	document := C.yaml_parser_load_document(parser)
	defer C.yaml_parser_destroy_document(document)

	root := C.yaml_document_get_root_node(document)
	if root == nil {
		return nil, errors.New(C.GoString(parser.problem))
	}
	return walkCTree(document, root), nil
	//fmt.Println(root._type)
}

type scalarNode struct {
	value  *C.yaml_char_t
	length C.size_t
	style  C.yaml_scalar_style_t
}

type sequenceNode struct {
	items struct {
		// actually *C.yaml_node_item_t
		start *C.int
		end   *C.int
		top   *C.int
	}
	style C.yaml_sequence_style_t
}

type mappingNode struct {
	pairs struct {
		start *C.yaml_node_pair_t
		end   *C.yaml_node_pair_t
		top   *C.yaml_node_pair_t
	}
	style C.yaml_mapping_style_t
}

func walkCTree(document *C.yaml_document_t, tree *C.yaml_node_t) YamlObject {
	if tree == nil {
		return nil
	}

	if tree._type == C.YAML_NO_NODE {
		return nil
	} else if tree._type == C.YAML_SCALAR_NODE {
		ystr := (*scalarNode)(unsafe.Pointer(&tree.data))
		rstr := C.GoBytes(unsafe.Pointer(ystr.value), C.int(ystr.length))
		return &YamlString{string(rstr)}
	} else if tree._type == C.YAML_MAPPING_NODE {
		ymap := (*mappingNode)(unsafe.Pointer(&tree.data))
		rmap := &YamlMap{}
		rmap.Inner = make(map[string]YamlObject)
		start := uintptr(unsafe.Pointer(ymap.pairs.start))
		end := uintptr(unsafe.Pointer(ymap.pairs.top))
		off := unsafe.Sizeof(*ymap.pairs.start)
		for i := start; i < end; i += off {
			pair := (*C.yaml_node_pair_t)(unsafe.Pointer(i))

			keyNode := C.yaml_document_get_node(document, pair.key)
			valNode := C.yaml_document_get_node(document, pair.value)
			key := walkCTree(document, keyNode)
			val := walkCTree(document, valNode)
			if skey, ok := key.(*YamlString); ok {
				rmap.Inner[skey.Inner] = val
			}
		}
		return rmap
	} else if tree._type == C.YAML_SEQUENCE_NODE {
		yseq := (*sequenceNode)(unsafe.Pointer(&tree.data))
		rseq := &YamlArray{}
		rseq.Inner = make([]YamlObject, 0)
		start := uintptr(unsafe.Pointer(yseq.items.start))
		end := uintptr(unsafe.Pointer(yseq.items.top))
		off := unsafe.Sizeof(*yseq.items.start)
		for i := start; i < end; i += off {
			item := (*C.int)(unsafe.Pointer(i))
			itemNode := C.yaml_document_get_node(document, *item)
			rseq.Inner = append(rseq.Inner, walkCTree(document, itemNode))
		}
		return rseq
	}

	return nil

}
