package core

import (
	"fmt"
	"github.com/asragi/yasoba-prototype/util"
	"math"
	"strconv"
)

type SkillId string

const (
	SkillIdLuneAttack       SkillId = "lune-attack"
	SkillIdLuneFireEnemy    SkillId = "lune-fire-enemy"
	SkillIdLuneFireAlly     SkillId = "lune-fire-ally"
	SkillIdLuneThunderEnemy SkillId = "lune-thunder-enemy"
)

type SkillType int

const (
	SkillTypePhysical SkillType = iota
	SkillTypeMagical
)

type SkillPower float64
type SkillTargetType int

type ServeSkillData func(id SkillId) *SkillData

func NewSkillServer() ServeSkillData {
	dict := map[SkillId]*SkillData{}
	register := func(id SkillId, rows []SkillDataRow) {
		dict[id] = &SkillData{
			Id:   id,
			Rows: rows,
		}
	}
	register(
		SkillIdLuneAttack, []SkillDataRow{
			&SkillSingleAttackRow{
				Power: 1.0,
				Type:  SkillTypePhysical,
			},
		},
	)
	register(
		SkillIdLuneFireEnemy, []SkillDataRow{
			&SkillSingleAttackRow{
				Power: 5,
				Type:  SkillTypeMagical,
			},
		},
	)
	return func(id SkillId) *SkillData {
		return dict[id]
	}
}

type SkillDataRow interface{}

type SkillSingleAttackRow struct {
	Power SkillPower
	Type  SkillType
}

type SkillSingleAttackResult struct {
	ActorId        ActorId
	TargetId       ActorId
	SkillId        SkillId
	Damage         Damage
	IsTargetBeaten bool
	AfterHp        HP
}

type SkillData struct {
	Id   SkillId
	Rows []SkillDataRow
}

type SelectedAction struct {
	Id       SkillId
	Actor    ActorId
	SubActor ActorId
	Target   []ActorId
}

type SkillApplyResultRow interface{}
type SkillApplyResult struct {
	Rows []SkillApplyResultRow
}

type SkillCalculationFunc func(*SelectedAction) *SkillApplyResult

func CreateSkillApply(
	skillServer ServeSkillData,
	supplyActor ActorSupplier,
	updateActor UpdateActorFunc,
	random util.EmitRandomFunc,
) SkillCalculationFunc {
	return func(args *SelectedAction) *SkillApplyResult {
		result := make([]SkillApplyResultRow, 0)
		data := skillServer(args.Id)
		for _, row := range data.Rows {
			switch r := row.(type) {
			case *SkillSingleAttackRow:
				actor := supplyActor(args.Actor)
				target := supplyActor(args.Target[0])
				damage := calculateNormalAttackDamage(
					actor.ATK,
					actor.MAG,
					target.DEF,
					target.MAG,
					r.Power,
					r.Type,
					random,
				)
				afterHP := damage.Apply(target.HP)
				target.HP = afterHP
				fmt.Printf("damage: %d, afterHP: %d\n", damage, afterHP)
				updateActor(target)
				result = append(
					result, &SkillSingleAttackResult{
						ActorId:        args.Actor,
						TargetId:       args.Target[0],
						SkillId:        args.Id,
						Damage:         damage,
						IsTargetBeaten: afterHP <= 0,
						AfterHp:        afterHP,
					},
				)
			}
		}
		return &SkillApplyResult{Rows: result}
	}
}

type Damage int

func (d Damage) String() string {
	return strconv.Itoa(int(d))
}

func (d Damage) Apply(hp HP) HP {
	return HP(int(hp) - int(d))
}

func calculateNormalAttackDamage(
	attackerATK ATK,
	attackerMAG MAG,
	defenderDEF DEF,
	defenderMAG MAG,
	power SkillPower,
	attackType SkillType,
	random util.EmitRandomFunc,
) Damage {
	// 同パラメータのキャラクターがPower=1.0の攻撃を受けた場合に失うHPの割合を表す係数
	const damageRate = 0.3
	// うまいことやって出た数字
	const baseValue = 7.0
	attackValue := func() float64 {
		if attackType == SkillTypePhysical {
			return float64(attackerATK)
		}
		return float64(attackerMAG)
	}()
	defenderValue := func() float64 {
		if attackType == SkillTypePhysical {
			return float64(defenderDEF)
		}
		return float64(defenderMAG)
	}()
	attackPower := float64(power) * baseValue * (math.Pow(attackValue, 3) / (math.Pow(attackValue, 2) + 1600))
	defencePower := (109 - defenderValue) / 109
	randomValue := random() * attackValue
	return Damage(attackPower*defencePower + randomValue)
}
