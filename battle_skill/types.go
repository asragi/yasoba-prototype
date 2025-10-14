package battle_skill

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/skilldata"
)

// SelectedAction describes the actor, partner, and targets of a skill.
type SelectedAction struct {
	Id       skilldata.SkillId
	Actor    actor.ActorId
	SubActor actor.ActorId
	Target   []actor.ActorId
}

// SkillApplyResultRow stores the result for a single target.
type SkillApplyResultRow struct {
	ActorId        actor.ActorId
	TargetId       actor.ActorId
	TargetSide     actor.ActorSide
	SkillId        skilldata.SkillId
	Damage         Damage
	IsTargetBeaten bool
	AfterHp        actor.HP
}

// SkillApplyResult aggregates all row results for a single skill cast.
type SkillApplyResult struct {
	Actor    actor.ActorId
	SubActor actor.ActorId
	SkillId  skilldata.SkillId
	Rows     []*SkillApplyResultRow
}
