#include "shim-go.h"

YAML_DECLARE(yaml_parser_t *)
yaml_parser_create() {
	yaml_parser_t *parser;
	parser = (yaml_parser_t *)malloc(sizeof(yaml_parser_t));
	if ( yaml_parser_initialize(parser) > 0 ) {
		return parser;
	} else {
		free(parser);
		return NULL;
	}
}

YAML_DECLARE(void)
yaml_parser_destroy(yaml_parser_t *parser) {
	if ( parser == NULL )
		return;
	yaml_parser_delete(parser);
	free(parser);
}


YAML_DECLARE(yaml_document_t *)
yaml_parser_load_document(yaml_parser_t *parser) {
	yaml_document_t *doc;
	doc = (yaml_document_t *)malloc(sizeof(yaml_document_t));
	yaml_parser_load(parser, doc);
	return doc;
}

YAML_DECLARE(void)
yaml_parser_destroy_document(yaml_document_t *document) {
	if ( document == NULL )
		return;
	yaml_document_delete(document);
	free(document);
}