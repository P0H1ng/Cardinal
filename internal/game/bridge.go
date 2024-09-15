package game

import (
	"github.com/P0H1ng/Cardinal/internal/asteroid"
	"github.com/P0H1ng/Cardinal/internal/db"
	"github.com/P0H1ng/Cardinal/internal/dynamic_config"
	"github.com/P0H1ng/Cardinal/internal/timer"
	"github.com/P0H1ng/Cardinal/internal/utils"
)

func AsteroidGreetData() (result asteroid.Greet) {
	var asteroidTeam []asteroid.Team
	var teams []db.Team
	var asteroidChallenge []asteroid.Challenge
	var challenges []db.Challenge
	var gameBoxes []db.GameBox
	db.MySQL.Model(&db.Team{}).Order("score DESC").Find(&teams)
	db.MySQL.Model(&db.Challenge{}).Find(&challenges)
	db.MySQL.Model(&db.GameBox{}).Find(&gameBoxes)
	for rank, team := range teams {
		asteroidTeam = append(asteroidTeam, asteroid.Team{
			Id:    int(team.ID),
			Name:  team.Name,
			Rank:  rank + 1,
			Image: team.Logo,
			Score: int(team.Score),
		})
	}
	// 建立一個 map 來保存 challenge_id 和對應的 visible 值
	visibleMap := make(map[uint]bool)
	// 對每個 GameBox 進行遍歷，處理相同 challenge_id 的 visible 值
	for _, gameBox := range gameBoxes {
		// 如果已經有該 challenge_id 的 visible 值為 true，則保持 true
		if currentVisible, exists := visibleMap[gameBox.ChallengeID]; exists {
			visibleMap[gameBox.ChallengeID] = currentVisible || gameBox.Visible
		} else {
			// 否則將該 GameBox 的 visible 設為該 challenge_id 的 visible
			visibleMap[gameBox.ChallengeID] = gameBox.Visible
		}
	}
	for _, challenge := range challenges {
		// 根據 challenge_id 從 map 中獲取 visible 值
		visible, exists := visibleMap[challenge.ID]
		if !exists {
			visible = false // 如果沒有匹配的 gameBox，默認為 false
		}
	
		asteroidChallenge = append(asteroidChallenge, asteroid.Challenge{
			ChallengeId:      uint(challenge.ID),
			ChallengeName:    challenge.Title,
			ChallengeVisible: visible,
		})
	}

	result.Title = dynamic_config.Get(utils.TITLE_CONF)
	result.Team = asteroidTeam
	result.Challenge = asteroidChallenge
	result.Time = timer.Get().RoundRemainTime
	result.Round = timer.Get().NowRound
	return
}
