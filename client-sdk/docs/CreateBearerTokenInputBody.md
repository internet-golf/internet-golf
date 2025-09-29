# CreateBearerTokenInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**FullPermissions** | **bool** |  | 

## Methods

### NewCreateBearerTokenInputBody

`func NewCreateBearerTokenInputBody(fullPermissions bool, ) *CreateBearerTokenInputBody`

NewCreateBearerTokenInputBody instantiates a new CreateBearerTokenInputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateBearerTokenInputBodyWithDefaults

`func NewCreateBearerTokenInputBodyWithDefaults() *CreateBearerTokenInputBody`

NewCreateBearerTokenInputBodyWithDefaults instantiates a new CreateBearerTokenInputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *CreateBearerTokenInputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *CreateBearerTokenInputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *CreateBearerTokenInputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *CreateBearerTokenInputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetFullPermissions

`func (o *CreateBearerTokenInputBody) GetFullPermissions() bool`

GetFullPermissions returns the FullPermissions field if non-nil, zero value otherwise.

### GetFullPermissionsOk

`func (o *CreateBearerTokenInputBody) GetFullPermissionsOk() (*bool, bool)`

GetFullPermissionsOk returns a tuple with the FullPermissions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFullPermissions

`func (o *CreateBearerTokenInputBody) SetFullPermissions(v bool)`

SetFullPermissions sets FullPermissions field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


