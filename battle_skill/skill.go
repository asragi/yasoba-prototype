package battle_skill

import (
	"math"

	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/util"
)

func decideAttackValue(atk actor.ATK, mag actor.MAG, skillType core.SkillType) attackValue {
	if skillType == core.SkillTypePhysical {
		return attackValue(atk)
	}
	return attackValue(mag)
}

type attackValue float64

// SkillApplyFunc calculates and applies skill effects to targets.
type SkillApplyFunc func(*SelectedAction) *SkillApplyResult

// CreateSkillApply creates a SkillApplyFunc given repositories and RNG.
func CreateSkillApply(
	skillServer core.ServeSkillData,
	supplyActor actor.ActorSupplier,
	updateActor actor.UpdateActorFunc,
	random util.EmitRandomFunc,
) SkillApplyFunc {
	applyAttack := func(
		args *SelectedAction,
		decidePower func(*core.SkillDataDetail) attackPower,
		decideAttack func(*core.SkillDataDetail) attackValue,
	) *SkillApplyResult {
		result := make([]*SkillApplyResultRow, 0)
		data := skillServer(args.Id)
		for _, row := range data.Rows {
			for _, targetId := range args.Target {
				target := supplyActor(targetId)
				attack := decideAttack(row)
				power := decidePower(row)
				randomDamageValue := newRandomDamage(attack, random)
				damage := calculateNormalAttackDamage(
					power,
					target.DEF,
					target.MAG,
					row.Type,
					randomDamageValue,
				)
				afterHP := damage.Apply(target.HP)
				target.HP = afterHP
				updateActor(target)
				result = append(
					result, &SkillApplyResultRow{
						ActorId:        args.Actor,
						TargetId:       args.Target[0],
						TargetSide:     target.Side,
						SkillId:        args.Id,
						Damage:         damage,
						IsTargetBeaten: afterHP <= 0,
						AfterHp:        afterHP,
					},
				)
			}
		}
		return &SkillApplyResult{
			Actor:   args.Actor,
			SkillId: args.Id,
			Rows:    result,
		}
	}
	normalAttack := func(args *SelectedAction) *SkillApplyResult {
		actorId := args.Actor
		mainActor := supplyActor(actorId)
		return applyAttack(
			args,
			func(row *core.SkillDataDetail) attackPower {
				attack := decideAttackValue(mainActor.ATK, mainActor.MAG, row.Type)
				return calculateNormalAttackPower(attack, row.Power)
			},
			func(row *core.SkillDataDetail) attackValue {
				return decideAttackValue(mainActor.ATK, mainActor.MAG, row.Type)
			},
		)
	}
	combinationAttack := func(args *SelectedAction) *SkillApplyResult {
		mainActor := supplyActor(args.Actor)
		subActor := supplyActor(args.SubActor)
		return applyAttack(
			args,
			func(row *core.SkillDataDetail) attackPower {
				return calculateCombinationAttackPower(
					mainActor.ATK,
					mainActor.MAG,
					subActor.ATK,
					subActor.MAG,
					row.Power,
					row.SubPower,
					row.Type,
					row.SubType,
				)
			},
			func(row *core.SkillDataDetail) attackValue {
				return decideAttackValue(mainActor.ATK, mainActor.MAG, row.Type)
			},
		)
	}

	return func(args *SelectedAction) *SkillApplyResult {
		skill := skillServer(args.Id)
		if skill.SkillFunctionId == core.SkillFunctionIdCombination {
			return combinationAttack(args)
		}
		return normalAttack(args)
	}
}

// randomDamage is a partial damage value that is calculated by random value.
type randomDamage int

func newRandomDamage(attack attackValue, emitRandom util.EmitRandomFunc) randomDamage {
	randomValue := emitRandom()
	return randomDamage(float64(attack) * randomValue)
}

const attackPowerBase = 7.0

func calculateNormalAttackPower(
	attackValue attackValue,
	power core.SkillPower,
) attackPower {
	return attackPower(
		float64(power) * attackPowerBase *
			(math.Pow(float64(attackValue), 3) /
				(math.Pow(float64(attackValue), 2) + 1600)),
	)
}

func calculateCombinationAttackPower(
	attackerATK actor.ATK,
	attackerMAG actor.MAG,
	subAttackerATK actor.ATK,
	subAttackerMAG actor.MAG,
	power core.SkillPower,
	subPower core.SkillPower,
	attackType core.SkillType,
	subAttackType core.SkillType,
) attackPower {
	mainAttackValue := decideAttackValue(attackerATK, attackerMAG, attackType)
	subAttackValue := decideAttackValue(subAttackerATK, subAttackerMAG, subAttackType)
	mainAttackPower := calculateNormalAttackPower(mainAttackValue, power)
	subAttackPower := calculateNormalAttackPower(subAttackValue, subPower)
	return mainAttackPower + subAttackPower
}

// attackPower is calculated based on attackValue.
type attackPower float64

func calculateNormalAttackDamage(
	attackPower attackPower,
	defenderDEF actor.DEF,
	defenderMAG actor.MAG,
	attackType core.SkillType,
	randomDamage randomDamage,
) Damage {
	defenderValue := func() float64 {
		if attackType == core.SkillTypePhysical {
			return float64(defenderDEF)
		}
		return float64(defenderMAG)
	}()
	defencePower := (109 - defenderValue) / 109
	return Damage(float64(attackPower)*defencePower + float64(randomDamage) + 1)
}
