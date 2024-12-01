CREATE MATERIALIZED VIEW mv_decision_logs TO decision_logs
AS
SELECT
    operation_id,
    client_id,
    amount,
    trace_id,
    decision,
    decline_reason
FROM kafka_decision_logs;
