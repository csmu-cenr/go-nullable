package nullable

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type address struct {
	AddressLine1 string           `json:"addressLine" xml:"AddressLine"`
	AddressLine2 Nullable[string] `json:"addressLine2,omitempty" xml:"AddressLine2"`
	PostNumber   string           `json:"postNumber" xml:"PostNumber"`
	PostCity     string           `json:"city" xml:"City"`
	County       Nullable[string] `json:"county,omitempty" xml:"County"`
}

type person struct {
	FirstName     string            `json:"firstName" xml:"first-name,attr"`
	LastName      Nullable[string]  `json:"lastName,omitempty" xml:"last-name,attr,omitempty"`
	PostAddress   address           `json:"postAddress" xml:"PostAddress"`
	OfficeAddress Nullable[address] `json:"officeAddress,omitempty" xml:"officeAddress,omitempty"`
}

func Test_Json_unmarshal_chained_struct_null_value(t *testing.T) {
	var john person
	err := json.Unmarshal([]byte(`{
		"firstName": "John",
		"lastName": "Smith",
		"postAddress": {
			"addressLine": "RoadStreet 1A",
			"postNumber": "1234",
			"county": "Texas"
		}
	}`), &john)
	assert.NoError(t, err)
	assert.Equal(t, "John", john.FirstName)
	assert.True(t, john.LastName.Valid)
	assert.Equal(t, "Smith", john.LastName.Data)
	assert.False(t, john.OfficeAddress.Valid)
	assert.Equal(t, "RoadStreet 1A", john.PostAddress.AddressLine1)
}

func Test_Json_unmarshal_chained_struct_with_value(t *testing.T) {
	var john person
	err := json.Unmarshal([]byte(`{
		"firstName": "John",
		"lastName": "Smith",
		"postAddress": {
			"addressLine": "RoadStreet 1A",
			"postNumber": "1234",
			"county": "Texas"
		},
		"officeAddress": {
			"addressLine": "RoadStreet 1B",
			"postNumber": "1234",
			"county": "Texas"
		}
	}`), &john)
	assert.NoError(t, err)
	assert.Equal(t, "John", john.FirstName)
	assert.True(t, john.LastName.Valid)
	assert.Equal(t, "Smith", john.LastName.Data)
	assert.Equal(t, "RoadStreet 1A", john.PostAddress.AddressLine1)
	assert.True(t, john.OfficeAddress.Valid)
	assert.Equal(t, "RoadStreet 1B", john.OfficeAddress.Data.AddressLine1)
	assert.False(t, john.OfficeAddress.Data.AddressLine2.Valid)
}

func Test_Json_marshal_chain_with_null_value(t *testing.T) {
	john := person{
		FirstName: "John",
		LastName:  Value("Smith"),
		PostAddress: address{
			AddressLine1: "RoadStreet 1A",
			PostNumber:   "1234",
			PostCity:     "Texas",
		},
		OfficeAddress: Null[address](),
	}
	assert.False(t, john.PostAddress.AddressLine2.Valid)
	assert.False(t, john.OfficeAddress.Valid)

	jsonData, err := json.Marshal(john)
	assert.NoError(t, err)
	assert.JSONEq(t, `{ "firstName": "John", "lastName": "Smith", "postAddress": { "addressLine": "RoadStreet 1A", "postNumber": "1234", "city": "Texas", "addressLine2": null, "county": null }, "officeAddress": null }`, string(jsonData), "chain struct marshal")
}
