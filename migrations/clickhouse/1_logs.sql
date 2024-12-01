CREATE TABLE decision_logs
(
    operation_id String,
    client_id String,
    amount Int64,
    trace_id String,
    decision String,
    decline_reason Nullable(String)
) ENGINE = MergeTree
ORDER BY operation_id;
