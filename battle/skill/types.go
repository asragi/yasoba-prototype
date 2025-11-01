package skill

import (
	"github.com/asragi/yasoba-prototype/battle/actor"
	"github.com/asragi/yasoba-prototype/common/character"
)

// SelectedAction describes the actor, partner, and targets of a skill.
type SelectedAction struct {
	Id       SkillId
	Actor    actor.ActorId
	SubActor actor.ActorId
	Target   []actor.ActorId
}

// SkillApplyResultRow stores the result for a single target.
type SkillApplyResultRow struct {
	ActorId        actor.ActorId
	TargetId       actor.ActorId
	TargetSide     actor.ActorSide
	SkillId        SkillId
	Damage         Damage
	IsTargetBeaten bool
	AfterHp        character.HP
}

// SkillApplyResult aggregates all row results for a single skill cast.
type SkillApplyResult struct {
	Actor    actor.ActorId
	SubActor actor.ActorId
	SkillId  SkillId
	Rows     []*SkillApplyResultRow
}
