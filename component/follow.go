package component

const FollowType Type = "follow"

type Follow struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewFollowComponent() *Component {
	return &Component{
		Type: FollowType,
		Data: &Follow{},
	}
}
