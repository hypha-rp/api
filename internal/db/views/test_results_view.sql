CREATE VIEW test_results_view AS
SELECT
    r.id::text AS result_id,
    r.product_id::text AS product_id,
    ts.id::text AS test_suite_id,
    ts.name AS test_suite_name,
    ts.tests AS test_suite_tests,
    ts.failures AS test_suite_failures,
    ts.errors AS test_suite_errors,
    ts.skipped AS test_suite_skipped,
    ts.assertions AS test_suite_assertions,
    ts.time AS test_suite_time,
    ts.file AS test_suite_file,
    ts.system_out AS test_suite_system_out,
    ts.system_err AS test_suite_system_err,
    tc.id::text AS test_case_id,
    tc.name AS test_case_name,
    tc.class_name AS test_case_class_name,
    tc.time AS test_case_time,
    tc.status AS test_case_status,
    tc.message AS test_case_message,
    tc.type AS test_case_type,
    tc.assertions AS test_case_assertions,
    tc.file AS test_case_file,
    tc.line AS test_case_line,
    tc.system_out AS test_case_system_out,
    tc.system_err AS test_case_system_err
FROM
    results r
JOIN
    test_suites ts ON r.id::text = ts.result_id::text
LEFT JOIN
    test_cases tc ON ts.id::text = tc.test_suite_id::text;