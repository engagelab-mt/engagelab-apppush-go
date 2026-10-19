//go:build integration

package engagelab

// Native compiled workflows; no external case interpreter.
func (c *fullQA) device(params ...any) any {
	var rid any
	_ = rid
	rid = qarg(params, 0)
	return c.rest("GET", qcat("/v4/devices/", qtxt(rid)), nil)
}

func (c *fullQA) relation(params ...any) any {
	var rid any
	_ = rid
	var tag any
	_ = tag
	tag = qarg(params, 0)
	rid = qarg(params, 1)
	return c.rest("GET", qcat(qcat(qcat("/v4/tags/", qtxt(tag)), "/registration_ids/"), qtxt(rid)), nil)
}

func (c *fullQA) seedDevice(params ...any) any {
	var r any
	_ = r
	var rid any
	_ = rid
	c.write()
	c.require(qeq(c.target(), "qa"), "Backend-only fixtures require QA")
	r = c.rest("POST", "/v4/devices/token/registration_id", qm("platform", "android", "tokens", qarr(qcat(qcat(c.prefix, "_"), c.unique()))))
	rid = qg(r, "data", "results", 0, "registration_id")
	c.require(qboth(qeq(qg(r, "status"), 200), qnonempty(rid), qyes(qg(r, "data", "results", 0, "is_new"))), "Cannot prepare exclusively owned backend device")
	c.own("device", rid)
	return rid
}

func (c *fullQA) seedTag(params ...any) any {
	var r any
	_ = r
	var rid any
	_ = rid
	var tag any
	_ = tag
	tag = qarg(params, 0)
	rid = qarg(params, 1)
	c.own("tag", tag)
	r = c.rest("POST", qcat("/v4/tags/", qtxt(tag)), qm("registration_ids", qm("add", qarr(rid))))
	c.require(qeq(qg(r, "status"), 200), "Isolated tag prerequisite rejected")
	return r
}

func (c *fullQA) seedAlias(params ...any) any {
	var alias any
	_ = alias
	var r any
	_ = r
	var rid any
	_ = rid
	alias = qarg(params, 0)
	rid = qarg(params, 1)
	c.rememberAlias(rid)
	c.own("alias", alias)
	r = c.rest("POST", qcat("/v4/devices/", qtxt(rid)), qm("alias", alias))
	c.require(qeq(qg(r, "status"), 200), "Isolated alias prerequisite rejected")
	return r
}

func (c *fullQA) listPlans(params ...any) any {
	var count any
	_ = count
	var page any
	_ = page
	var prefix any
	_ = prefix
	var q any
	_ = q
	var source any
	_ = source
	prefix = qarg(params, 0)
	count = qarg(params, 1)
	page = qarg(params, 2)
	source = qarg(params, 3)
	q = qm("page_index", page, "page_size", count, "search_description", prefix)
	if qyes(qnonempty(source)) {
		qput(q, "send_source", source)
	}
	return c.rest("GET", qcat("/v4/push_plan/list?", qquery(q)), nil)
}

func (c *fullQA) seedPlan(params ...any) any {
	var description any
	_ = description
	var pid any
	_ = pid
	var r any
	_ = r
	pid = qarg(params, 0)
	description = qarg(params, 1)
	c.own("plan", pid)
	r = c.sdk("plan_create", qarr(qm("plan_id", pid, "plan_description", description)))
	c.require(qboth(qyes(qg(r, "ok")), qeq(qg(r, "data", "plan_id"), pid)), "Cannot prepare owned plan through native SDK")
	return r
}

func (c *fullQA) seedPlans(params ...any) any {
	var count any
	_ = count
	var first any
	_ = first
	var i any
	_ = i
	var r any
	_ = r
	count = qarg(params, 0)
	first = nil
	for _, item1 := range qlistOf(qspan(count)) {
		i = item1
		r = c.seedPlan(qcat(qcat(c.prefix, "_"), qtxt(i)), c.prefix)
		if qyes(qeq(i, 0)) {
			first = r
		}
	}
	c.waitFor("plans", qm("prefix", c.prefix, "count", count))
	return first
}

func (c *fullQA) planIDs(params ...any) any {
	return c.ownedValues("plan", 101)
}

func (c *fullQA) batchPlanPair(params ...any) any {
	var exclude any
	_ = exclude
	exclude = qarg(params, 0)
	return qtake(qwithout(c.planIDs(), exclude), 2)
}

func (c *fullQA) hasKeys(params ...any) any {
	var keys any
	_ = keys
	var value any
	_ = value
	value = qarg(params, 0)
	keys = qarg(params, 1)
	return qhasKeysDeep(value, keys)
}

func (c *fullQA) join(params ...any) any {
	var values any
	_ = values
	values = qarg(params, 0)
	return qjoined(values)
}

func (c *fullQA) noneIn(params ...any) any {
	var left any
	_ = left
	var right any
	_ = right
	left = qarg(params, 0)
	right = qarg(params, 1)
	return qno(qintersects(left, right))
}

func (c *fullQA) errorsEqual(params ...any) any {
	var raw any
	_ = raw
	var sdk any
	_ = sdk
	sdk = qarg(params, 0)
	raw = qarg(params, 1)
	return qboth(qno(qg(sdk, "ok")), qeq(qg(sdk, "status"), qg(raw, "status")), qeq(qg(sdk, "error_code"), qg(raw, "data", "error", "code")), qeq(qg(sdk, "error_message"), qg(raw, "data", "error", "message")))
}

func (c *fullQA) planBoundary(params ...any) any {
	var candidate any
	_ = candidate
	var length any
	_ = length
	var payload any
	_ = payload
	var pid any
	_ = pid
	var r any
	_ = r
	var raw any
	_ = raw
	length = qarg(params, 0)
	pid = nil
	if qyes(qeq(length, 1)) {
		for _, item2 := range qlistOf(qarr("a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9")) {
			candidate = item2
			if qyes(qno(qnonempty(pid))) {
				r = c.listPlans(candidate, 10000, 1, nil)
				if qyes(qboth(qge(qsize(qg(r, "data", "push_plan_info")), qnum(qg(r, "data", "total"))), qno(qincludes(qids(qg(r, "data", "push_plan_info"), "plan_id"), candidate)))) {
					pid = candidate
				}
			}
		}
	} else {
		pid = qleft(qcat(qcat(c.prefix, "_"), qrepeat("x", length)), length)
	}
	if qyes(qno(qnonempty(pid))) {
		return qm("ready", false, "verified", false)
	}
	payload = qm("plan_id", pid, "plan_description", c.prefix)
	c.own("plan", pid)
	r = c.sdk("plan_create", qarr(payload))
	if qyes(qeq(length, 51)) {
		raw = c.rest("POST", "/v4/push_plan", payload)
		return qm("ready", true, "verified", qboth(qyes(c.errorsEqual(r, raw)), qeq(qg(r, "error_code"), 27304)))
	}
	raw = c.listPlans(c.prefix, 101, 1, nil)
	return qm("ready", true, "verified", qboth(qyes(qg(r, "ok")), qincludes(qids(qg(raw, "data", "push_plan_info"), "plan_id"), pid), qeq(qsize(pid), length)))
}

func (c *fullQA) targetRID(params ...any) any {
	var index any
	_ = index
	var rid any
	_ = rid
	index = qarg(params, 0)
	rid = c.fixtureRID(index)
	if qyes(qno(qnonempty(rid))) {
		rid = c.seedDevice()
		c.addFixture(rid)
	}
	c.require(qeq(qg(c.device(rid), "status"), 200), "Declared QA target absent from application")
	return rid
}

func (c *fullQA) batchResultsValid(params ...any) any {
	var k any
	_ = k
	var r any
	_ = r
	var targets any
	_ = targets
	var v any
	_ = v
	var valid any
	_ = valid
	var values any
	_ = values
	r = qarg(params, 0)
	targets = qarg(params, 1)
	values = qg(r, "data", "results")
	valid = qsets(qkeys(values), targets)
	for _, item3 := range qlistOf(qkeys(values)) {
		k = item3
		v = qg(values, k)
		valid = qboth(valid, qyes(qg(v, "success")), qnonempty(qg(v, "msg_id")), qeq(qg(v, "target"), k))
	}
	return valid
}

func (c *fullQA) batch(params ...any) any {
	var k any
	_ = k
	var linked any
	_ = linked
	var mids any
	_ = mids
	var op any
	_ = op
	var r any
	_ = r
	var requests any
	_ = requests
	var t any
	_ = t
	var targets any
	_ = targets
	var v any
	_ = v
	op = qarg(params, 0)
	targets = qarg(params, 1)
	requests = qarr()
	for _, item4 := range qlistOf(targets) {
		t = item4
		qappend(requests, qm("target", t, "platform", "android", "notification", qm("alert", c.prefix), "options", qm("time_to_live", 30, "third_party_channel", qm("distribution_new", "mtpush"))))
	}
	r = c.sdk(op, qarr(qm("requests", requests)))
	mids = qarr()
	for _, item5 := range qlistOf(qkeys(qg(r, "data", "results"))) {
		k = item5
		v = qg(r, "data", "results", k)
		if qyes(qyes(qg(v, "success"))) {
			qappend(mids, qtxt(qg(v, "msg_id")))
		}
	}
	linked = false
	if qyes(qeq(qsize(mids), qsize(targets))) {
		linked = c.waitFor("messages", mids)
	}
	qput(r, "messages_persisted", linked)
	return r
}

func (c *fullQA) negative(params ...any) any {
	var base any
	_ = base
	var control any
	_ = control
	var controls any
	_ = controls
	var kind any
	_ = kind
	var name any
	_ = name
	var op any
	_ = op
	var p any
	_ = p
	var path any
	_ = path
	var payload any
	_ = payload
	var r any
	_ = r
	var raw any
	_ = raw
	var types any
	_ = types
	kind = qarg(params, 0)
	base = qm("to", qm("registration_id", qarr("sdk_full_negative_nonexistent")), "body", qm("platform", qarr("ios"), "options", qm("apns_production", false)))
	types = qm("voip", qm("call_id", c.prefix), "message", qm("msg_content", c.prefix), "notification", qm("alert", c.prefix), "live_activity", qm("ios", qm("event", "start", "attributes-type", "SdkFull", "attributes", qm(), "content-state", qm("n", 1))))
	controls = true
	op = "validate"
	path = "/v4/push/validate"
	payload = qm()
	if qyes(qeq(kind, "platform")) {
		payload = qm("to", qm("registration_id", qarr("sdk_full_negative_nonexistent")), "body", qm("platform", "bad-platform", "notification", qm("alert", c.prefix)))
		op = "push_raw"
		path = "/v4/push"
	} else {
		if qyes(qno(qeq(kind, "empty"))) {
			qput(qg(base, "to"), "registration_id", qarr(c.negativeTarget()))
			for _, item6 := range qlistOf(qarr("voip", kind)) {
				name = item6
				p = qclone(base)
				qput(qg(p, "body"), name, qg(types, name))
				control = c.rest("POST", path, p)
				controls = qboth(controls, qeq(qg(control, "status"), 200))
			}
			payload = base
			qput(qg(payload, "body"), "voip", qg(types, "voip"))
			qput(qg(payload, "body"), kind, qg(types, kind))
		}
	}
	raw = c.rest("POST", path, payload)
	r = c.sdk(op, qarr(payload))
	return qm("verified", qboth(qyes(c.errorsEqual(r, raw)), qno(qg(r, "ok"))), "controls_valid", controls)
}

func (c *fullQA) vipValid(params ...any) any {
	var r any
	_ = r
	r = qarg(params, 0)
	return qboth(qisInteger(qg(r, "data", "vip_status")), qisInteger(qg(r, "data", "vip_end_time")))
}

func (c *fullQA) payload(params ...any) any {
	var marker any
	_ = marker
	var options any
	_ = options
	var plan any
	_ = plan
	var rid any
	_ = rid
	rid = qarg(params, 0)
	marker = qarg(params, 1)
	plan = qarg(params, 2)
	options = qm("time_to_live", 30, "third_party_channel", qm("distribution_new", "mtpush"))
	if qyes(qnonempty(plan)) {
		qput(options, "plan_id", plan)
	}
	return qm("request_id", marker, "to", qm("registration_id", qarr(rid)), "body", qm("platform", qarr("android"), "notification", qm("alert", marker, "android", qm("title", marker)), "options", options))
}

func (c *fullQA) scheduleOne(params ...any) any {
	var created any
	_ = created
	var kind any
	_ = kind
	var param any
	_ = param
	var pid any
	_ = pid
	var push any
	_ = push
	var raw any
	_ = raw
	var rid any
	_ = rid
	var sid any
	_ = sid
	var trigger any
	_ = trigger
	var when any
	_ = when
	kind = qarg(params, 0)
	when = qarg(params, 1)
	rid = qarg(params, 2)
	pid = qarg(params, 3)
	trigger = qm("single", qm("time", when, "zone_type", 0))
	if qyes(qeq(kind, "intelligent")) {
		trigger = qm("intelligent", qm("backup_time", when))
	}
	push = c.payload(rid, qcat(qcat(c.prefix, "_"), qtxt(kind)), pid)
	qdrop(push, "request_id")
	param = qm("name", qcat(qcat(c.prefix, "_"), qtxt(kind)), "enabled", true, "trigger", trigger, "push", push)
	created = c.sdk("schedule_create", qarr(param))
	sid = qg(created, "data", "schedule_id")
	raw = qm()
	if qyes(qnonempty(sid)) {
		c.own("schedule", sid)
		raw = c.rest("GET", qcat("/v4/schedules/", qtxt(sid)), nil)
	}
	return qm("sdk", created, "id", sid, "param", param, "read", raw, "detail", qfirstOrSelf(qg(raw, "data")))
}

func (c *fullQA) prepareSchedules(params ...any) any {
	var afterDelete any
	_ = afterDelete
	var afterUpdate any
	_ = afterUpdate
	var deadline any
	_ = deadline
	var delete any
	_ = delete
	var deleted any
	_ = deleted
	var deletedMessages any
	_ = deletedMessages
	var detail any
	_ = detail
	var disabledMessages any
	_ = disabledMessages
	var known any
	_ = known
	var listed any
	_ = listed
	var messages any
	_ = messages
	var pid any
	_ = pid
	var positive any
	_ = positive
	var positiveIDs any
	_ = positiveIDs
	var positiveReady any
	_ = positiveReady
	var records any
	_ = records
	var rid any
	_ = rid
	var sid any
	_ = sid
	var smart any
	_ = smart
	var tomorrow any
	_ = tomorrow
	var update any
	_ = update
	var updateParam any
	_ = updateParam
	var updated any
	_ = updated
	var when any
	_ = when
	rid = c.targetRID(0)
	pid = qcat(c.prefix, "_schedule_plan")
	c.seedPlan(pid, c.prefix)
	deadline = qadd(c.now(), qmaxi(45, c.envNumber("ENGAGELAB_SCHEDULE_DELAY", 45)))
	when = c.dateTime(deadline)
	tomorrow = c.dateTime(qadd(deadline, 86400))
	positive = c.scheduleOne("positive", when, rid, pid)
	update = c.scheduleOne("update", when, rid, pid)
	delete = c.scheduleOne("delete", when, rid, pid)
	smart = c.scheduleOne("intelligent", tomorrow, rid, pid)
	sid = qg(positive, "id")
	detail = c.sdk("schedule_get", qarr(sid))
	listed = c.listOwnedSchedule(sid)
	updateParam = qclone(qg(update, "param"))
	qput(updateParam, "enabled", false)
	qput(updateParam, "name", qcat(qtxt(qg(updateParam, "name")), "_updated"))
	qput(qg(updateParam, "push", "body", "notification"), "alert", qcat(c.prefix, "_updated"))
	updated = c.sdk("schedule_update", qarr(qg(update, "id"), updateParam))
	afterUpdate = c.rest("GET", qcat("/v4/schedules/", qtxt(qg(update, "id"))), nil)
	deleted = c.sdk("schedule_delete", qarr(qg(delete, "id")))
	afterDelete = c.rest("GET", qcat("/v4/schedules/", qtxt(qg(delete, "id"))), nil)
	c.sleepUntil(qadd(deadline, qmaxi(15, c.envNumber("ENGAGELAB_SCHEDULE_OBSERVATION", 15))))
	c.waitFor("schedule", sid)
	messages = c.sdk("schedule_msgs", qarr(sid))
	records = qg(messages, "data", "MsgIds")
	positiveIDs = c.scheduleIDs(records)
	positiveReady = qnonempty(positiveIDs)
	disabledMessages = c.rest("GET", qcat(qcat("/v4/schedules/", qtxt(qg(update, "id"))), "/msg-ids"), nil)
	deletedMessages = c.rest("GET", qcat(qcat("/v4/schedules/", qtxt(qg(delete, "id"))), "/msg-ids"), nil)
	known = qarr(qg(positive, "id"), qg(update, "id"))
	if qyes(qeq(qg(smart, "read", "status"), 200)) {
		qappend(known, qg(smart, "id"))
	}
	return qm("create", qboth(qyes(qg(positive, "sdk", "ok")), qyes(c.scheduleMatches(qg(positive, "detail"), qg(positive, "param")))), "create_ready", true, "get", qboth(qyes(qg(detail, "ok")), qyes(c.scheduleMatches(qfirstOrSelf(qg(detail, "data")), qg(positive, "param")))), "get_ready", true, "list", qboth(qeq(qg(listed, "status"), 200), qincludes(qids(qg(listed, "data", "schedules"), "schedule_id"), sid), qyes(c.scheduleListMatches(qg(listed, "data"), qg(listed, "raw")))), "list_ready", true, "intelligent", qboth(qyes(qg(smart, "sdk", "ok")), qeq(qg(smart, "detail", "trigger", "intelligent", "backup_time"), tomorrow)), "intelligent_ready", qeq(qg(smart, "read", "status"), 200), "messages", qboth(qeq(qg(messages, "status"), 200), positiveReady, qeq(qnum(qg(messages, "data", "count")), qsize(records)), qeq(records, qg(messages, "raw", "MsgIds")), c.waitFor("messages", positiveIDs)), "messages_ready", positiveReady, "update", qboth(qyes(qg(updated, "ok")), qyes(c.scheduleMatches(qfirstOrSelf(qg(afterUpdate, "data")), updateParam)), qeq(qsize(qcoalesce(qg(disabledMessages, "data", "msg_ids"), qg(disabledMessages, "data", "msgIds"))), 0)), "update_ready", positiveReady, "delete", qboth(qyes(qg(deleted, "ok")), qeq(qg(afterDelete, "status"), 404), qeither(qeq(qg(deletedMessages, "status"), 404), qeq(qsize(qcoalesce(qg(deletedMessages, "data", "msg_ids"), qg(deletedMessages, "data", "msgIds"))), 0))), "delete_ready", positiveReady)
}

func (c *fullQA) prepareStatistics(params ...any) any {
	var batch any
	_ = batch
	var batchReady any
	_ = batchReady
	var end any
	_ = end
	var i any
	_ = i
	var life any
	_ = life
	var lifeMID any
	_ = lifeMID
	var lifeRID any
	_ = lifeRID
	var lifeReady any
	_ = lifeReady
	var linked any
	_ = linked
	var linkedIDs any
	_ = linkedIDs
	var message any
	_ = message
	var messageReady any
	_ = messageReady
	var mid any
	_ = mid
	var mids any
	_ = mids
	var p any
	_ = p
	var pids any
	_ = pids
	var plan any
	_ = plan
	var planReady any
	_ = planReady
	var rids any
	_ = rids
	var sent any
	_ = sent
	var start any
	_ = start
	var users any
	_ = users
	var v any
	_ = v
	rids = qarr(c.targetRID(0), c.targetRID(1))
	pids = qarr(qcat(c.prefix, "_stats_a"), qcat(c.prefix, "_stats_b"))
	mids = qarr()
	for _, item7 := range qlistOf(qspan(2)) {
		i = item7
		c.seedPlan(qg(pids, i), c.prefix)
		sent = c.sdk("push", qarr(c.payload(qg(rids, i), qcat(qcat(c.prefix, "_stats_"), qtxt(i)), qg(pids, i))))
		c.require(qboth(qyes(qg(sent, "ok")), qnonempty(qg(sent, "data", "msg_id"))), "Statistics requires accepted owned messages")
		qappend(mids, qtxt(qg(sent, "data", "msg_id")))
	}
	c.waitFor("messages", mids)
	start = c.date(qadd(c.now(), -172800))
	end = c.date(c.now())
	message = c.sdk("status", qarr(mids))
	lifeMID = qcoalesce(c.env("ENGAGELAB_LIFECYCLE_MESSAGE_ID", ""), qg(mids, 0))
	lifeRID = qcoalesce(c.env("ENGAGELAB_LIFECYCLE_REGISTRATION_ID", ""), qg(rids, 0))
	life = c.sdk("lifecycle", qarr(lifeMID, qarr(lifeRID)))
	batch = c.sdk("batch_lifecycle", qarr(mids))
	users = c.sdk("users", qarr("DAY", start, 3))
	plan = c.sdk("plan_status", qarr(pids, start, end))
	linked = c.sdk("plan_msgs", qarr(qg(pids, 0), start, end))
	lifeReady = c.lifecycleReady(qg(life, "raw", lifeRID))
	batchReady = qnonempty(qg(batch, "raw"))
	for _, item8 := range qlistOf(qlistOf(qg(batch, "raw"))) {
		v = item8
		batchReady = qboth(batchReady, qincludes(mids, qtxt(qg(v, "message_id"))), qincludes(rids, qg(v, "registration_id")))
	}
	planReady = true
	for _, item9 := range qlistOf(pids) {
		p = item9
		planReady = qboth(planReady, qgt(qnum(qcoalesce(qg(plan, "raw", p, "targets"), qg(plan, "raw", p, "target"))), 0))
	}
	messageReady = true
	for _, item10 := range qlistOf(mids) {
		mid = item10
		messageReady = qboth(messageReady, qnonempty(qg(message, "raw", mid)))
	}
	linkedIDs = qg(linked, "data", qg(pids, 0), "msg_ids")
	return qm("message", qboth(qeq(qg(message, "status"), 200), messageReady, qsame(qg(message, "data"), qg(message, "raw"))), "message_ready", messageReady, "lifecycle", qboth(qeq(qg(life, "status"), 200), qsame(qg(life, "data"), qg(life, "raw")), c.lifecycleDefaults(qg(life, "data", lifeRID), qg(life, "raw", lifeRID))), "lifecycle_ready", lifeReady, "batch_lifecycle", qboth(qeq(qg(batch, "status"), 200), qsame(qg(batch, "data"), qg(batch, "raw"))), "batch_lifecycle_ready", batchReady, "users", qboth(qeq(qg(users, "status"), 200), qsame(qg(users, "data"), qg(users, "raw"))), "users_ready", qhmosNonzero(qg(users, "raw")), "plan", qboth(qeq(qg(plan, "status"), 200), qsame(qg(plan, "data"), qg(plan, "raw"))), "plan_ready", planReady, "plan_messages", qboth(qeq(qg(linked, "status"), 200), qincludes(linkedIDs, qg(mids, 0)), qsame(qg(linked, "data"), qg(linked, "raw"))), "plan_messages_ready", qnonempty(linkedIDs))
}

func (c *fullQA) invalidAuth(params ...any) any {
	return c.authControl()
}

func (c *fullQA) apnsToken(params ...any) any {
	var again any
	_ = again
	var baseline any
	_ = baseline
	var cleanup any
	_ = cleanup
	var config any
	_ = config
	var created any
	_ = created
	var parsed any
	_ = parsed
	var payload any
	_ = payload
	var r any
	_ = r
	var raw any
	_ = raw
	var rid any
	_ = rid
	var verified any
	_ = verified
	config = c.apnsConfig()
	payload = qm("platform", "ios", "tokens", qarr(qg(config, "token")), "apns_production", qg(config, "production"))
	r = c.sdk("token", qarr(payload))
	raw = qg(r, "raw", "results")
	parsed = qg(r, "data", "results")
	rid = qg(raw, 0, "registration_id")
	created = qyes(qg(raw, 0, "is_new"))
	if qyes(qboth(qnonempty(rid), created)) {
		c.own("device", rid)
	}
	verified = qboth(qeq(qg(r, "status"), 200), qeq(qsize(raw), 1), qeq(qsize(parsed), 1), qnonempty(rid), qeq(qg(parsed, 0, "token"), qg(config, "token")), qeq(qg(parsed, 0, "registration_id"), rid), qeq(qg(parsed, 0, "is_new"), created), qeither(qeq(qg(parsed, 0, "code"), nil), qeq(qg(parsed, 0, "code"), 0)), qeq(qg(r, "request", "apns_production"), qg(config, "production")))
	baseline = c.device(rid)
	again = c.sdk("token", qarr(payload))
	verified = qboth(verified, qeq(qg(baseline, "status"), 200), qeq(qg(again, "data", "results", 0, "registration_id"), rid), qno(qg(again, "data", "results", 0, "is_new")))
	cleanup = false
	if qyes(qboth(qnonempty(rid), created)) {
		c.rest("DELETE", qcat("/v4/devices/", qtxt(rid)), nil)
		cleanup = c.waitFor("device_absent", rid)
	} else {
		if qyes(qnonempty(rid)) {
			cleanup = qeq(c.device(rid), baseline)
		}
	}
	return qm("verified", verified, "cleanup", cleanup)
}

func (c *fullQA) groupPush(params ...any) any {
	var alias any
	_ = alias
	var app any
	_ = app
	var apps any
	_ = apps
	var binding any
	_ = binding
	var config any
	_ = config
	var k any
	_ = k
	var parsed any
	_ = parsed
	var payload any
	_ = payload
	var r any
	_ = r
	var valid any
	_ = valid
	config = c.groupConfig()
	apps = qg(config, "apps")
	alias = qg(config, "alias")
	c.require(qge(qsize(apps), 2), "Two explicit group members required")
	for _, item11 := range qlistOf(apps) {
		app = item11
		binding = c.memberAliasBinding(app, alias)
		c.require(qeq(qsize(qg(binding, "data", "registration_ids")), 1), "Group alias must resolve exclusively to one prepared test target per app")
	}
	payload = qm("request_id", c.prefix, "to", qm("alias", qarr(alias)), "body", qm("platform", "android", "notification", qm("alert", c.prefix), "options", qm("time_to_live", 30, "third_party_channel", qm("distribution_new", "mtpush"))))
	r = c.sdk("group", qarr(payload))
	for _, item12 := range qlistOf(qkeys(qg(r, "raw"))) {
		k = item12
		c.require(qno(qeq(qg(r, "raw", k, "error", "code"), 21038)), "Group member lacks push permission: 21038")
	}
	parsed = qg(r, "data", "successes")
	valid = qboth(qeq(qg(r, "status"), 200), qsets(qkeys(parsed), apps), qeq(qsize(qg(r, "data", "errors")), 0), qnonempty(qg(r, "data", "group_msgid")), qeq(qg(r, "data", "group_msgid"), qg(r, "raw", "group_msgid")))
	for _, item13 := range qlistOf(apps) {
		app = item13
		valid = qboth(valid, qnonempty(qg(parsed, app, "msg_id")), qsame(qg(parsed, app), qg(r, "raw", app)))
	}
	return qm("verified", valid, "explicit_target", qeq(qg(r, "request", "to"), qm("alias", qarr(alias))))
}

func (c *fullQA) imageUpload(params ...any) any {
	var r any
	_ = r
	var resource any
	_ = resource
	resource = c.imageResource()
	r = c.sdk("image", qarr(qm("big_picture_url", qg(resource, "url"))))
	return qm("resource_valid", qg(resource, "valid"), "verified", qboth(qeq(qg(r, "status"), 200), qnonempty(qg(r, "data", "big_picture_id")), qeq(qg(r, "data", "big_picture_id"), qg(r, "raw", "big_picture_id"))), "prerequisite", qeq(qg(r, "status"), 200))
}

func (c *fullQA) prepareVoice(params ...any) any {
	var detail any
	_ = detail
	var fileURL any
	_ = fileURL
	var found any
	_ = found
	var item any
	_ = item
	var language any
	_ = language
	var listed any
	_ = listed
	var path any
	_ = path
	var r any
	_ = r
	var upload any
	_ = upload
	r = c.sdk("voices", qarr())
	c.require(qyes(qg(r, "ok")), "Voice list prerequisite unavailable")
	language = c.env("ENGAGELAB_VOICE_LANGUAGE", "en")
	c.require(qno(qincludes(qids(qg(r, "data"), "language"), language)), "Voice language occupied; existing template preserved")
	c.own("voice", language)
	c.state("voiceLanguage", language)
	path = c.voiceFile()
	upload = c.sdk("voice_create", qarr(language, path))
	fileURL = qg(upload, "data", "file_url")
	detail = c.sdk("voice_get", qarr(language))
	listed = c.sdk("voices", qarr())
	found = false
	for _, item14 := range qlistOf(qlistOf(qg(listed, "data"))) {
		item = item14
		found = qeither(found, qboth(qeq(qg(item, "language"), language), qeq(qg(item, "file_url"), fileURL)))
	}
	return qm("upload_verified", qboth(qeq(qg(upload, "status"), 200), qnonempty(fileURL), qeq(fileURL, qg(upload, "raw", "file_url")), qyes(qg(upload, "multipart_valid"))), "detail_verified", qboth(qeq(qg(detail, "status"), 200), qeq(qg(detail, "data", "file_url"), fileURL)), "list_verified", qboth(qeq(qg(listed, "status"), 200), found))
}

func (c *fullQA) deleteVoice(params ...any) any {
	var detail any
	_ = detail
	var language any
	_ = language
	var listed any
	_ = listed
	var r any
	_ = r
	language = c.state("voiceLanguage", nil)
	r = c.sdk("voice_delete", qarr(language))
	detail = c.rest("GET", qcat("/v4/voices/", qtxt(language)), nil)
	listed = c.sdk("voices", qarr())
	return qm("sdk_ok", qg(r, "ok"), "verified", qboth(qeq(qg(detail, "status"), 404), qno(qincludes(qids(qg(listed, "data"), "language"), language))))
}

func (c *fullQA) negativeTarget(params ...any) any {
	var config any
	_ = config
	var r any
	_ = r
	var rid any
	_ = rid
	c.write()
	rid = c.state("negativeRID", nil)
	if qyes(qno(qnonempty(rid))) {
		config = c.apnsConfig()
		r = c.rest("POST", "/v4/devices/token/registration_id", qm("platform", "ios", "tokens", qarr(qg(config, "token")), "apns_production", qg(config, "production")))
		rid = qg(r, "data", "results", 0, "registration_id")
		c.require(qboth(qeq(qg(r, "status"), 200), qnonempty(rid)), "Cannot establish legal APNs sandbox validation target")
		if qyes(qyes(qg(r, "data", "results", 0, "is_new"))) {
			c.own("device", rid)
		}
		c.state("negativeRID", rid)
	}
	return rid
}

func (c *fullQA) scheduleMatches(params ...any) any {
	var actual any
	_ = actual
	var audience any
	_ = audience
	var body any
	_ = body
	var push any
	_ = push
	var want any
	_ = want
	actual = qarg(params, 0)
	want = qarg(params, 1)
	push = qg(actual, "push")
	audience = qcoalesce(qg(push, "to"), qg(push, "audience"))
	body = qcoalesce(qg(push, "body"), push)
	return qboth(qeq(qg(actual, "name"), qg(want, "name")), qeq(qg(actual, "enabled"), qg(want, "enabled")), qsame(qg(actual, "trigger"), qg(want, "trigger")), qsame(audience, qg(want, "push", "to")), qsame(body, qg(want, "push", "body")))
}

func (c *fullQA) scheduleListMatches(params ...any) any {
	var a any
	_ = a
	var actual any
	_ = actual
	var b any
	_ = b
	var i any
	_ = i
	var k any
	_ = k
	var raw any
	_ = raw
	var valid any
	_ = valid
	actual = qarg(params, 0)
	raw = qarg(params, 1)
	valid = qboth(qeq(qg(actual, "page"), qg(raw, "page")), qeq(qg(actual, "total_count"), qg(raw, "total_count")), qeq(qg(actual, "total_pages"), qg(raw, "total_pages")), qeq(qsize(qg(actual, "schedules")), qsize(qg(raw, "schedules"))))
	for _, item15 := range qlistOf(qspan(qsize(qg(raw, "schedules")))) {
		i = item15
		a = qg(actual, "schedules", i)
		b = qg(raw, "schedules", i)
		for _, item16 := range qlistOf(qarr("schedule_id", "name", "enabled", "trigger", "push")) {
			k = item16
			valid = qboth(valid, qsame(qg(a, k), qg(b, k)))
		}
	}
	return valid
}

func (c *fullQA) listOwnedSchedule(params ...any) any {
	var listed any
	_ = listed
	var pages any
	_ = pages
	var sid any
	_ = sid
	sid = qarg(params, 0)
	listed = c.sdk("schedule_list", qarr(1))
	pages = qnum(qg(listed, "data", "total_pages"))
	if qyes(qboth(qno(qincludes(qids(qg(listed, "data", "schedules"), "schedule_id"), sid)), qgt(pages, 1))) {
		listed = c.sdk("schedule_list", qarr(pages))
		if qyes(qboth(qno(qincludes(qids(qg(listed, "data", "schedules"), "schedule_id"), sid)), qgt(pages, 2))) {
			listed = c.sdk("schedule_list", qarr(qadd(pages, -1)))
		}
	}
	return listed
}

func (c *fullQA) scheduleIDs(params ...any) any {
	var decoded any
	_ = decoded
	var item any
	_ = item
	var mid any
	_ = mid
	var records any
	_ = records
	var result any
	_ = result
	records = qarg(params, 0)
	result = qarr()
	for _, item17 := range qlistOf(qlistOf(records)) {
		item = item17
		decoded = qdecodeJSON(qtxt(item))
		mid = qg(decoded, "msg_id")
		if qyes(qeither(qeq(decoded, nil), qisInteger(decoded))) {
			mid = item
		}
		if qyes(qboth(qnonempty(mid), qno(qeq(qtxt(mid), "0")))) {
			qappend(result, qtxt(mid))
		}
	}
	return result
}

func (c *fullQA) lifecycleReady(params ...any) any {
	var record any
	_ = record
	record = qarg(params, 0)
	return qboth(qnonempty(record), qeither(qnonempty(qg(record, "status")), qboth(qgt(qnum(qg(record, "itime")), 0), qnonempty(qg(record, "channel"))), qboth(qgt(qnum(qg(record, "error_code")), 0), qno(qincludes(qarr(21003, 21061), qg(record, "error_code"))))))
}

func (c *fullQA) lifecycleDefaults(params ...any) any {
	var parsed any
	_ = parsed
	var raw any
	_ = raw
	var valid any
	_ = valid
	parsed = qarg(params, 0)
	raw = qarg(params, 1)
	valid = true
	if qyes(qno(qincludes(qkeys(raw), "error_code"))) {
		valid = qeither(qeq(qg(parsed, "error_code"), nil), qeq(qg(parsed, "error_code"), 0))
	}
	if qyes(qno(qincludes(qkeys(raw), "error_message"))) {
		valid = qboth(valid, qeq(qtxt(qg(parsed, "error_message")), ""))
	}
	return valid
}
