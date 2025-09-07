# HealthCheckOutputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**Ok** | **bool** |  | 

## Methods

### NewHealthCheckOutputBody

`func NewHealthCheckOutputBody(ok bool, ) *HealthCheckOutputBody`

NewHealthCheckOutputBody instantiates a new HealthCheckOutputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewHealthCheckOutputBodyWithDefaults

`func NewHealthCheckOutputBodyWithDefaults() *HealthCheckOutputBody`

NewHealthCheckOutputBodyWithDefaults instantiates a new HealthCheckOutputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *HealthCheckOutputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *HealthCheckOutputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *HealthCheckOutputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *HealthCheckOutputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetOk

`func (o *HealthCheckOutputBody) GetOk() bool`

GetOk returns the Ok field if non-nil, zero value otherwise.

### GetOkOk

`func (o *HealthCheckOutputBody) GetOkOk() (*bool, bool)`

GetOkOk returns a tuple with the Ok field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOk

`func (o *HealthCheckOutputBody) SetOk(v bool)`

SetOk sets Ok field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


