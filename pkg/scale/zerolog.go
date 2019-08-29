package scale

import "github.com/rs/zerolog"

func (g *GroupReq) MarshalZerologObject(e *zerolog.Event) {
	e.Str("direction", g.Direction.String()).Int("count", g.Count).Str("group", g.GroupName)
}
