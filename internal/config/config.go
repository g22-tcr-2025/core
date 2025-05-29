package config

import "time"

const (
	MsgLogin         = "LOGIN"
	MsgLoginResponse = "LOGIN_RESPONSE"

	MsgTick        = "TICK"
	MsgMatchStart  = "MATCH_START"
	MsgMatchUpdate = "MATCH_UPDATE"
	MsgMatchEnd    = "MATCH_END"

	MsgUpdatePlayerMnana  = "UPDATE_PLAYER_MANA"
	MsgUpdateOpponentMana = "UPDATE_OPPONENT_MANA"

	MsgAttack       = "ATTACK"
	MsgAttackResult = "ATTACK_RS"

	MsgError = "ERROR"

	MatchDuration = 3 * time.Minute
)
