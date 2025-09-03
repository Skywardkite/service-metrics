package repository

type Gauge struct {
	Name 	string 	`db:"name"`
	Value 	float64 `db:"value"`
}

type Counter struct {
	Name 	string 	`db:"name"`
	Value 	int64 	`db:"value"`
}