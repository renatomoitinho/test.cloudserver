"transforms": "unwrap,copyHeaders",

"transforms.unwrap.type": "io.debezium.transforms.ExtractNewRecordState",
"transforms.unwrap.add.fields": "schema,table",

"transforms.copyHeaders.type": "org.apache.kafka.connect.transforms.HeaderFrom$Value",
"transforms.copyHeaders.fields": "__schema,__table",
"transforms.copyHeaders.headers": "schema_name,table_name",
"transforms.copyHeaders.operation": "copy"



"transforms": "unwrap,copyHeaders,fixSchemaHeader",

"transforms.unwrap.type": "io.debezium.transforms.ExtractNewRecordState",
"transforms.unwrap.add.fields": "schema,table",

"transforms.copyHeaders.type": "org.apache.kafka.connect.transforms.HeaderFrom$Value",
"transforms.copyHeaders.fields": "__schema,__table",
"transforms.copyHeaders.headers": "schema_name,table_name",
"transforms.copyHeaders.operation": "copy",

"transforms.fixSchemaHeader.type": "com.sheeps.kafka.connect.transforms.RegexHeader",
"transforms.fixSchemaHeader.input.header": "schema_name",
"transforms.fixSchemaHeader.output.header": "schema_name",
"transforms.fixSchemaHeader.regex": "^.*_([^_]+)$",
"transforms.fixSchemaHeader.replacement": "$1"