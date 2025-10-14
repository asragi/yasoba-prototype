package skill

import (
	"github.com/asragi/yasoba-prototype/actor"
	gameskill "github.com/asragi/yasoba-prototype/game/skill"
)

// SelectedAction describes the actor, partner, and targets of a skill.
type SelectedAction struct {
	Id       gameskill.SkillId
	Actor    actor.ActorId
	SubActor actor.ActorId
	Target   []actor.ActorId
}

// SkillApplyResultRow stores the result for a single target.
type SkillApplyResultRow struct {
	ActorId        actor.ActorId
	TargetId       actor.ActorId
	TargetSide     actor.ActorSide
	SkillId        gameskill.SkillId
	Damage         Damage
	IsTargetBeaten bool
	AfterHp        actor.HP
}

// SkillApplyResult aggregates all row results for a single skill cast.
type SkillApplyResult struct {
	Actor    actor.ActorId
	SubActor actor.ActorId
	SkillId  gameskill.SkillId
	Rows     []*SkillApplyResultRow
}
