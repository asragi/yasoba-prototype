package core

type ActorId string

const (
	ActorLuneId  = ActorId("lune")
	ActorSunnyId = ActorId("sunny")
)

type Actor struct {
	Id          ActorId
	CharacterId CharacterId
	MaxHP       MaxHP
	HP          HP
	ATK         ATK
	MAG         MAG
	DEF         DEF
	SPD         SPD
}

type ChosenAction struct {
	Actor   ActorId
	Target  []ActorId
	SkillId SkillId
}

type ServeActorCharacterFunc func(CharacterId) *Actor

func CreateServeActorCharacter(
	serveCharacter ServeCharacterFunc,
) ServeActorCharacterFunc {
	return func(id CharacterId) *Actor {
		character := serveCharacter(id)
		return &Actor{
			CharacterId: id,
			MaxHP:       character.MaxHP,
			HP:          character.HP,
			ATK:         character.ATK,
			MAG:         character.MAG,
			DEF:         character.DEF,
			SPD:         character.SPD,
		}
	}
}

type ActorServer interface {
	Get(ActorId) *Actor
	Upsert(*Actor)
}

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

	}
}
