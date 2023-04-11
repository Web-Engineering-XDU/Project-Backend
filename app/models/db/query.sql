-- name: GetAgentList :many
SELECT 
	agent.id id,
	agent.name name,
	agent.enable enable,
	agent.type_id type_id,
	agent_type.name type_name,
	agent.event_max_age event_max_age,
	agent.prop_json_str prop_json_str,
	agent.create_time create_time
FROM (
	agent INNER JOIN agent_type
	ON agent.type_id = agent_type.id
)
WHERE agent.deleted = 0 AND
	agent_type.deleted = 0;