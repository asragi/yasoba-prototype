package core

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

type SkillFunctionId string

const (
	SkillFunctionIdNormal      SkillFunctionId = "normal"
	SkillFunctionIdCombination SkillFunctionId = "combination"
)

// SkillType represents which function will be used to calculate the damage.
type SkillType int

const (
	SkillTypePhysical SkillType = iota
	SkillTypeMagical
)

type SkillPower float64

// AttackTargetType is a type that represents how actually targets are selected.
type AttackTargetType int

// SkillSelectTargetType is a type that represents how to select the target of the skill.
type SkillSelectTargetType int

const (
	SkillTargetTypeNone SkillSelectTargetType = iota
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
	register := func(id SkillId, targetType SkillSelectTargetType, rows []*SkillDataDetail) {
		dict[id] = &SkillData{
			SkillId:    id,
			TargetType: targetType,
			Rows:       rows,
		}
	}
	register(
		SkillIdLuneAttack, SkillTargetTypeSingleOther, []*SkillDataDetail{
			{
				Power: 1.0,
				Type:  SkillTypePhysical,
			},
		},
	)
	register(
		SkillIdLuneFireEnemy, SkillTargetTypeSingleOther, []*SkillDataDetail{
			{
				Power: FirePower,
				Type:  SkillTypeMagical,
			},
		},
	)
	register(
		SkillIdNormalTackle, SkillTargetTypeSingleOther, []*SkillDataDetail{
			{
				Power: 1.0,
				Type:  SkillTypePhysical,
			},
		},
	)
	register(
		SkillIdCombinationThunder, SkillTargetTypeNone, []*SkillDataDetail{
			{
				Power:    ThunderPower * CombinationEfficiency,
				Type:     SkillTypeMagical,
				SubPower: KickPower * CombinationEfficiency,
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

type SkillDataDetail struct {
	Power            SkillPower
	Type             SkillType
	SubPower         SkillPower
	SubType          SkillType
	AttackTargetType AttackTargetType
	SkillFunctionId  SkillFunctionId
}

type SkillData struct {
	SkillId    SkillId
	TargetType SkillSelectTargetType
	Rows       []*SkillDataDetail
}
