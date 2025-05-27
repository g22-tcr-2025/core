package config

import "time"

const (
	MsgLogin         = "LOGIN"
	MsgLoginResponse = "LOGIN_RESPONSE"

	MsgMatchStart = "MATCH_START"

	MsgUpdateMnana = "UPDATE_MANA"

	MsgAttack       = "ATTACK"
	MsgAttackResult = "ATTACK_RS"

	MsgError = "ERROR"

	MatchDuration = 5 * time.Minute
)
