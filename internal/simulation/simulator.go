package simulation

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/services"
	"bitbucket.org/electronicjaw/asgardparty-server/utils"
	"math"
)

type SimulatorService struct {
	factory          *services.CoreService
	progressListener func(percent float64)
}

func NewSimulatorService() *SimulatorService {
	return &SimulatorService{factory: services.NewCoreService()}
}

func (s SimulatorService) WithProgressListener(progressListener func(percent float64)) *SimulatorService {
	s.progressListener = progressListener

	return &s
}

func (s *SimulatorService) Simulate(game string, count int, wager int64) *Result {
	res := &Result{
		Wager: wager,
		Count: count,
		Game:  game,
		Spent: wager * int64(count),
	}

	var (
		factory = s.factory
		percent = count / 100
	)

	user := s.factory.BaseUser()
	user.State.WagerAmount = wager

	// base, bonus, total
	awardSample := make([][3]float64, count)

	for i := 0; i < count; i++ {
		spin := factory.RollBaseSpin(user)

		baseAward, bonusAward := spin.BaseAward(), spin.BonusAward()
		award := spin.Amount

		awardSample[i] = [3]float64{float64(baseAward), float64(bonusAward), float64(award)}

		res.BaseAward += baseAward
		res.BonusAward += bonusAward
		res.Award += award
		res.CascadeBaseAward += spin.CascadeBaseAward()
		res.FirstStageBaseAward += spin.FirstStageBaseAward()

		res.BaseAwardSquareSum += baseAward * baseAward
		res.BonusAwardSquareSum += bonusAward * bonusAward
		res.AwardSquareSum += award * award

		if s.progressListener != nil && percent != 0 && i%percent == 0 {
			s.progressListener(float64(i) / float64(count))
		}

		res.BonusGameCount += spin.BonusGameCount()
		res.BonusGameWithWinCount += spin.BonusGameWithWinCount()

		if spin.Amount > 0 {
			res.WinSpinCount++
		}

		if award/wager >= 1 {
			res.CountX1++
		}

		if award/wager >= 10 {
			res.CountX10++
		}

		if award/wager >= 100 {
			res.CountX100++
		}
	}

	avgBaseAward := float64(res.BaseAward) / float64(count)
	avgBonusAward := float64(res.BonusAward) / float64(count)
	avgAward := float64(res.Award) / float64(count)

	for i := 0; i < count; i++ {
		res.BaseAwardStandardDeviation += (awardSample[i][0] - avgBaseAward) * (awardSample[i][0] - avgBaseAward)
		res.BonusAwardStandardDeviation += (awardSample[i][1] - avgBonusAward) * (awardSample[i][1] - avgBonusAward)
		res.AwardStandardDeviation += (awardSample[i][2] - avgAward) * (awardSample[i][2] - avgAward)
	}

	res.BaseAwardStandardDeviation = math.Sqrt(res.BaseAwardStandardDeviation / float64(count-1))
	res.BonusAwardStandardDeviation = math.Sqrt(res.BonusAwardStandardDeviation / float64(count-1))
	res.AwardStandardDeviation = math.Sqrt(res.AwardStandardDeviation / float64(count-1))

	res.BaseAwardStandardDeviation = utils.SetPrecision(res.BaseAwardStandardDeviation, 4)
	res.BonusAwardStandardDeviation = utils.SetPrecision(res.BonusAwardStandardDeviation, 4)
	res.AwardStandardDeviation = utils.SetPrecision(res.AwardStandardDeviation, 4)

	RTP := float64(res.Award) / float64(res.Spent) * 100
	RTPBaseGame := float64(res.BaseAward) / float64(res.Spent) * 100
	RTPBonusGame := float64(res.BonusAward) / float64(res.Spent) * 100
	RTPCascadeBase := float64(res.CascadeBaseAward) / float64(res.Spent) * 100
	RTPFirstStageBase := float64(res.FirstStageBaseAward) / float64(res.Spent) * 100

	res.RTP = utils.PrecisionString(RTP) + "%"
	res.RTPBaseGame = utils.PrecisionString(RTPBaseGame) + "%"
	res.RTPBonusGame = utils.PrecisionString(RTPBonusGame) + "%"
	res.RTPCascadeBase = utils.PrecisionString(RTPCascadeBase) + "%"
	res.RTPFirstStageBase = utils.PrecisionString(RTPFirstStageBase) + "%"

	res.RateX1 = utils.PrecisionString(float64(res.CountX1)/float64(res.Count)*100) + "%"
	res.RateX10 = utils.PrecisionString(float64(res.CountX10)/float64(res.Count)*100) + "%"
	res.RateX100 = utils.PrecisionString(float64(res.CountX100)/float64(res.Count)*100) + "%"

	res.WinSpinRate = utils.PrecisionString(float64(res.WinSpinCount)/float64(count)*100.0) + "%"
	res.BonusGameRate = float64(count) / math.Max(float64(res.BonusGameCount), 1)

	return res
}
