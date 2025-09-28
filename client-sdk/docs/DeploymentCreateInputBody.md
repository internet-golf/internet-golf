# DeploymentCreateInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**ExternalSource** | Pointer to **string** |  | [optional] 
**ExternalSourceType** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** | The primary identifier for the deployment. Defaults to the deployment&#39;s URL if it only has one URL; otherwise, the name must be specified when creating the deployment. | [optional] 
**PreserveExternalPath** | Pointer to **bool** |  | [optional] 
**Tags** | Pointer to **[]string** |  | [optional] 
**Urls** | [**[]Url**](Url.md) |  | 

## Methods

### NewDeploymentCreateInputBody

`func NewDeploymentCreateInputBody(urls []Url, ) *DeploymentCreateInputBody`

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

### GetName

`func (o *DeploymentCreateInputBody) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DeploymentCreateInputBody) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DeploymentCreateInputBody) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DeploymentCreateInputBody) HasName() bool`

HasName returns a boolean if a field has been set.

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
### GetUrls

`func (o *DeploymentCreateInputBody) GetUrls() []Url`

GetUrls returns the Urls field if non-nil, zero value otherwise.

### GetUrlsOk

`func (o *DeploymentCreateInputBody) GetUrlsOk() (*[]Url, bool)`

GetUrlsOk returns a tuple with the Urls field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrls

`func (o *DeploymentCreateInputBody) SetUrls(v []Url)`

SetUrls sets Urls field to given value.


### SetUrlsNil

`func (o *DeploymentCreateInputBody) SetUrlsNil(b bool)`

 SetUrlsNil sets the value for Urls to be an explicit nil

### UnsetUrls
`func (o *DeploymentCreateInputBody) UnsetUrls()`

UnsetUrls ensures that no value is present for Urls, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


