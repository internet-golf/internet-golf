# DeploymentCreateInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**ExternalSource** | Pointer to **string** |  | [optional] 
**ExternalSourceType** | Pointer to **string** |  | [optional] 
**PreserveExternalPath** | Pointer to **bool** |  | [optional] 
**Tags** | Pointer to **[]string** |  | [optional] 
**Url** | [**Url**](Url.md) |  | 

## Methods

### NewDeploymentCreateInputBody

`func NewDeploymentCreateInputBody(url Url, ) *DeploymentCreateInputBody`

NewDeploymentCreateInputBody instantiates a new DeploymentCreateInputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeploymentCreateInputBodyWithDefaults

`func NewDeploymentCreateInputBodyWithDefaults() *DeploymentCreateInputBody`

NewDeploymentCreateInputBodyWithDefaults instantiates a new DeploymentCreateInputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *DeploymentCreateInputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *DeploymentCreateInputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *DeploymentCreateInputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *DeploymentCreateInputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetExternalSource

`func (o *DeploymentCreateInputBody) GetExternalSource() string`

GetExternalSource returns the ExternalSource field if non-nil, zero value otherwise.

### GetExternalSourceOk

`func (o *DeploymentCreateInputBody) GetExternalSourceOk() (*string, bool)`

GetExternalSourceOk returns a tuple with the ExternalSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSource

`func (o *DeploymentCreateInputBody) SetExternalSource(v string)`

SetExternalSource sets ExternalSource field to given value.

### HasExternalSource

`func (o *DeploymentCreateInputBody) HasExternalSource() bool`

HasExternalSource returns a boolean if a field has been set.

### GetExternalSourceType

`func (o *DeploymentCreateInputBody) GetExternalSourceType() string`

GetExternalSourceType returns the ExternalSourceType field if non-nil, zero value otherwise.

### GetExternalSourceTypeOk

`func (o *DeploymentCreateInputBody) GetExternalSourceTypeOk() (*string, bool)`

GetExternalSourceTypeOk returns a tuple with the ExternalSourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSourceType

`func (o *DeploymentCreateInputBody) SetExternalSourceType(v string)`

SetExternalSourceType sets ExternalSourceType field to given value.

### HasExternalSourceType

`func (o *DeploymentCreateInputBody) HasExternalSourceType() bool`

HasExternalSourceType returns a boolean if a field has been set.

### GetPreserveExternalPath

`func (o *DeploymentCreateInputBody) GetPreserveExternalPath() bool`

GetPreserveExternalPath returns the PreserveExternalPath field if non-nil, zero value otherwise.

### GetPreserveExternalPathOk

`func (o *DeploymentCreateInputBody) GetPreserveExternalPathOk() (*bool, bool)`

GetPreserveExternalPathOk returns a tuple with the PreserveExternalPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreserveExternalPath

`func (o *DeploymentCreateInputBody) SetPreserveExternalPath(v bool)`

SetPreserveExternalPath sets PreserveExternalPath field to given value.

### HasPreserveExternalPath

`func (o *DeploymentCreateInputBody) HasPreserveExternalPath() bool`

HasPreserveExternalPath returns a boolean if a field has been set.

### GetTags

`func (o *DeploymentCreateInputBody) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *DeploymentCreateInputBody) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *DeploymentCreateInputBody) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *DeploymentCreateInputBody) HasTags() bool`

HasTags returns a boolean if a field has been set.

### SetTagsNil

`func (o *DeploymentCreateInputBody) SetTagsNil(b bool)`

 SetTagsNil sets the value for Tags to be an explicit nil

### UnsetTags
`func (o *DeploymentCreateInputBody) UnsetTags()`

UnsetTags ensures that no value is present for Tags, not even an explicit nil
### GetUrl

`func (o *DeploymentCreateInputBody) GetUrl() Url`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *DeploymentCreateInputBody) GetUrlOk() (*Url, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *DeploymentCreateInputBody) SetUrl(v Url)`

SetUrl sets Url field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


