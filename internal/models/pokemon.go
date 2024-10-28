package models

type Pokemon struct {
	Identifier     string `csv:"identifier"`
	Id             uint   `csv:"id"`
	SpeciesId      uint   `csv:"species_id"`
	Height         uint   `csv:"height"`
	Weight         uint   `csv:"weight"`
	BaseExperience uint   `csv:"base_experience"`
	Order          uint   `csv:"order"`
	IsDefault      bool   `csv:"is_default"`
}

// Implement the Model interface
func (pokemon *Pokemon) GetID() uint {
	return pokemon.Id
}
