"transforms": "unwrap,copyHeaders",

"transforms.unwrap.type": "io.debezium.transforms.ExtractNewRecordState",
"transforms.unwrap.add.fields": "schema,table",

"transforms.copyHeaders.type": "org.apache.kafka.connect.transforms.HeaderFrom$Value",
"transforms.copyHeaders.fields": "__schema,__table",
"transforms.copyHeaders.headers": "schema_name,table_name",
"transforms.copyHeaders.operation": "copy"