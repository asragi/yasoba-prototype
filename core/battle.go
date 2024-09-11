package core

type ActionId string

type ActionCommand struct {
	ActionId ActionId
	ActorId  ActorId
	TargetId []ActorId
	SkillId  SkillId
}

type ActionResult struct {
	ActionId           ActionId
	CombinationActorId ActorId
	IsCombination      bool
	ActorId            ActorId
}

type ExecuteActionFunc func(*ActionCommand) *ActionResult

func ExecuteAction(
	actorServer ActorServer,
) ExecuteActionFunc {
	return func(command *ActionCommand) *ActionResult {
		actor := actorServer.Get(command.ActorId)
		return &ActionResult{ActorId: actor.Id}
	}
}
