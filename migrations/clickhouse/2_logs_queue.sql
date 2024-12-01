CREATE TABLE kafka_decision_logs
(
    operation_id String,
    client_id String,
    amount Int64,
    trace_id String,
    decision String,
    decline_reason Nullable(String)
) ENGINE = Kafka SETTINGS   kafka_broker_list = 'localhost:9091',
                            kafka_topic_list = 'logs',
                            kafka_group_name = 'click/audit_logs',
                            kafka_format = 'JSONEachRow',
                            kafka_skip_broken_messages = 10;
