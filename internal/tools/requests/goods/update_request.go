package goods

type Attributes struct {
	Amount int32 `json:"amount" validate:"required,ne=0"`
}

type Data struct {
	Type       string     `json:"type"       validate:"required"`
	ID         string     `json:"id"         validate:"required"`
	Attributes Attributes `json:"attributes" validate:"required"`
}

type Request struct {
	Data []Data `json:"data" validate:"required"`
}
