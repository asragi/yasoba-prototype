package core

import (
	"fmt"
	"github.com/asragi/yasoba-prototype/util"
	"math"
	"strconv"
)

type SkillId string

const (
	SkillIdLuneAttack         SkillId = "lune-attack"
	SkillIdLuneFireEnemy      SkillId = "lune-fire-enemy"
	SkillIdLuneFireAlly       SkillId = "lune-fire-ally"
	SkillIdLuneThunderEnemy   SkillId = "lune-thunder-enemy"
	SkillIdSunnyKick          SkillId = "sunny-attack"
	SkillIdSunnyUppercut      SkillId = "sunny-uppercut"
	SkillIdCombinationFire    SkillId = "combination-fire"
	SkillIdCombinationThunder SkillId = "combination-thunder"
	SkillIdNormalTackle       SkillId = "normal-tackle"
)

type SkillType int

const (
	SkillTypePhysical SkillType = iota
	SkillTypeMagical
)

type SkillPower float64
type SkillTargetType int

const (
	SkillTargetTypeNone SkillTargetType = iota
	SkillTargetTypeSingleOther
)

type ServeSkillData func(id SkillId) *SkillData

func NewSkillServer() ServeSkillData {
	const (
		FirePower             = 6.0
		ThunderPower          = 7.5
		KickPower             = 1.0
		CombinationEfficiency = 1.5
	)
	dict := map[SkillId]*SkillData{}
	register := func(id SkillId, targetType SkillTargetType, rows []*SkillDataDetail) {
		dict[id] = &SkillData{
			SkillId:    id,
			TargetType: targetType,
			Rows:       rows,
		}
	}
	register(
		SkillIdLuneAttack, SkillTargetTypeSingleOther, []*SkillDataDetail{
			{
				Power:     1.0,
				Type:      SkillTypePhysical,
				ActorType: ActorTypeMain,
			},
		},
	)
	register(
		SkillIdLuneFireEnemy, SkillTargetTypeSingleOther, []*SkillDataDetail{
			{
				Power:     FirePower,
				Type:      SkillTypeMagical,
				ActorType: ActorTypeMain,
			},
		},
	)
	register(
		SkillIdNormalTackle, SkillTargetTypeSingleOther, []*SkillDataDetail{
			{
				Power:     1.0,
				Type:      SkillTypePhysical,
				ActorType: ActorTypeMain,
			},
		},
	)
	register(
		SkillIdCombinationThunder, SkillTargetTypeNone, []*SkillDataDetail{
			{
				Power:     ThunderPower * CombinationEfficiency,
				Type:      SkillTypeMagical,
				ActorType: ActorTypeMain,
			},
			{
				Power:     KickPower * CombinationEfficiency,
				Type:      SkillTypePhysical,
				ActorType: ActorTypeSub,
			},
		},
	)
	return func(id SkillId) *SkillData {
		skill, ok := dict[id]
		if !ok {
			panic("skill not found")
		}
		return skill
	}
}

type ActorType int

const (
	ActorTypeMain ActorType = iota
	ActorTypeSub
)

type SkillDataDetail struct {
	Power     SkillPower
	Type      SkillType
	ActorType ActorType
}

type SkillData struct {
	SkillId    SkillId
	TargetType SkillTargetType
	Rows       []*SkillDataDetail
}

type SkillDataRow interface{}

type SkillSingleAttackRow struct {
	Power SkillPower
	Type  SkillType
}

type SkillSingleAttackResult struct {
	ActorId        ActorId
	TargetId       ActorId
	TargetSide     ActorSide
	SkillId        SkillId
	Damage         Damage
	IsTargetBeaten bool
	AfterHp        HP
}

type SkillSingleCombinationAttackResult struct {
	ActorId        ActorId
	SubActorId     ActorId
	TargetId       ActorId
	TargetSide     ActorSide
	SkillId        SkillId
	Damage         Damage
	IsTargetBeaten bool
	AfterHp        HP
}

// deprecated: use SkillData
type SkillDataOld struct {
	Id         SkillId
	TargetType SkillTargetType
	Rows       []SkillDataRow
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
	return func(args *SelectedAction) *SkillApplyResult {
		result := make([]*SkillApplyResultRow, 0)
		data := skillServer(args.Id)
		for _, row := range data.Rows {
			actualActorId := func() ActorId {
				if row.ActorType == ActorTypeMain {
					return args.Actor
				}
				return args.SubActor
			}()
			actualActor := supplyActor(actualActorId)
			target := supplyActor(args.Target[0])
			damage := calculateNormalAttackDamage(
				actualActor.ATK,
				actualActor.MAG,
				target.DEF,
				target.MAG,
				row.Power,
				row.Type,
				random,
			)
			afterHP := damage.Apply(target.HP)
			target.HP = afterHP
			fmt.Printf("actualActor: %s, target: %s\n", actualActor.Id, target.Id)
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
		return &SkillApplyResult{
			Actor:   args.Actor,
			SkillId: args.Id,
			Rows:    result,
		}
	}
}

type Damage int

func (d Damage) String() string {
	return strconv.Itoa(int(d))
}

func (d Damage) Apply(hp HP) HP {
	return HP(math.Max(0, float64(hp)-float64(d)))
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
