package models

type Profile struct {
	Age                int                `json:"age" bson:"age"`
	Education          []Education        `json:"education" bson:"education"`
	Perks              []Perk             `json:"perks" bson:"perks"`
	SelfRepresentation SelfRepresentation `json:"description" bson:"description"`
	Pronounce          string             `json:"pronounce" bson:"pronounce"`
}

type Education struct {
	Organization string `json:"organization" bson:"organization"`
	Level        string `json:"level" bson:"level"`
	Accomplished bool   `json:"is_accomplished" bson:"is_accomplished"`
}

type Perk struct {
	Name        string `json:"name_of_perk" bson:"name_of_perk"`
	Description string `json:"description_of_perk" bson:"description_of_perk"`
}

type SelfRepresentation struct {
	Text string `json:"text" bson:"text"`
}
