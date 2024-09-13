package core

import "math"

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
				Power: 0.6,
				Type:  SkillTypePhysical,
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
	ActorId  ActorId
	TargetId ActorId
	SkillId  SkillId
	Damage   Damage
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

type SkillApplyFunc func(*SelectedAction) *SkillApplyResult

func CreateSkillApply(
	skillServer ServeSkillData,
	supplyActor ActorSupplier,
) SkillApplyFunc {
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
				)
				result = append(
					result, &SkillSingleAttackResult{
						ActorId:  args.Actor,
						TargetId: args.Target[0],
						SkillId:  args.Id,
						Damage:   damage,
					},
				)
			}
		}
		return &SkillApplyResult{Rows: result}
	}
}

type SkillResultRow interface{}

type SkillResult struct {
	Rows []SkillResultRow
}

type Damage int

func calculateNormalAttackDamage(
	attackerATK ATK,
	attackerMAG MAG,
	defenderDEF DEF,
	defenderMAG MAG,
	power SkillPower,
	attackType SkillType,
) Damage {
	// 同パラメータのキャラクターがPower=1.0の攻撃を受けた場合に失うHPの割合を表す係数
	const damageRate = 0.3
	// うまいことやって出た数字
	const baseValue = 4.0
	attackValue := func() int {
		if attackType == SkillTypePhysical {
			return int(attackerATK)
		}
		return int(attackerMAG)
	}()
	defenderValue := func() int {
		if attackType == SkillTypePhysical {
			return int(defenderDEF)
		}
		return int(defenderMAG)
	}()
	estimatedHP := attackValue * 10
	ratio := float64(attackValue) / float64(defenderValue)
	return Damage(int(float64(estimatedHP) * damageRate * float64(power) * math.Pow(baseValue, ratio-1)))
}
