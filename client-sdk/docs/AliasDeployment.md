# AliasDeployment

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AliasedTo** | Pointer to **string** | The URL that this deployment is an alias for. | [optional] 
**CreatedAt** | **string** | When the deployment was created (string in ISO-8601 format.) | 
**ExternalSource** | Pointer to **string** | Original repository for this deployment&#39;s source. Can include a branch name. | [optional] 
**ExternalSourceType** | Pointer to **string** | Place where the original repository lives. | [optional] 
**Meta** | [**SiteMeta**](SiteMeta.md) |  | 
**Name** | **string** | Name for the deployment. This is just metadata; make it whatever you want. | 
**PreserveExternalPath** | Pointer to **bool** | If this is true and the deployment url has a path like \&quot;/thing\&quot;, then the \&quot;/thing\&quot; in the path will be transparently passed through to the underlying resource instead of being removed (which is the default) | [optional] 
**Redirect** | Pointer to **bool** | If this is true, visitors to this deployment&#39;s URL will be completely redirected to the URL that this alias is for. | [optional] 
**Tags** | Pointer to **[]string** | Tags used for metadata. | [optional] 
**Type** | **string** | Type of deployment contents. | 
**UpdatedAt** | **string** | When the deployment was last updated (string in ISO-8601 format.) | 
**Url** | **string** | URL that this deployment will appear at. The DNS for the domain has to be set up first. | 

## Methods

### NewAliasDeployment

`func NewAliasDeployment(createdAt string, meta SiteMeta, name string, type_ string, updatedAt string, url string, ) *AliasDeployment`

NewAliasDeployment instantiates a new AliasDeployment object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAliasDeploymentWithDefaults

`func NewAliasDeploymentWithDefaults() *AliasDeployment`

NewAliasDeploymentWithDefaults instantiates a new AliasDeployment object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAliasedTo

`func (o *AliasDeployment) GetAliasedTo() string`

GetAliasedTo returns the AliasedTo field if non-nil, zero value otherwise.

### GetAliasedToOk

`func (o *AliasDeployment) GetAliasedToOk() (*string, bool)`

GetAliasedToOk returns a tuple with the AliasedTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAliasedTo

`func (o *AliasDeployment) SetAliasedTo(v string)`

SetAliasedTo sets AliasedTo field to given value.

### HasAliasedTo

`func (o *AliasDeployment) HasAliasedTo() bool`

HasAliasedTo returns a boolean if a field has been set.

### GetCreatedAt

`func (o *AliasDeployment) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *AliasDeployment) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *AliasDeployment) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.


### GetExternalSource

`func (o *AliasDeployment) GetExternalSource() string`

GetExternalSource returns the ExternalSource field if non-nil, zero value otherwise.

### GetExternalSourceOk

`func (o *AliasDeployment) GetExternalSourceOk() (*string, bool)`

GetExternalSourceOk returns a tuple with the ExternalSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSource

`func (o *AliasDeployment) SetExternalSource(v string)`

SetExternalSource sets ExternalSource field to given value.

### HasExternalSource

`func (o *AliasDeployment) HasExternalSource() bool`

HasExternalSource returns a boolean if a field has been set.

### GetExternalSourceType

`func (o *AliasDeployment) GetExternalSourceType() string`

GetExternalSourceType returns the ExternalSourceType field if non-nil, zero value otherwise.

### GetExternalSourceTypeOk

`func (o *AliasDeployment) GetExternalSourceTypeOk() (*string, bool)`

GetExternalSourceTypeOk returns a tuple with the ExternalSourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSourceType

`func (o *AliasDeployment) SetExternalSourceType(v string)`

SetExternalSourceType sets ExternalSourceType field to given value.

### HasExternalSourceType

`func (o *AliasDeployment) HasExternalSourceType() bool`

HasExternalSourceType returns a boolean if a field has been set.

### GetMeta

`func (o *AliasDeployment) GetMeta() SiteMeta`

GetMeta returns the Meta field if non-nil, zero value otherwise.

### GetMetaOk

`func (o *AliasDeployment) GetMetaOk() (*SiteMeta, bool)`

GetMetaOk returns a tuple with the Meta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMeta

`func (o *AliasDeployment) SetMeta(v SiteMeta)`

SetMeta sets Meta field to given value.


### GetName

`func (o *AliasDeployment) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AliasDeployment) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AliasDeployment) SetName(v string)`

SetName sets Name field to given value.


### GetPreserveExternalPath

`func (o *AliasDeployment) GetPreserveExternalPath() bool`

GetPreserveExternalPath returns the PreserveExternalPath field if non-nil, zero value otherwise.

### GetPreserveExternalPathOk

`func (o *AliasDeployment) GetPreserveExternalPathOk() (*bool, bool)`

GetPreserveExternalPathOk returns a tuple with the PreserveExternalPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreserveExternalPath

`func (o *AliasDeployment) SetPreserveExternalPath(v bool)`

SetPreserveExternalPath sets PreserveExternalPath field to given value.

### HasPreserveExternalPath

`func (o *AliasDeployment) HasPreserveExternalPath() bool`

HasPreserveExternalPath returns a boolean if a field has been set.

### GetRedirect

`func (o *AliasDeployment) GetRedirect() bool`

GetRedirect returns the Redirect field if non-nil, zero value otherwise.

### GetRedirectOk

`func (o *AliasDeployment) GetRedirectOk() (*bool, bool)`

GetRedirectOk returns a tuple with the Redirect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirect

`func (o *AliasDeployment) SetRedirect(v bool)`

SetRedirect sets Redirect field to given value.

### HasRedirect

`func (o *AliasDeployment) HasRedirect() bool`

HasRedirect returns a boolean if a field has been set.

### GetTags

`func (o *AliasDeployment) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *AliasDeployment) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *AliasDeployment) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *AliasDeployment) HasTags() bool`

HasTags returns a boolean if a field has been set.

### SetTagsNil

`func (o *AliasDeployment) SetTagsNil(b bool)`

 SetTagsNil sets the value for Tags to be an explicit nil

### UnsetTags
`func (o *AliasDeployment) UnsetTags()`

UnsetTags ensures that no value is present for Tags, not even an explicit nil
### GetType

`func (o *AliasDeployment) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *AliasDeployment) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *AliasDeployment) SetType(v string)`

SetType sets Type field to given value.


### GetUpdatedAt

`func (o *AliasDeployment) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *AliasDeployment) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *AliasDeployment) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetUrl

`func (o *AliasDeployment) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *AliasDeployment) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *AliasDeployment) SetUrl(v string)`

SetUrl sets Url field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


