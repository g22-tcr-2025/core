package config

import "time"

const (
	MsgLogin         = "LOGIN"
	MsgLoginResponse = "LOGIN_RESPONSE"
	MsgMatchFound    = "MATCH_FOUND"
	MsgMatchStart    = "MATCH_START"
	MsgStateUpdate   = "STATE_UPDATE"
	MsgMatchEnd      = "MATCH_END"

	MatchDuration = 3 * time.Minute
)
