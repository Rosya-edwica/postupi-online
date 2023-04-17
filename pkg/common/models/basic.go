package models

type Basic struct {
	// BudgetPlaces - string, а не int, т.к не получается вырезать юникод в строке
	Name        string
	Url         string
	Description string
	Direction   string
	Image       string
	Logo        string
	Cost        int
	Scores      Scores
}

type Scores struct {
	PointsBudget  float64
	PointsPayment float64
	PlacesBudget  int
	PlacesPayment int
}
