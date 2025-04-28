package models

// Ciclo representa un ciclo académico en el sistema
type Ciclo struct {
	ID      string `json:"id_"`
	IDCiclo string `json:"id_ciclos"`
	Ciclo   string `json:"ciclo"`
	Version int    `json:"version"`
}
