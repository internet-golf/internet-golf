# EmptyDeployment

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CreatedAt** | **string** | When the deployment was created (string in ISO-8601 format.) | 
**ExternalSource** | Pointer to **string** | Original repository for this deployment&#39;s source. Can include a branch name. | [optional] 
**ExternalSourceType** | Pointer to **string** | Place where the original repository lives. | [optional] 
**Meta** | [**SiteMeta**](SiteMeta.md) |  | 
**Name** | **string** | Name for the deployment. This is just metadata; make it whatever you want. | 
**NoContentYet** | Pointer to **bool** | Set to true to indicate that this deployment has not yet been set up. | [optional] 
**PreserveExternalPath** | Pointer to **bool** | If this is true and the deployment url has a path like \&quot;/thing\&quot;, then the \&quot;/thing\&quot; in the path will be transparently passed through to the underlying resource instead of being removed (which is the default) | [optional] 
**Tags** | Pointer to **[]string** | Tags used for metadata. | [optional] 
**Type** | **string** | Type of deployment contents. | 
**UpdatedAt** | **string** | When the deployment was last updated (string in ISO-8601 format.) | 
**Url** | **string** | URL that this deployment will appear at. The DNS for the domain has to be set up first. | 

## Methods

### NewEmptyDeployment

`func NewEmptyDeployment(createdAt string, meta SiteMeta, name string, type_ string, updatedAt string, url string, ) *EmptyDeployment`

NewEmptyDeployment instantiates a new EmptyDeployment object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEmptyDeploymentWithDefaults

`func NewEmptyDeploymentWithDefaults() *EmptyDeployment`

NewEmptyDeploymentWithDefaults instantiates a new EmptyDeployment object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCreatedAt

`func (o *EmptyDeployment) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *EmptyDeployment) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *EmptyDeployment) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.


### GetExternalSource

`func (o *EmptyDeployment) GetExternalSource() string`

GetExternalSource returns the ExternalSource field if non-nil, zero value otherwise.

### GetExternalSourceOk

`func (o *EmptyDeployment) GetExternalSourceOk() (*string, bool)`

GetExternalSourceOk returns a tuple with the ExternalSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSource

`func (o *EmptyDeployment) SetExternalSource(v string)`

SetExternalSource sets ExternalSource field to given value.

### HasExternalSource

`func (o *EmptyDeployment) HasExternalSource() bool`

HasExternalSource returns a boolean if a field has been set.

### GetExternalSourceType

`func (o *EmptyDeployment) GetExternalSourceType() string`

GetExternalSourceType returns the ExternalSourceType field if non-nil, zero value otherwise.

### GetExternalSourceTypeOk

`func (o *EmptyDeployment) GetExternalSourceTypeOk() (*string, bool)`

GetExternalSourceTypeOk returns a tuple with the ExternalSourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSourceType

`func (o *EmptyDeployment) SetExternalSourceType(v string)`

SetExternalSourceType sets ExternalSourceType field to given value.

### HasExternalSourceType

`func (o *EmptyDeployment) HasExternalSourceType() bool`

HasExternalSourceType returns a boolean if a field has been set.

### GetMeta

`func (o *EmptyDeployment) GetMeta() SiteMeta`

GetMeta returns the Meta field if non-nil, zero value otherwise.

### GetMetaOk

`func (o *EmptyDeployment) GetMetaOk() (*SiteMeta, bool)`

GetMetaOk returns a tuple with the Meta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMeta

`func (o *EmptyDeployment) SetMeta(v SiteMeta)`

SetMeta sets Meta field to given value.


### GetName

`func (o *EmptyDeployment) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *EmptyDeployment) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *EmptyDeployment) SetName(v string)`

SetName sets Name field to given value.


### GetNoContentYet

`func (o *EmptyDeployment) GetNoContentYet() bool`

GetNoContentYet returns the NoContentYet field if non-nil, zero value otherwise.

### GetNoContentYetOk

`func (o *EmptyDeployment) GetNoContentYetOk() (*bool, bool)`

GetNoContentYetOk returns a tuple with the NoContentYet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNoContentYet

`func (o *EmptyDeployment) SetNoContentYet(v bool)`

SetNoContentYet sets NoContentYet field to given value.

### HasNoContentYet

`func (o *EmptyDeployment) HasNoContentYet() bool`

HasNoContentYet returns a boolean if a field has been set.

### GetPreserveExternalPath

`func (o *EmptyDeployment) GetPreserveExternalPath() bool`

GetPreserveExternalPath returns the PreserveExternalPath field if non-nil, zero value otherwise.

### GetPreserveExternalPathOk

`func (o *EmptyDeployment) GetPreserveExternalPathOk() (*bool, bool)`

GetPreserveExternalPathOk returns a tuple with the PreserveExternalPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreserveExternalPath

`func (o *EmptyDeployment) SetPreserveExternalPath(v bool)`

SetPreserveExternalPath sets PreserveExternalPath field to given value.

### HasPreserveExternalPath

`func (o *EmptyDeployment) HasPreserveExternalPath() bool`

HasPreserveExternalPath returns a boolean if a field has been set.

### GetTags

`func (o *EmptyDeployment) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *EmptyDeployment) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *EmptyDeployment) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *EmptyDeployment) HasTags() bool`

HasTags returns a boolean if a field has been set.

### SetTagsNil

`func (o *EmptyDeployment) SetTagsNil(b bool)`

 SetTagsNil sets the value for Tags to be an explicit nil

### UnsetTags
`func (o *EmptyDeployment) UnsetTags()`

UnsetTags ensures that no value is present for Tags, not even an explicit nil
### GetType

`func (o *EmptyDeployment) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *EmptyDeployment) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *EmptyDeployment) SetType(v string)`

SetType sets Type field to given value.


### GetUpdatedAt

`func (o *EmptyDeployment) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *EmptyDeployment) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *EmptyDeployment) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetUrl

`func (o *EmptyDeployment) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *EmptyDeployment) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *EmptyDeployment) SetUrl(v string)`

SetUrl sets Url field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


