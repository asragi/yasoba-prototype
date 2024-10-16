package core

import (
	"fmt"
	"github.com/asragi/yasoba-prototype/util"
	"math"
	"strconv"
)

func decideAttackValue(atk ATK, mag MAG, skillType SkillType) attackerValue {
	if skillType == SkillTypePhysical {
		return atk.toAttackValue()
	}
	return mag.toAttackValue()
}

type SelectedAction struct {
	Id       SkillId
	Actor    ActorId
	SubActor ActorId
	Target   []ActorId
}

type SkillApplyResultRow struct {
	ActorId        ActorId
	TargetId       ActorId
	TargetSide     ActorSide
	SkillId        SkillId
	Damage         Damage
	IsTargetBeaten bool
	AfterHp        HP
}

type SkillApplyResult struct {
	Actor    ActorId
	SubActor ActorId
	SkillId  SkillId
	Rows     []*SkillApplyResultRow
}

// SkillApplyFunc is a function that calculate skill effect to the target
// and UPDATE actors status.
type SkillApplyFunc func(*SelectedAction) *SkillApplyResult

func CreateSkillApply(
	skillServer ServeSkillData,
	supplyActor ActorSupplier,
	updateActor UpdateActorFunc,
	random util.EmitRandomFunc,
) SkillApplyFunc {
	applyAttack := func(
		args *SelectedAction,
		decidePower func(*SkillDataDetail) attackPower,
		decideAttack func(*SkillDataDetail) attackerValue,
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
				fmt.Printf("damage: %d, afterHP: %d\n", damage, afterHP)
				fmt.Println("---")
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
		actor := supplyActor(actorId)
		return applyAttack(
			args,
			func(row *SkillDataDetail) attackPower {
				attack := decideAttackValue(actor.ATK, actor.MAG, row.Type)
				return calculateNormalAttackPower(attack, row.Power)
			},
			func(row *SkillDataDetail) attackerValue {
				return decideAttackValue(actor.ATK, actor.MAG, row.Type)
			},
		)
	}
	combinationAttack := func(args *SelectedAction) *SkillApplyResult {
		actorId := args.Actor
		subActorId := args.SubActor
		mainActor := supplyActor(actorId)
		subActor := supplyActor(subActorId)
		return applyAttack(
			args,
			func(row *SkillDataDetail) attackPower {
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
			func(row *SkillDataDetail) attackerValue {
				return decideAttackValue(mainActor.ATK, mainActor.MAG, row.Type)
			},
		)
	}

	return func(args *SelectedAction) *SkillApplyResult {
		skill := skillServer(args.Id)
		if skill.SkillFunctionId == SkillFunctionIdCombination {
			return combinationAttack(args)
		}
		return normalAttack(args)
	}
}

type Damage int

func (d Damage) String() string {
	return strconv.Itoa(int(d))
}

func (d Damage) Apply(hp HP) HP {
	return HP(math.Max(0, float64(hp)-float64(d)))
}

// randomDamage is a partial damage value that is calculated by random value.
type randomDamage int

func newRandomDamage(attack attackerValue, emitRandom util.EmitRandomFunc) randomDamage {
	randomValue := emitRandom()
	return randomDamage(float64(attack) * randomValue)
}

const attackPowerBase = 7.0

func calculateNormalAttackPower(
	attackValue attackerValue,
	power SkillPower,
) attackPower {
	return attackPower(
		float64(power) *
			attackPowerBase *
			(math.Pow(float64(attackValue), 3) /
				(math.Pow(float64(attackValue), 2) + 1600)),
	)
}

func calculateCombinationAttackPower(
	attackerATK ATK,
	attackerMAG MAG,
	subAttackerATK ATK,
	subAttackerMAG MAG,
	power SkillPower,
	subPower SkillPower,
	attackType SkillType,
	subAttackType SkillType,
) attackPower {
	mainAttackValue := decideAttackValue(attackerATK, attackerMAG, attackType)
	subAttackValue := decideAttackValue(subAttackerATK, subAttackerMAG, subAttackType)
	mainAttackPower := calculateNormalAttackPower(mainAttackValue, power)
	subAttackPower := calculateNormalAttackPower(subAttackValue, subPower)
	return mainAttackPower + subAttackPower
}

// attackPower is calculated based on attackerValue
type attackPower float64

func calculateNormalAttackDamage(
	attackPower attackPower,
	defenderDEF DEF,
	defenderMAG MAG,
	attackType SkillType,
	randomDamage randomDamage,
) Damage {
	defenderValue := func() float64 {
		if attackType == SkillTypePhysical {
			return float64(defenderDEF)
		}
		return float64(defenderMAG)
	}()
	defencePower := (109 - defenderValue) / 109
	return Damage(float64(attackPower)*defencePower + float64(randomDamage) + 1)
}
