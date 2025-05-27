package config

import "time"

const (
	MsgLogin         = "LOGIN"
	MsgLoginResponse = "LOGIN_RESPONSE"

	MsgMatchStart = "MATCH_START"

	MsgUpdatePlayerMnana  = "UPDATE_PLAYER_MANA"
	MsgUpdateOpponentMana = "UPDATE_OPPONENT_MANA"

	MsgAttack       = "ATTACK"
	MsgAttackResult = "ATTACK_RS"

	MsgError = "ERROR"

	MatchDuration = 5 * time.Minute
)
