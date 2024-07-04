package avro

type SimpleRecord struct {
	StringField  string  `avro:"StringField" json:"StringField"`
	FloatField   float32 `avro:"FloatField" json:"FloatField"`
	BooleanField bool    `avro:"BooleanField" json:"BooleanField"`
}

type InvalidSimpleRecord struct {
	StringField  bool    `avro:"StringField" json:"StringField"`
	FloatField   string  `avro:"FloatField" json:"FloatField"`
	BooleanField float32 `avro:"BooleanField" json:"BooleanField"`
}
