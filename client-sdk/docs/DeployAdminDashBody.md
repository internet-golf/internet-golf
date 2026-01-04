# DeployAdminDashBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**Url** | **string** | The URL that you want to deploy the admin dashboard to. | 

## Methods

### NewDeployAdminDashBody

`func NewDeployAdminDashBody(url string, ) *DeployAdminDashBody`

NewDeployAdminDashBody instantiates a new DeployAdminDashBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeployAdminDashBodyWithDefaults

`func NewDeployAdminDashBodyWithDefaults() *DeployAdminDashBody`

NewDeployAdminDashBodyWithDefaults instantiates a new DeployAdminDashBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *DeployAdminDashBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *DeployAdminDashBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *DeployAdminDashBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *DeployAdminDashBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetUrl

`func (o *DeployAdminDashBody) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *DeployAdminDashBody) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *DeployAdminDashBody) SetUrl(v string)`

SetUrl sets Url field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


