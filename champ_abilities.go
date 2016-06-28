package lol

type spellRange int
type spellEffect string
type spellLabel string
type spelleffectBurn string
type spellcooldown float32
type spellcost int

type spellleveltip struct {
	Effect []spellEffect `json:"effect"`
	Label  []spellLabel  `json:"label"`
}

type spellvars struct {
	Link  string    `json:"link"`
	Coeff []float32 `json:"coeff"`
	Key   string    `json:"key"`
}

type spellsData struct {
	Range                []spellRange      `json:"range"`
	Leveltip             spellleveltip     `json:"leveltip"`
	Resource             string            `json:"resource"`
	Maxrank              int               `json:"maxrank"`
	EffectBurn           []spelleffectBurn `json:"effectBurn"`
	Image                ImageDto          `json:"image"`
	Cooldown             []spellcooldown   `json:"cooldown"`
	Cost                 []spellcost       `json:"cost"`
	Vars                 []spellvars       `json:"vars"`
	SanitizedDescription string            `json:"sanitizedDescription"`
	RangeBurn            string            `json:"sanitizedDescription"`
	costType             string            `json:"costType"`
	Effect               []interface{}     `json:"effect"`
	CooldownBurn         string            `json:"cooldownBurn"`
	Description          string            `json:"description"`
	Name                 string            `json:"name"`
	SanitizedTooltip     string            `json:"sanitizedTooltip"`
	Key                  string            `json:"key"`
	CostBurn             string            `json:"costBurn"`
	Tooltip              string            `json:"tooltip"`
}
