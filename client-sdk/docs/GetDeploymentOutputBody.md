# GetDeploymentOutputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**ExternalSource** | Pointer to **string** |  | [optional] 
**ExternalSourceType** | Pointer to **string** |  | [optional] 
**HasContent** | **bool** |  | 
**Name** | Pointer to **string** | The primary identifier for the deployment. Defaults to the deployment&#39;s URL if it only has one URL; otherwise, the name must be specified when creating the deployment. | [optional] 
**PreserveExternalPath** | Pointer to **bool** |  | [optional] 
**ServedThing** | **string** |  | 
**ServedThingType** | **string** |  | 
**Tags** | Pointer to **[]string** |  | [optional] 
**Urls** | [**[]Url**](Url.md) |  | 

## Methods

### NewGetDeploymentOutputBody

`func NewGetDeploymentOutputBody(hasContent bool, servedThing string, servedThingType string, urls []Url, ) *GetDeploymentOutputBody`

NewGetDeploymentOutputBody instantiates a new GetDeploymentOutputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetDeploymentOutputBodyWithDefaults

`func NewGetDeploymentOutputBodyWithDefaults() *GetDeploymentOutputBody`

NewGetDeploymentOutputBodyWithDefaults instantiates a new GetDeploymentOutputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *GetDeploymentOutputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *GetDeploymentOutputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *GetDeploymentOutputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *GetDeploymentOutputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetExternalSource

`func (o *GetDeploymentOutputBody) GetExternalSource() string`

GetExternalSource returns the ExternalSource field if non-nil, zero value otherwise.

### GetExternalSourceOk

`func (o *GetDeploymentOutputBody) GetExternalSourceOk() (*string, bool)`

GetExternalSourceOk returns a tuple with the ExternalSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSource

`func (o *GetDeploymentOutputBody) SetExternalSource(v string)`

SetExternalSource sets ExternalSource field to given value.

### HasExternalSource

`func (o *GetDeploymentOutputBody) HasExternalSource() bool`

HasExternalSource returns a boolean if a field has been set.

### GetExternalSourceType

`func (o *GetDeploymentOutputBody) GetExternalSourceType() string`

GetExternalSourceType returns the ExternalSourceType field if non-nil, zero value otherwise.

### GetExternalSourceTypeOk

`func (o *GetDeploymentOutputBody) GetExternalSourceTypeOk() (*string, bool)`

GetExternalSourceTypeOk returns a tuple with the ExternalSourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSourceType

`func (o *GetDeploymentOutputBody) SetExternalSourceType(v string)`

SetExternalSourceType sets ExternalSourceType field to given value.

### HasExternalSourceType

`func (o *GetDeploymentOutputBody) HasExternalSourceType() bool`

HasExternalSourceType returns a boolean if a field has been set.

### GetHasContent

`func (o *GetDeploymentOutputBody) GetHasContent() bool`

GetHasContent returns the HasContent field if non-nil, zero value otherwise.

### GetHasContentOk

`func (o *GetDeploymentOutputBody) GetHasContentOk() (*bool, bool)`

GetHasContentOk returns a tuple with the HasContent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHasContent

`func (o *GetDeploymentOutputBody) SetHasContent(v bool)`

SetHasContent sets HasContent field to given value.


### GetName

`func (o *GetDeploymentOutputBody) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *GetDeploymentOutputBody) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *GetDeploymentOutputBody) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *GetDeploymentOutputBody) HasName() bool`

HasName returns a boolean if a field has been set.

### GetPreserveExternalPath

`func (o *GetDeploymentOutputBody) GetPreserveExternalPath() bool`

GetPreserveExternalPath returns the PreserveExternalPath field if non-nil, zero value otherwise.

### GetPreserveExternalPathOk

`func (o *GetDeploymentOutputBody) GetPreserveExternalPathOk() (*bool, bool)`

GetPreserveExternalPathOk returns a tuple with the PreserveExternalPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreserveExternalPath

`func (o *GetDeploymentOutputBody) SetPreserveExternalPath(v bool)`

SetPreserveExternalPath sets PreserveExternalPath field to given value.

### HasPreserveExternalPath

`func (o *GetDeploymentOutputBody) HasPreserveExternalPath() bool`

HasPreserveExternalPath returns a boolean if a field has been set.

### GetServedThing

`func (o *GetDeploymentOutputBody) GetServedThing() string`

GetServedThing returns the ServedThing field if non-nil, zero value otherwise.

### GetServedThingOk

`func (o *GetDeploymentOutputBody) GetServedThingOk() (*string, bool)`

GetServedThingOk returns a tuple with the ServedThing field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServedThing

`func (o *GetDeploymentOutputBody) SetServedThing(v string)`

SetServedThing sets ServedThing field to given value.


### GetServedThingType

`func (o *GetDeploymentOutputBody) GetServedThingType() string`

GetServedThingType returns the ServedThingType field if non-nil, zero value otherwise.

### GetServedThingTypeOk

`func (o *GetDeploymentOutputBody) GetServedThingTypeOk() (*string, bool)`

GetServedThingTypeOk returns a tuple with the ServedThingType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServedThingType

`func (o *GetDeploymentOutputBody) SetServedThingType(v string)`

SetServedThingType sets ServedThingType field to given value.


### GetTags

`func (o *GetDeploymentOutputBody) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *GetDeploymentOutputBody) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *GetDeploymentOutputBody) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *GetDeploymentOutputBody) HasTags() bool`

HasTags returns a boolean if a field has been set.

### SetTagsNil

`func (o *GetDeploymentOutputBody) SetTagsNil(b bool)`

 SetTagsNil sets the value for Tags to be an explicit nil

### UnsetTags
`func (o *GetDeploymentOutputBody) UnsetTags()`

UnsetTags ensures that no value is present for Tags, not even an explicit nil
### GetUrls

`func (o *GetDeploymentOutputBody) GetUrls() []Url`

GetUrls returns the Urls field if non-nil, zero value otherwise.

### GetUrlsOk

`func (o *GetDeploymentOutputBody) GetUrlsOk() (*[]Url, bool)`

GetUrlsOk returns a tuple with the Urls field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrls

`func (o *GetDeploymentOutputBody) SetUrls(v []Url)`

SetUrls sets Urls field to given value.


### SetUrlsNil

`func (o *GetDeploymentOutputBody) SetUrlsNil(b bool)`

 SetUrlsNil sets the value for Urls to be an explicit nil

### UnsetUrls
`func (o *GetDeploymentOutputBody) UnsetUrls()`

UnsetUrls ensures that no value is present for Urls, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


