-- name: GetAgentRuntimeInfoList :many
SELECT 
	agent.id id,
	agent.enable enable,
	agent.event_forever event_forever,
	agent.event_max_age event_max_age,
	agent.prop_json_str prop_json_str,
	agent.type_id type_id,
	agent_type.allow_input allow_input,
	agent_type.allow_output allow_output
FROM (
	agent INNER JOIN agent_type
	ON agent.type_id = agent_type.id
)
WHERE
	agent.deleted = 0 AND
	agent_type.deleted = 0;

-- name: GetAgentBasicInfoList :many
SELECT 
	agent.id id,
	agent.name name,
	agent.enable enable,
	agent.type_id type_id,
	agent_type.name type_name
FROM (
	agent INNER JOIN agent_type
	ON agent.type_id = agent_type.id
)
WHERE agent.deleted = 0 AND
	agent_type.deleted = 0
ORDER BY agent.create_at DESC
LIMIT ?, ?;

-- name: GetAgentDetail :one
SELECT 
	agent.id id,
	agent.name name,
	agent.enable enable,
	agent.event_forever event_forever,
	agent.type_id type_id,
	agent_type.name type_name,
	agent.event_max_age event_max_age,
	agent.prop_json_str prop_json_str,
	agent.create_at create_at,
	agent.description description
FROM (
	agent INNER JOIN agent_type
	ON agent.type_id = agent_type.id
)
WHERE
	agent.deleted = 0 AND
	agent_type.deleted = 0 AND
	agent.id = ?
LIMIT 1;

-- name: AddAgent :one
INSERT INTO
	agent(
		name,
		enable,
		type_id,
		event_forever,
		event_max_age,
		prop_json_str,
		create_at,
		description
	)
VALUES(?, ?, ?, ?, ?, ?, ?, ?)
-- VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: SoftDeleteAgent :exec
UPDATE agent
SET deleted = 1
WHERE
	id = ? AND
	deleted = 0;

-- name: GetAgentTypes :many
SELECT
	id,
	name,
	allow_input,
	allow_output
FROM agent_type
WHERE
	deleted = 0;

-- name: GetAgentType :one
SELECT
	id,
	name,
	allow_input,
	allow_output
FROM agent_type
WHERE
	deleted = 0 AND
	id = ?;

-- name: GetEventList :many
SELECT
	event.id id,
	event.src_agent_id src_agent_id,
	agent_type.name src_agent_name,
	event.json_str json_str,
	event.error error,
	event.create_at create_at
FROM(
	event INNER JOIN agent_type
	ON event.src_agent_id = agent_type.id
)
WHERE event.deleted = 0
ORDER BY event.create_at DESC
LIMIT ? OFFSET ?;

-- name: GetEventDetail :one
SELECT
	event.id id,
	event.src_agent_id src_agent_id,
	agent_type.name src_agent_name,
	event.json_str json_str,
	event.error error,
	event.log log,
	event.create_at create_at
FROM(
	event INNER JOIN agent_type
	ON event.src_agent_id = agent_type.id
)
WHERE
	event.id = ? AND
	event.deleted = 0
LIMIT 1;

-- name: AddEvent :exec
INSERT INTO
	event (
		src_agent_id,
		json_str,
		content_hash,
		error,
		log,
		create_at,
		delete_at
	)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: SoftDeleteEventById :exec
UPDATE event
SET deleted = 1
WHERE 
	id = ? AND
	deleted = 0;

-- name: SoftDeleteAllEventAbout :exec
UPDATE event
SET deleted = 1
WHERE
	src_agent_id = ? AND
	deleted = 0;
	
-- name: GetAgentRelationList :many
SELECT
	src_agent_id,
	dst_agent_id
FROM agent_relation;

-- name: AddAgentRelation :exec
INSERT INTO
	agent_relation (
		src_agent_id,
		dst_agent_id
	)
VALUES (?, ?);

-- name: DeleteAgentRelation :exec
DELETE FROM
	agent_relation
WHERE
	src_agent_id = ? AND
	dst_agent_id = ?;
	
-- name: DeleteAllAgentRelationAbout :exec
DELETE FROM
	agent_relation
WHERE
	src_agent_id = ? OR
	dst_agent_id = ?;