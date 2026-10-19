//go:build integration

package engagelab

import "testing"

func fullQA_devices(c *fullQA) {
	var a any
	_ = a
	var alias any
	_ = alias
	var b any
	_ = b
	var binding any
	_ = binding
	var d1 any
	_ = d1
	var d2 any
	_ = d2
	var r any
	_ = r
	var raw any
	_ = raw
	var tag1 any
	_ = tag1
	var tag2 any
	_ = tag2
	var tag3 any
	_ = tag3
	c.write()
	d1 = c.seedDevice()
	d2 = c.seedDevice()
	tag1 = qcat(c.prefix, "_a")
	tag2 = qcat(c.prefix, "_b")
	tag3 = qcat(c.prefix, "_del")
	alias = qcat(c.prefix, "_alias")
	c.rest("POST", qcat("/v4/devices/", qtxt(d1)), qm("tags", "", "alias", ""))
	c.seedTag(tag1, d1)
	c.seedAlias(alias, d1)
	r = c.sdk("device_get", qarr(d1))
	raw = c.device(d1)
	c.record("qa_48e10c8c3688", qm("HTTP200", qeq(qg(r, "status"), 200), "known_tags", qsets(qg(r, "data", "tags"), qarr(tag1)), "known_alias", qeq(qg(r, "data", "alias"), alias), "wire_fields", qsame(qg(r, "data"), qg(r, "raw"))), qyes(true), "")
	r = c.sdk("alias_get", qarr(alias, qarr("android")))
	c.record("qa_4228ae60b590", qm("HTTP200", qeq(qg(r, "status"), 200), "exact_binding", qsets(qg(r, "data", "registration_ids"), qarr(d1)), "wire_fields", qsame(qg(r, "data"), qg(r, "raw"))), qyes(true), "")
	r = c.sdk("alias_delete", qarr(alias, qarr("android")))
	raw = c.device(d1)
	binding = c.rest("GET", qcat(qcat("/v4/aliases/", alias), "?platform=android"), nil)
	c.record("qa_f76445835b7d", qm("SDK_success", qyes(qg(r, "ok")), "device_alias_empty", qeq(qtxt(qg(raw, "data", "alias")), ""), "alias_unbound", qno(qincludes(qg(binding, "data", "registration_ids"), d1))), qyes(true), "")
	c.rest("POST", qcat("/v4/devices/", qtxt(d1)), qm("tags", ""))
	c.seedTag(tag2, d1)
	r = c.sdk("device_set", qarr(d1, qm("tags", qm("add", qarr(tag1), "remove", qarr(tag2)), "alias", alias)))
	raw = c.device(d1)
	c.record("qa_3a16670b310f", qm("SDK_success", qyes(qg(r, "ok")), "tags_add_remove", qsets(qg(raw, "data", "tags"), qarr(tag1)), "alias_written", qeq(qg(raw, "data", "alias"), alias)), qyes(true), "")
	c.seedTag(tag2, d1)
	r = c.sdk("device_set", qarr(d1, qm("tags", "")))
	raw = c.device(d1)
	a = c.relation(tag1, d1)
	b = c.relation(tag2, d1)
	c.record("qa_418e812084ee", qm("SDK_success", qyes(qg(r, "ok")), "tags_empty", qeq(qsize(qg(raw, "data", "tags")), 0), "both_false", qboth(qno(qg(a, "data", "result")), qno(qg(b, "data", "result")))), qyes(true), "")
	c.seedTag(tag1, d1)
	c.seedTag(tag2, d2)
	r = c.sdk("tags", qarr())
	c.record("qa_00cbd3d800bd", qm("HTTP200", qeq(qg(r, "status"), 200), "two_seeded_tags", qboth(qincludes(qg(r, "data", "tags"), tag1), qincludes(qg(r, "data", "tags"), tag2))), qyes(true), "")
	r = c.sdk("tag_status", qarr(tag1, d1))
	c.record("qa_a925a3248bcd", qm("HTTP200", qeq(qg(r, "status"), 200), "typed_true", qyes(qg(r, "data", "result"))), qyes(true), "")
	r = c.sdk("tag_count", qarr(qarr(tag1, tag2), "android"))
	c.record("qa_d5581d7d5601", qm("HTTP200", qeq(qg(r, "status"), 200), "two_known_single_counts", qboth(qeq(qg(r, "data", "tagsCount", tag1), 1), qeq(qg(r, "data", "tagsCount", tag2), 1)), "repeated_tags_query", qsets(qg(r, "query", "tags"), qarr(tag1, tag2))), qyes(true), "")
	r = c.sdk("tag_quota", qarr(qarr(tag1, tag2), "android"))
	c.record("qa_5f8b9b9d7f78", qm("HTTP200", qeq(qg(r, "status"), 200), "repeated_tags_query", qsets(qg(r, "query", "tags"), qarr(tag1, tag2)), "nonempty_known_quotas", qboth(qincludes(qids(qg(r, "raw", "data", "tagUidQuotaDetail"), "tagName"), tag1), qincludes(qids(qg(r, "raw", "data", "tagUidQuotaDetail"), "tagName"), tag2)), "wire_fields", qsame(qg(r, "data"), qg(r, "raw"))), qyes(true), "")
	c.rest("POST", qcat("/v4/devices/", qtxt(d1)), qm("tags", ""))
	c.rest("POST", qcat("/v4/devices/", qtxt(d2)), qm("tags", ""))
	c.seedTag(tag1, d2)
	r = c.sdk("tag_set", qarr(tag1, qm("registration_ids", qm("add", qarr(d1), "remove", qarr(d2)))))
	a = c.relation(tag1, d1)
	b = c.relation(tag1, d2)
	c.record("qa_83ba432b89af", qm("SDK_success", qyes(qg(r, "ok")), "A_true_B_false", qboth(qyes(qg(a, "data", "result")), qno(qg(b, "data", "result")))), qyes(true), "")
	c.seedTag(tag3, d1)
	r = c.sdk("tag_delete", qarr(tag3, qarr("android")))
	raw = c.rest("GET", "/v4/tags/", nil)
	a = c.relation(tag3, d1)
	c.record("qa_ff6fe44cc46d", qm("SDK_success", qyes(qg(r, "ok")), "tag_absent", qno(qincludes(qg(raw, "data", "tags"), tag3)), "relation_false", qno(qg(a, "data", "result"))), qyes(true), "")
	r = c.sdk("device_status", qarr(qm("registration_ids", qarr(d1, d2))))
	c.record("qa_edfec6f68242", qm("HTTP200", qeq(qg(r, "status"), 200), "two_known_offline_records", qboth(qeq(qsize(qg(r, "data")), 2), qsets(qids(qg(r, "data"), "regid"), qarr(d1, d2)), qno(qg(r, "data", 0, "online")), qno(qg(r, "data", 1, "online")))), qyes(true), "")
	r = c.sdk("device_delete", qarr(d2))
	raw = c.device(d2)
	a = c.relation(tag2, d2)
	c.record("qa_42caca08cb39", qm("SDK_success", qyes(qg(r, "ok")), "device_absent", qeq(qg(raw, "data", "error", "code"), 27002), "old_relations_inaccessible", qeq(qg(a, "data", "error", "code"), 27002)), qyes(true), "")
}

func fullQA_plans(c *fullQA) {
	var a any
	_ = a
	var b any
	_ = b
	var control any
	_ = control
	var description any
	_ = description
	var first any
	_ = first
	var pair any
	_ = pair
	var pid any
	_ = pid
	var prefix any
	_ = prefix
	var r any
	_ = r
	var raw any
	_ = raw
	c.write()
	first = c.seedPlans(101)
	prefix = c.prefix
	r = c.sdk("plan_list", qarr(1, 101, nil, prefix))
	c.record("qa_bbef5f838820", qm("SDK_HTTP200", qeq(qg(first, "status"), 200), "created_readback", qincludes(qids(qg(r, "data", "push_plan_info"), "plan_id"), qg(first, "data", "plan_id"))), qyes(true), "")
	c.record("qa_0e01ec2f278a", qm("HTTP200", qeq(qg(r, "status"), 200), "count101", qeq(qg(r, "data", "total"), 101), "all_fields", qsame(qg(r, "data"), qg(r, "raw"))), qyes(true), "")
	c.record("qa_2d67bd95bf58", qm("HTTP200", qeq(qg(r, "status"), 200), "total101", qeq(qg(r, "data", "total"), 101), "count101", qeq(qsize(qg(r, "data", "push_plan_info")), 101), "owned_ids_exact", qsets(qids(qg(r, "data", "push_plan_info"), "plan_id"), c.planIDs()), "wire_fields", qsame(qg(r, "data"), qg(r, "raw")), "query101", qeq(qg(r, "query", "page_size", 0), "101")), qyes(true), "")
	r = c.sdk("plan_list", qarr(1, 100, nil, prefix))
	c.record("qa_452338f7f0cb", qm("HTTP200", qeq(qg(r, "status"), 200), "page100", qeq(qsize(qg(r, "data", "push_plan_info")), 100), "wire_fields", qsame(qg(r, "data"), qg(r, "raw"))), qyes(true), "")
	a = c.sdk("plan_list", qarr(1, 1, nil, prefix))
	c.record("qa_d1c2dc88ee2a", qm("HTTP200", qeq(qg(a, "status"), 200), "one_item", qeq(qsize(qg(a, "data", "push_plan_info")), 1), "wire_fields", qsame(qg(a, "data"), qg(a, "raw"))), qyes(true), "")
	b = c.sdk("plan_list", qarr(2, 1, nil, prefix))
	c.record("qa_1068750e639d", qm("HTTP200", qeq(qg(b, "status"), 200), "one_item", qeq(qsize(qg(b, "data", "push_plan_info")), 1), "different_from_first", qno(qeq(qg(a, "data", "push_plan_info", 0, "plan_id"), qg(b, "data", "push_plan_info", 0, "plan_id"))), "wire_fields", qsame(qg(b, "data"), qg(b, "raw"))), qyes(true), "")
	r = c.sdk("plan_list", qarr(1, 0, nil, prefix))
	control = c.listPlans(prefix, 0, 1, nil)
	c.record("qa_c32bda243404", qm("SDK_rejected", qno(qg(r, "ok")), "HTTP_code_message_match", c.errorsEqual(r, control)), qyes(true), "")
	pid = qg(first, "data", "plan_id")
	description = qcat(prefix, "_更新")
	r = c.sdk("plan_create", qarr(qm("plan_id", pid, "plan_description", description)))
	raw = c.listPlans(description, 101, 1, nil)
	c.record("qa_ee67b55b2f74", qm("SDK_success", qyes(qg(r, "ok")), "single_updated", qboth(qeq(qg(raw, "data", "total"), 1), qeq(qg(raw, "data", "push_plan_info", 0, "plan_description"), description))), qyes(true), "")
	r = c.sdk("plan_list", qarr(1, 101, nil, description))
	c.record("qa_3e98ca439abe", qm("HTTP200", qeq(qg(r, "status"), 200), "one_UTF8_match", qboth(qeq(qg(r, "data", "total"), 1), qeq(qg(r, "data", "push_plan_info", 0, "plan_description"), description)), "UTF8_query", qeq(qg(r, "query", "search_description", 0), description)), qyes(true), "")
	r = c.sdk("plan_list", qarr(1, 101, 0, prefix))
	raw = c.listPlans(prefix, 101, 1, 0)
	c.record("qa_a8afa2ba63c2", qm("HTTP200", qeq(qg(r, "status"), 200), "API_known_records", qeq(qg(r, "data", "total"), 101), "source_query", qeq(qg(r, "query", "send_source", 0), "0"), "REST_agrees", qsame(qg(r, "data"), qg(raw, "data"))), qyes(true), "")
	r = c.sdk("plan_list", qarr(1, 1, 0, description))
	c.record("qa_6df4ebb6c583", qm("HTTP200", qeq(qg(r, "status"), 200), "all_four_queries", qboth(qeq(qg(r, "query", "page_index", 0), "1"), qeq(qg(r, "query", "page_size", 0), "1"), qeq(qg(r, "query", "send_source", 0), "0"), qeq(qg(r, "query", "search_description", 0), description)), "single_correct_item", qboth(qeq(qg(r, "data", "total"), 1), qeq(qg(r, "data", "push_plan_info", 0, "plan_id"), pid))), qyes(true), "")
	r = c.sdk("plan_delete", qarr(pid))
	raw = c.listPlans(prefix, 101, 1, nil)
	c.record("qa_793b07888b66", qm("SDK_success", qyes(qg(r, "ok")), "plan_absent", qno(qincludes(qids(qg(raw, "data", "push_plan_info"), "plan_id"), pid))), qyes(true), "")
	pair = c.batchPlanPair(pid)
	r = c.sdk("plan_batch_delete", qarr(c.join(pair)))
	raw = c.listPlans(prefix, 101, 1, nil)
	c.record("qa_7c9f085a13ea", qm("SDK_success", qyes(qg(r, "ok")), "comma_path", qeq(qg(r, "path"), qcat("/v4/push_plan/batch/", c.join(pair))), "both_absent", c.noneIn(pair, qids(qg(raw, "data", "push_plan_info"), "plan_id"))), qyes(true), "")
	r = c.planBoundary(1)
	c.record("qa_0df5c6ed2605", qm("expected_outcome", qyes(qg(r, "verified"))), qyes(qyes(qg(r, "ready"))), "No unused one-character plan ID available")
	r = c.planBoundary(49)
	c.record("qa_2e34db0cb38c", qm("expected_outcome", qyes(qg(r, "verified"))), qyes(qyes(qg(r, "ready"))), "No unused one-character plan ID available")
	r = c.planBoundary(50)
	c.record("qa_eeb26ebb2b07", qm("expected_outcome", qyes(qg(r, "verified"))), qyes(qyes(qg(r, "ready"))), "No unused one-character plan ID available")
	r = c.planBoundary(51)
	c.record("qa_d71d62850538", qm("expected_outcome", qyes(qg(r, "verified"))), qyes(qyes(qg(r, "ready"))), "No unused one-character plan ID available")
}

func fullQA_batch(c *fullQA) {
	var alias1 any
	_ = alias1
	var alias2 any
	_ = alias2
	var d1 any
	_ = d1
	var d2 any
	_ = d2
	var r any
	_ = r
	c.write()
	d1 = c.targetRID(0)
	d2 = c.targetRID(1)
	alias1 = qcat(c.prefix, "_1")
	alias2 = qcat(c.prefix, "_2")
	c.seedAlias(alias1, d1)
	c.seedAlias(alias2, d2)
	r = c.batch("batch_regid", qarr(d1, d2))
	c.record("qa_81f620c4cc20", qm("HTTP200", qeq(qg(r, "status"), 200), "all_target_successes", c.batchResultsValid(r, qarr(d1, d2)), "wire_fields", qsame(qg(r, "data"), qg(r, "raw")), "message_linkage", qyes(qg(r, "messages_persisted"))), qyes(true), "")
	r = c.batch("batch_alias", qarr(alias1, alias2))
	c.record("qa_d7b1b6768d39", qm("HTTP200", qeq(qg(r, "status"), 200), "all_target_successes", c.batchResultsValid(r, qarr(alias1, alias2)), "wire_fields", qsame(qg(r, "data"), qg(r, "raw")), "message_linkage", qyes(qg(r, "messages_persisted"))), qyes(true), "")
}

func fullQA_negative(c *fullQA) {
	var r any
	_ = r
	r = c.negative("empty")
	c.record("qa_f0179b307f6b", qm("rejected_like_REST", qyes(qg(r, "verified"))), qyes(true), "")
	r = c.negative("platform")
	c.record("qa_65322fbec684", qm("legacy_error_like_REST", qyes(qg(r, "verified"))), qyes(true), "")
	r = c.negative("message")
	c.record("qa_a62eda72d0a4", qm("single_type_controls_valid", qyes(qg(r, "controls_valid")), "same_HTTP_code_message", qyes(qg(r, "verified"))), qyes(qyes(qg(r, "controls_valid"))), "Cannot validate mutual exclusion without accepted single-type controls")
	r = c.negative("notification")
	c.record("qa_358fd717fbad", qm("single_type_controls_valid", qyes(qg(r, "controls_valid")), "same_HTTP_code_message", qyes(qg(r, "verified"))), qyes(qyes(qg(r, "controls_valid"))), "Cannot validate mutual exclusion without accepted single-type controls")
	r = c.negative("live_activity")
	c.record("qa_f27bda2df7c2", qm("single_type_controls_valid", qyes(qg(r, "controls_valid")), "same_HTTP_code_message", qyes(qg(r, "verified"))), qyes(qyes(qg(r, "controls_valid"))), "Cannot validate mutual exclusion without accepted single-type controls")
}

func fullQA_auth(c *fullQA) {
	var r any
	_ = r
	r = c.sdk("vip", qarr())
	c.record("qa_861607f5c661", qm("HTTP200", qeq(qg(r, "status"), 200), "integer_status_and_64bit_expiry", c.vipValid(r), "wire_fields", qsame(qg(r, "data"), qg(r, "raw"))), qyes(true), "")
	r = c.invalidAuth()
	c.record("qa_99b351df6c47", qm("REST_and_SDK_reject", qyes(qg(r, "verified"))), qyes(true), "")
}

func fullQA_token(c *fullQA) {
	var r any
	_ = r
	c.write()
	r = c.apnsToken()
	c.record("qa_113173a66889", qm("fields_match_wire", qyes(qg(r, "verified")), "mapping_preserved_or_cleaned", qyes(qg(r, "cleanup"))), qyes(true), "")
}

func fullQA_group(c *fullQA) {
	var r any
	_ = r
	c.write()
	r = c.groupPush()
	c.record("qa_992ab0e0747a", qm("two_dynamic_members", qyes(qg(r, "verified")), "explicit_alias_only", qyes(qg(r, "explicit_target"))), qyes(true), "")
}

func fullQA_image(c *fullQA) {
	var r any
	_ = r
	c.write()
	r = c.imageUpload()
	c.record("qa_b0478ec2b424", qm("eligible_984x369_resource", qyes(qg(r, "resource_valid")), "HTTP_and_picture_ID_exact", qyes(qg(r, "verified"))), qyes(qyes(qg(r, "prerequisite"))), "QA/vendor image capability unavailable; inspect native response code")
}

func fullQA_voice(c *fullQA) {
	var r any
	_ = r
	c.write()
	r = c.prepareVoice()
	c.record("qa_bd03c2bacf73", qm("multipart_and_fields", qyes(qg(r, "upload_verified")), "detail_matches_upload", qyes(qg(r, "detail_verified"))), qyes(true), "")
	c.record("qa_a4be39a4d66b", qm("direct_array_and_known_template", qyes(qg(r, "list_verified"))), qyes(true), "")
	c.record("qa_c7377d6722fc", qm("file_url_matches", qyes(qg(r, "detail_verified"))), qyes(true), "")
	r = c.deleteVoice()
	c.record("qa_95007431d632", qm("SDK_success", qyes(qg(r, "sdk_ok")), "list_and_detail_absent", qyes(qg(r, "verified"))), qyes(true), "")
}

func fullQA_schedule(c *fullQA) {
	var r any
	_ = r
	c.write()
	r = c.prepareSchedules()
	c.record("qa_cc992174acf7", qm("native_business_assertions", qyes(qg(r, "create"))), qyes(qyes(qg(r, "create_ready"))), "Scheduled task, timezone/activity, or positive trigger prerequisite not established")
	c.record("qa_fbe62ef8b3e7", qm("native_business_assertions", qyes(qg(r, "get"))), qyes(qyes(qg(r, "get_ready"))), "Scheduled task, timezone/activity, or positive trigger prerequisite not established")
	c.record("qa_c69b89ef022d", qm("native_business_assertions", qyes(qg(r, "list"))), qyes(qyes(qg(r, "list_ready"))), "Scheduled task, timezone/activity, or positive trigger prerequisite not established")
	c.record("qa_f47503d90583", qm("native_business_assertions", qyes(qg(r, "intelligent"))), qyes(qyes(qg(r, "intelligent_ready"))), "Scheduled task, timezone/activity, or positive trigger prerequisite not established")
	c.record("qa_b1336ed19443", qm("native_business_assertions", qyes(qg(r, "messages"))), qyes(qyes(qg(r, "messages_ready"))), "Scheduled task, timezone/activity, or positive trigger prerequisite not established")
	c.record("qa_0359444e46a9", qm("native_business_assertions", qyes(qg(r, "update"))), qyes(qyes(qg(r, "update_ready"))), "Scheduled task, timezone/activity, or positive trigger prerequisite not established")
	c.record("qa_cbbfd7d59d01", qm("native_business_assertions", qyes(qg(r, "delete"))), qyes(qyes(qg(r, "delete_ready"))), "Scheduled task, timezone/activity, or positive trigger prerequisite not established")
}

func fullQA_stats(c *fullQA) {
	var r any
	_ = r
	c.write()
	r = c.prepareStatistics()
	c.record("qa_06327bbf436b", qm("nonempty_known_business_data_and_fields", qyes(qg(r, "message"))), qyes(qyes(qg(r, "message_ready"))), "Known message/RID/plan relation or nonzero HMOS statistics not available in bounded observation window")
	c.record("qa_cc281fbd1340", qm("nonempty_known_business_data_and_fields", qyes(qg(r, "lifecycle"))), qyes(qyes(qg(r, "lifecycle_ready"))), "Known message/RID/plan relation or nonzero HMOS statistics not available in bounded observation window")
	c.record("qa_6f9ed781a410", qm("nonempty_known_business_data_and_fields", qyes(qg(r, "batch_lifecycle"))), qyes(qyes(qg(r, "batch_lifecycle_ready"))), "Known message/RID/plan relation or nonzero HMOS statistics not available in bounded observation window")
	c.record("qa_0c64e79794cc", qm("nonempty_known_business_data_and_fields", qyes(qg(r, "users"))), qyes(qyes(qg(r, "users_ready"))), "Known message/RID/plan relation or nonzero HMOS statistics not available in bounded observation window")
	c.record("qa_29349c0683f5", qm("nonempty_known_business_data_and_fields", qyes(qg(r, "plan"))), qyes(qyes(qg(r, "plan_ready"))), "Known message/RID/plan relation or nonzero HMOS statistics not available in bounded observation window")
	c.record("qa_cb8dc4d08772", qm("nonempty_known_business_data_and_fields", qyes(qg(r, "plan_messages"))), qyes(qyes(qg(r, "plan_messages_ready"))), "Known message/RID/plan relation or nonzero HMOS statistics not available in bounded observation window")
}

func fullQA_physical(c *fullQA) {
	var r any
	_ = r
	r = c.physicalNotifications()
	c.record("qa_e9617155db83", qm("native_API_and_physical_notification_assertions", qyes(qg(r, "push"))), qyes(qyes(qg(r, "push_ready"))), "Matching QA application and physical notification observation are required")
	c.record("qa_f5308cfc2d09", qm("native_API_and_physical_notification_assertions", qyes(qg(r, "validate"))), qyes(qyes(qg(r, "validate_ready"))), "Matching QA application and physical notification observation are required")
	c.record("qa_db1a22ee70da", qm("native_API_and_physical_notification_assertions", qyes(qg(r, "withdraw"))), qyes(qyes(qg(r, "withdraw_ready"))), "Matching QA application and physical notification observation are required")
}

func TestFullIntegration(t *testing.T) {
	cache := map[string]map[string]any{}
	t.Run("qa_f76445835b7d", func(t *testing.T) {
		fullExecute(t, cache, "devices", "qa_f76445835b7d", []string{"qa_48e10c8c3688", "qa_4228ae60b590", "qa_f76445835b7d", "qa_3a16670b310f", "qa_418e812084ee", "qa_00cbd3d800bd", "qa_a925a3248bcd", "qa_d5581d7d5601", "qa_5f8b9b9d7f78", "qa_83ba432b89af", "qa_ff6fe44cc46d", "qa_edfec6f68242", "qa_42caca08cb39"}, fullQA_devices)
	})
	t.Run("qa_4228ae60b590", func(t *testing.T) {
		fullExecute(t, cache, "devices", "qa_4228ae60b590", []string{"qa_48e10c8c3688", "qa_4228ae60b590", "qa_f76445835b7d", "qa_3a16670b310f", "qa_418e812084ee", "qa_00cbd3d800bd", "qa_a925a3248bcd", "qa_d5581d7d5601", "qa_5f8b9b9d7f78", "qa_83ba432b89af", "qa_ff6fe44cc46d", "qa_edfec6f68242", "qa_42caca08cb39"}, fullQA_devices)
	})
	t.Run("qa_d7b1b6768d39", func(t *testing.T) {
		fullExecute(t, cache, "batch", "qa_d7b1b6768d39", []string{"qa_81f620c4cc20", "qa_d7b1b6768d39"}, fullQA_batch)
	})
	t.Run("qa_81f620c4cc20", func(t *testing.T) {
		fullExecute(t, cache, "batch", "qa_81f620c4cc20", []string{"qa_81f620c4cc20", "qa_d7b1b6768d39"}, fullQA_batch)
	})
	t.Run("qa_42caca08cb39", func(t *testing.T) {
		fullExecute(t, cache, "devices", "qa_42caca08cb39", []string{"qa_48e10c8c3688", "qa_4228ae60b590", "qa_f76445835b7d", "qa_3a16670b310f", "qa_418e812084ee", "qa_00cbd3d800bd", "qa_a925a3248bcd", "qa_d5581d7d5601", "qa_5f8b9b9d7f78", "qa_83ba432b89af", "qa_ff6fe44cc46d", "qa_edfec6f68242", "qa_42caca08cb39"}, fullQA_devices)
	})
	t.Run("qa_edfec6f68242", func(t *testing.T) {
		fullExecute(t, cache, "devices", "qa_edfec6f68242", []string{"qa_48e10c8c3688", "qa_4228ae60b590", "qa_f76445835b7d", "qa_3a16670b310f", "qa_418e812084ee", "qa_00cbd3d800bd", "qa_a925a3248bcd", "qa_d5581d7d5601", "qa_5f8b9b9d7f78", "qa_83ba432b89af", "qa_ff6fe44cc46d", "qa_edfec6f68242", "qa_42caca08cb39"}, fullQA_devices)
	})
	t.Run("qa_48e10c8c3688", func(t *testing.T) {
		fullExecute(t, cache, "devices", "qa_48e10c8c3688", []string{"qa_48e10c8c3688", "qa_4228ae60b590", "qa_f76445835b7d", "qa_3a16670b310f", "qa_418e812084ee", "qa_00cbd3d800bd", "qa_a925a3248bcd", "qa_d5581d7d5601", "qa_5f8b9b9d7f78", "qa_83ba432b89af", "qa_ff6fe44cc46d", "qa_edfec6f68242", "qa_42caca08cb39"}, fullQA_devices)
	})
	t.Run("qa_418e812084ee", func(t *testing.T) {
		fullExecute(t, cache, "devices", "qa_418e812084ee", []string{"qa_48e10c8c3688", "qa_4228ae60b590", "qa_f76445835b7d", "qa_3a16670b310f", "qa_418e812084ee", "qa_00cbd3d800bd", "qa_a925a3248bcd", "qa_d5581d7d5601", "qa_5f8b9b9d7f78", "qa_83ba432b89af", "qa_ff6fe44cc46d", "qa_edfec6f68242", "qa_42caca08cb39"}, fullQA_devices)
	})
	t.Run("qa_3a16670b310f", func(t *testing.T) {
		fullExecute(t, cache, "devices", "qa_3a16670b310f", []string{"qa_48e10c8c3688", "qa_4228ae60b590", "qa_f76445835b7d", "qa_3a16670b310f", "qa_418e812084ee", "qa_00cbd3d800bd", "qa_a925a3248bcd", "qa_d5581d7d5601", "qa_5f8b9b9d7f78", "qa_83ba432b89af", "qa_ff6fe44cc46d", "qa_edfec6f68242", "qa_42caca08cb39"}, fullQA_devices)
	})
	t.Run("qa_992ab0e0747a", func(t *testing.T) {
		fullExecute(t, cache, "group", "qa_992ab0e0747a", []string{"qa_992ab0e0747a"}, fullQA_group)
	})
	t.Run("qa_b0478ec2b424", func(t *testing.T) {
		fullExecute(t, cache, "image", "qa_b0478ec2b424", []string{"qa_b0478ec2b424"}, fullQA_image)
	})
	t.Run("qa_6df4ebb6c583", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_6df4ebb6c583", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_1068750e639d", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_1068750e639d", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_bbef5f838820", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_bbef5f838820", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_793b07888b66", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_793b07888b66", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_7c9f085a13ea", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_7c9f085a13ea", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_3e98ca439abe", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_3e98ca439abe", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_ee67b55b2f74", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_ee67b55b2f74", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_a8afa2ba63c2", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_a8afa2ba63c2", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_0e01ec2f278a", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_0e01ec2f278a", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_cb8dc4d08772", func(t *testing.T) {
		fullExecute(t, cache, "stats", "qa_cb8dc4d08772", []string{"qa_06327bbf436b", "qa_cc281fbd1340", "qa_6f9ed781a410", "qa_0c64e79794cc", "qa_29349c0683f5", "qa_cb8dc4d08772"}, fullQA_stats)
	})
	t.Run("qa_0df5c6ed2605", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_0df5c6ed2605", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_2e34db0cb38c", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_2e34db0cb38c", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_eeb26ebb2b07", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_eeb26ebb2b07", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_d71d62850538", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_d71d62850538", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_c32bda243404", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_c32bda243404", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_d1c2dc88ee2a", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_d1c2dc88ee2a", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_452338f7f0cb", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_452338f7f0cb", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_2d67bd95bf58", func(t *testing.T) {
		fullExecute(t, cache, "plans", "qa_2d67bd95bf58", []string{"qa_bbef5f838820", "qa_0e01ec2f278a", "qa_2d67bd95bf58", "qa_452338f7f0cb", "qa_d1c2dc88ee2a", "qa_1068750e639d", "qa_c32bda243404", "qa_ee67b55b2f74", "qa_3e98ca439abe", "qa_a8afa2ba63c2", "qa_6df4ebb6c583", "qa_793b07888b66", "qa_7c9f085a13ea", "qa_0df5c6ed2605", "qa_2e34db0cb38c", "qa_eeb26ebb2b07", "qa_d71d62850538"}, fullQA_plans)
	})
	t.Run("qa_f27bda2df7c2", func(t *testing.T) {
		fullExecute(t, cache, "negative", "qa_f27bda2df7c2", []string{"qa_f0179b307f6b", "qa_65322fbec684", "qa_a62eda72d0a4", "qa_358fd717fbad", "qa_f27bda2df7c2"}, fullQA_negative)
	})
	t.Run("qa_a62eda72d0a4", func(t *testing.T) {
		fullExecute(t, cache, "negative", "qa_a62eda72d0a4", []string{"qa_f0179b307f6b", "qa_65322fbec684", "qa_a62eda72d0a4", "qa_358fd717fbad", "qa_f27bda2df7c2"}, fullQA_negative)
	})
	t.Run("qa_358fd717fbad", func(t *testing.T) {
		fullExecute(t, cache, "negative", "qa_358fd717fbad", []string{"qa_f0179b307f6b", "qa_65322fbec684", "qa_a62eda72d0a4", "qa_358fd717fbad", "qa_f27bda2df7c2"}, fullQA_negative)
	})
	t.Run("qa_e9617155db83", func(t *testing.T) {
		fullExecute(t, cache, "physical", "qa_e9617155db83", []string{"qa_e9617155db83", "qa_f5308cfc2d09", "qa_db1a22ee70da"}, fullQA_physical)
	})
	t.Run("qa_65322fbec684", func(t *testing.T) {
		fullExecute(t, cache, "negative", "qa_65322fbec684", []string{"qa_f0179b307f6b", "qa_65322fbec684", "qa_a62eda72d0a4", "qa_358fd717fbad", "qa_f27bda2df7c2"}, fullQA_negative)
	})
	t.Run("qa_db1a22ee70da", func(t *testing.T) {
		fullExecute(t, cache, "physical", "qa_db1a22ee70da", []string{"qa_e9617155db83", "qa_f5308cfc2d09", "qa_db1a22ee70da"}, fullQA_physical)
	})
	t.Run("qa_f5308cfc2d09", func(t *testing.T) {
		fullExecute(t, cache, "physical", "qa_f5308cfc2d09", []string{"qa_e9617155db83", "qa_f5308cfc2d09", "qa_db1a22ee70da"}, fullQA_physical)
	})
	t.Run("qa_f0179b307f6b", func(t *testing.T) {
		fullExecute(t, cache, "negative", "qa_f0179b307f6b", []string{"qa_f0179b307f6b", "qa_65322fbec684", "qa_a62eda72d0a4", "qa_358fd717fbad", "qa_f27bda2df7c2"}, fullQA_negative)
	})
	t.Run("qa_c69b89ef022d", func(t *testing.T) {
		fullExecute(t, cache, "schedule", "qa_c69b89ef022d", []string{"qa_cc992174acf7", "qa_fbe62ef8b3e7", "qa_c69b89ef022d", "qa_f47503d90583", "qa_b1336ed19443", "qa_0359444e46a9", "qa_cbbfd7d59d01"}, fullQA_schedule)
	})
	t.Run("qa_cc992174acf7", func(t *testing.T) {
		fullExecute(t, cache, "schedule", "qa_cc992174acf7", []string{"qa_cc992174acf7", "qa_fbe62ef8b3e7", "qa_c69b89ef022d", "qa_f47503d90583", "qa_b1336ed19443", "qa_0359444e46a9", "qa_cbbfd7d59d01"}, fullQA_schedule)
	})
	t.Run("qa_cbbfd7d59d01", func(t *testing.T) {
		fullExecute(t, cache, "schedule", "qa_cbbfd7d59d01", []string{"qa_cc992174acf7", "qa_fbe62ef8b3e7", "qa_c69b89ef022d", "qa_f47503d90583", "qa_b1336ed19443", "qa_0359444e46a9", "qa_cbbfd7d59d01"}, fullQA_schedule)
	})
	t.Run("qa_f47503d90583", func(t *testing.T) {
		fullExecute(t, cache, "schedule", "qa_f47503d90583", []string{"qa_cc992174acf7", "qa_fbe62ef8b3e7", "qa_c69b89ef022d", "qa_f47503d90583", "qa_b1336ed19443", "qa_0359444e46a9", "qa_cbbfd7d59d01"}, fullQA_schedule)
	})
	t.Run("qa_0359444e46a9", func(t *testing.T) {
		fullExecute(t, cache, "schedule", "qa_0359444e46a9", []string{"qa_cc992174acf7", "qa_fbe62ef8b3e7", "qa_c69b89ef022d", "qa_f47503d90583", "qa_b1336ed19443", "qa_0359444e46a9", "qa_cbbfd7d59d01"}, fullQA_schedule)
	})
	t.Run("qa_fbe62ef8b3e7", func(t *testing.T) {
		fullExecute(t, cache, "schedule", "qa_fbe62ef8b3e7", []string{"qa_cc992174acf7", "qa_fbe62ef8b3e7", "qa_c69b89ef022d", "qa_f47503d90583", "qa_b1336ed19443", "qa_0359444e46a9", "qa_cbbfd7d59d01"}, fullQA_schedule)
	})
	t.Run("qa_b1336ed19443", func(t *testing.T) {
		fullExecute(t, cache, "schedule", "qa_b1336ed19443", []string{"qa_cc992174acf7", "qa_fbe62ef8b3e7", "qa_c69b89ef022d", "qa_f47503d90583", "qa_b1336ed19443", "qa_0359444e46a9", "qa_cbbfd7d59d01"}, fullQA_schedule)
	})
	t.Run("qa_6f9ed781a410", func(t *testing.T) {
		fullExecute(t, cache, "stats", "qa_6f9ed781a410", []string{"qa_06327bbf436b", "qa_cc281fbd1340", "qa_6f9ed781a410", "qa_0c64e79794cc", "qa_29349c0683f5", "qa_cb8dc4d08772"}, fullQA_stats)
	})
	t.Run("qa_cc281fbd1340", func(t *testing.T) {
		fullExecute(t, cache, "stats", "qa_cc281fbd1340", []string{"qa_06327bbf436b", "qa_cc281fbd1340", "qa_6f9ed781a410", "qa_0c64e79794cc", "qa_29349c0683f5", "qa_cb8dc4d08772"}, fullQA_stats)
	})
	t.Run("qa_06327bbf436b", func(t *testing.T) {
		fullExecute(t, cache, "stats", "qa_06327bbf436b", []string{"qa_06327bbf436b", "qa_cc281fbd1340", "qa_6f9ed781a410", "qa_0c64e79794cc", "qa_29349c0683f5", "qa_cb8dc4d08772"}, fullQA_stats)
	})
	t.Run("qa_0c64e79794cc", func(t *testing.T) {
		fullExecute(t, cache, "stats", "qa_0c64e79794cc", []string{"qa_06327bbf436b", "qa_cc281fbd1340", "qa_6f9ed781a410", "qa_0c64e79794cc", "qa_29349c0683f5", "qa_cb8dc4d08772"}, fullQA_stats)
	})
	t.Run("qa_29349c0683f5", func(t *testing.T) {
		fullExecute(t, cache, "stats", "qa_29349c0683f5", []string{"qa_06327bbf436b", "qa_cc281fbd1340", "qa_6f9ed781a410", "qa_0c64e79794cc", "qa_29349c0683f5", "qa_cb8dc4d08772"}, fullQA_stats)
	})
	t.Run("qa_ff6fe44cc46d", func(t *testing.T) {
		fullExecute(t, cache, "devices", "qa_ff6fe44cc46d", []string{"qa_48e10c8c3688", "qa_4228ae60b590", "qa_f76445835b7d", "qa_3a16670b310f", "qa_418e812084ee", "qa_00cbd3d800bd", "qa_a925a3248bcd", "qa_d5581d7d5601", "qa_5f8b9b9d7f78", "qa_83ba432b89af", "qa_ff6fe44cc46d", "qa_edfec6f68242", "qa_42caca08cb39"}, fullQA_devices)
	})
	t.Run("qa_00cbd3d800bd", func(t *testing.T) {
		fullExecute(t, cache, "devices", "qa_00cbd3d800bd", []string{"qa_48e10c8c3688", "qa_4228ae60b590", "qa_f76445835b7d", "qa_3a16670b310f", "qa_418e812084ee", "qa_00cbd3d800bd", "qa_a925a3248bcd", "qa_d5581d7d5601", "qa_5f8b9b9d7f78", "qa_83ba432b89af", "qa_ff6fe44cc46d", "qa_edfec6f68242", "qa_42caca08cb39"}, fullQA_devices)
	})
	t.Run("qa_d5581d7d5601", func(t *testing.T) {
		fullExecute(t, cache, "devices", "qa_d5581d7d5601", []string{"qa_48e10c8c3688", "qa_4228ae60b590", "qa_f76445835b7d", "qa_3a16670b310f", "qa_418e812084ee", "qa_00cbd3d800bd", "qa_a925a3248bcd", "qa_d5581d7d5601", "qa_5f8b9b9d7f78", "qa_83ba432b89af", "qa_ff6fe44cc46d", "qa_edfec6f68242", "qa_42caca08cb39"}, fullQA_devices)
	})
	t.Run("qa_5f8b9b9d7f78", func(t *testing.T) {
		fullExecute(t, cache, "devices", "qa_5f8b9b9d7f78", []string{"qa_48e10c8c3688", "qa_4228ae60b590", "qa_f76445835b7d", "qa_3a16670b310f", "qa_418e812084ee", "qa_00cbd3d800bd", "qa_a925a3248bcd", "qa_d5581d7d5601", "qa_5f8b9b9d7f78", "qa_83ba432b89af", "qa_ff6fe44cc46d", "qa_edfec6f68242", "qa_42caca08cb39"}, fullQA_devices)
	})
	t.Run("qa_a925a3248bcd", func(t *testing.T) {
		fullExecute(t, cache, "devices", "qa_a925a3248bcd", []string{"qa_48e10c8c3688", "qa_4228ae60b590", "qa_f76445835b7d", "qa_3a16670b310f", "qa_418e812084ee", "qa_00cbd3d800bd", "qa_a925a3248bcd", "qa_d5581d7d5601", "qa_5f8b9b9d7f78", "qa_83ba432b89af", "qa_ff6fe44cc46d", "qa_edfec6f68242", "qa_42caca08cb39"}, fullQA_devices)
	})
	t.Run("qa_83ba432b89af", func(t *testing.T) {
		fullExecute(t, cache, "devices", "qa_83ba432b89af", []string{"qa_48e10c8c3688", "qa_4228ae60b590", "qa_f76445835b7d", "qa_3a16670b310f", "qa_418e812084ee", "qa_00cbd3d800bd", "qa_a925a3248bcd", "qa_d5581d7d5601", "qa_5f8b9b9d7f78", "qa_83ba432b89af", "qa_ff6fe44cc46d", "qa_edfec6f68242", "qa_42caca08cb39"}, fullQA_devices)
	})
	t.Run("qa_113173a66889", func(t *testing.T) {
		fullExecute(t, cache, "token", "qa_113173a66889", []string{"qa_113173a66889"}, fullQA_token)
	})
	t.Run("qa_99b351df6c47", func(t *testing.T) {
		fullExecute(t, cache, "auth", "qa_99b351df6c47", []string{"qa_861607f5c661", "qa_99b351df6c47"}, fullQA_auth)
	})
	t.Run("qa_861607f5c661", func(t *testing.T) {
		fullExecute(t, cache, "auth", "qa_861607f5c661", []string{"qa_861607f5c661", "qa_99b351df6c47"}, fullQA_auth)
	})
	t.Run("qa_bd03c2bacf73", func(t *testing.T) {
		fullExecute(t, cache, "voice", "qa_bd03c2bacf73", []string{"qa_bd03c2bacf73", "qa_a4be39a4d66b", "qa_c7377d6722fc", "qa_95007431d632"}, fullQA_voice)
	})
	t.Run("qa_95007431d632", func(t *testing.T) {
		fullExecute(t, cache, "voice", "qa_95007431d632", []string{"qa_bd03c2bacf73", "qa_a4be39a4d66b", "qa_c7377d6722fc", "qa_95007431d632"}, fullQA_voice)
	})
	t.Run("qa_a4be39a4d66b", func(t *testing.T) {
		fullExecute(t, cache, "voice", "qa_a4be39a4d66b", []string{"qa_bd03c2bacf73", "qa_a4be39a4d66b", "qa_c7377d6722fc", "qa_95007431d632"}, fullQA_voice)
	})
	t.Run("qa_c7377d6722fc", func(t *testing.T) {
		fullExecute(t, cache, "voice", "qa_c7377d6722fc", []string{"qa_bd03c2bacf73", "qa_a4be39a4d66b", "qa_c7377d6722fc", "qa_95007431d632"}, fullQA_voice)
	})
}
