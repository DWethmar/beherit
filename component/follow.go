package component

const FollowType Type = "follow"

type Follow struct {
	X, Y int
}

func NewFollowComponent() *Component {
	return &Component{
		Type: FollowType,
		Data: &Follow{},
	}
}
