package bootstrap

import (
	"github.com/P0H1ng/Cardinal/internal/game"
	"github.com/P0H1ng/Cardinal/internal/timer"
)

func GameToTimerBridge() {
	timer.SetRankListTitle = game.SetRankListTitle
	timer.SetRankList = game.SetRankList
	timer.CleanGameBoxStatus = game.CleanGameBoxStatus
	timer.GetLatestScoreRound = game.GetLatestScoreRound
	timer.RefreshFlag = game.RefreshFlag
	timer.CalculateRoundScore = game.CalculateRoundScore
}
