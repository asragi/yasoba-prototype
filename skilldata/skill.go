package skilldata

// SkillId identifies a battle skill.
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

// SkillFunctionId classifies skill execution logic.
type SkillFunctionId string

const (
	SkillFunctionIdNormal      SkillFunctionId = "normal"
	SkillFunctionIdCombination SkillFunctionId = "combination"
)

// SkillType represents damage calculation mode.
type SkillType int

const (
	SkillTypePhysical SkillType = iota
	SkillTypeMagical
)

// SkillPower is the base strength of a skill row.
type SkillPower float64

// AttackTargetType describes actual target application.
type AttackTargetType int

// SkillSelectTargetType describes target selection behaviour.
type SkillSelectTargetType int

const (
	SkillTargetTypeNone SkillSelectTargetType = iota
	SkillTargetTypeSingleOther
)

// ServeSkillData retrieves skill data by ID.
type ServeSkillData func(id SkillId) *SkillData

// NewSkillServer creates a simple in-memory skill repository.
func NewSkillServer() ServeSkillData {
	const (
		firePower             = 6.0
		thunderPower          = 7.5
		kickPower             = 1.0
		combinationEfficiency = 1.5
	)
	dict := map[SkillId]*SkillData{}
	register := func(id SkillId, targetType SkillSelectTargetType, funcId SkillFunctionId, rows []*SkillDataDetail) {
		dict[id] = &SkillData{
			SkillId:         id,
			SkillFunctionId: funcId,
			TargetType:      targetType,
			Rows:            rows,
		}
	}
	register(
		SkillIdLuneAttack, SkillTargetTypeSingleOther, SkillFunctionIdNormal, []*SkillDataDetail{
			{Power: 1.0, Type: SkillTypePhysical},
		},
	)
	register(
		SkillIdLuneFireEnemy, SkillTargetTypeSingleOther, SkillFunctionIdNormal, []*SkillDataDetail{
			{Power: firePower, Type: SkillTypeMagical},
		},
	)
	register(
		SkillIdNormalTackle, SkillTargetTypeSingleOther, SkillFunctionIdNormal, []*SkillDataDetail{
			{Power: 1.0, Type: SkillTypePhysical},
		},
	)
	register(
		SkillIdCombinationThunder, SkillTargetTypeNone, SkillFunctionIdCombination, []*SkillDataDetail{
			{
				Power:    thunderPower * combinationEfficiency,
				Type:     SkillTypeMagical,
				SubPower: kickPower * combinationEfficiency,
				SubType:  SkillTypePhysical,
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

// SkillDataDetail represents a single row in a skill definition.
type SkillDataDetail struct {
	Power            SkillPower
	Type             SkillType
	SubPower         SkillPower
	SubType          SkillType
	AttackTargetType AttackTargetType
	SkillFunctionId  SkillFunctionId
}

// SkillData binds metadata with execution details.
type SkillData struct {
	SkillId         SkillId
	SkillFunctionId SkillFunctionId
	TargetType      SkillSelectTargetType
	Rows            []*SkillDataDetail
}

