package model

type AthleteId int64
type RelationshipState string

const (
	Unset    = RelationshipState("")
	Pending  = RelationshipState("pending")
	Accepted = RelationshipState("accepted")
	Blocked  = RelationshipState("blocked")
)

type Athlete struct {
	Id        AthleteId         `json:"id"`
	FirstName string            `json:"firstname"`
	LastName  string            `json:"lastname"`
	Friend    RelationshipState `json:"friend"`
	Follower  RelationshipState `json:"follower"`
}
