# AddExternalUserInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**ExternalUserHandle** | Pointer to **string** |  | [optional] 
**ExternalUserId** | Pointer to **string** |  | [optional] 
**ExternalUserSource** | **string** |  | 

## Methods

### NewAddExternalUserInputBody

`func NewAddExternalUserInputBody(externalUserSource string, ) *AddExternalUserInputBody`

NewAddExternalUserInputBody instantiates a new AddExternalUserInputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAddExternalUserInputBodyWithDefaults

`func NewAddExternalUserInputBodyWithDefaults() *AddExternalUserInputBody`

NewAddExternalUserInputBodyWithDefaults instantiates a new AddExternalUserInputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *AddExternalUserInputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *AddExternalUserInputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *AddExternalUserInputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *AddExternalUserInputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetExternalUserHandle

`func (o *AddExternalUserInputBody) GetExternalUserHandle() string`

GetExternalUserHandle returns the ExternalUserHandle field if non-nil, zero value otherwise.

### GetExternalUserHandleOk

`func (o *AddExternalUserInputBody) GetExternalUserHandleOk() (*string, bool)`

GetExternalUserHandleOk returns a tuple with the ExternalUserHandle field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalUserHandle

`func (o *AddExternalUserInputBody) SetExternalUserHandle(v string)`

SetExternalUserHandle sets ExternalUserHandle field to given value.

### HasExternalUserHandle

`func (o *AddExternalUserInputBody) HasExternalUserHandle() bool`

HasExternalUserHandle returns a boolean if a field has been set.

### GetExternalUserId

`func (o *AddExternalUserInputBody) GetExternalUserId() string`

GetExternalUserId returns the ExternalUserId field if non-nil, zero value otherwise.

### GetExternalUserIdOk

`func (o *AddExternalUserInputBody) GetExternalUserIdOk() (*string, bool)`

GetExternalUserIdOk returns a tuple with the ExternalUserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalUserId

`func (o *AddExternalUserInputBody) SetExternalUserId(v string)`

SetExternalUserId sets ExternalUserId field to given value.

### HasExternalUserId

`func (o *AddExternalUserInputBody) HasExternalUserId() bool`

HasExternalUserId returns a boolean if a field has been set.

### GetExternalUserSource

`func (o *AddExternalUserInputBody) GetExternalUserSource() string`

GetExternalUserSource returns the ExternalUserSource field if non-nil, zero value otherwise.

### GetExternalUserSourceOk

`func (o *AddExternalUserInputBody) GetExternalUserSourceOk() (*string, bool)`

GetExternalUserSourceOk returns a tuple with the ExternalUserSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalUserSource

`func (o *AddExternalUserInputBody) SetExternalUserSource(v string)`

SetExternalUserSource sets ExternalUserSource field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


