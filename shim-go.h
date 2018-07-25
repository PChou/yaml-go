#ifndef SHIM_GO_H
#define SHIM_GO_H

#ifdef __cplusplus
extern "C" {
#endif

#include "yaml.h"


YAML_DECLARE(yaml_parser_t *)
yaml_parser_create();

YAML_DECLARE(void)
yaml_parser_destroy(yaml_parser_t *parser);

YAML_DECLARE(yaml_document_t *)
yaml_parser_load_document(yaml_parser_t *parser);

YAML_DECLARE(void)
yaml_parser_destroy_document(yaml_document_t *document);

#ifdef __cplusplus
}
#endif

#endif /* #ifndef SHIM_GO_H */