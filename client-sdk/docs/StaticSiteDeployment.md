# StaticSiteDeployment

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CreatedAt** | **string** | When the deployment was created (string in ISO-8601 format.) | 
**ExternalSource** | Pointer to **string** | Original repository for this deployment&#39;s source. Can include a branch name. | [optional] 
**ExternalSourceType** | Pointer to **string** | Place where the original repository lives. | [optional] 
**Meta** | [**SiteMeta**](SiteMeta.md) |  | 
**Name** | **string** | Name for the deployment. This is just metadata; make it whatever you want. | 
**PreserveExternalPath** | Pointer to **bool** | If this is true and the deployment url has a path like \&quot;/thing\&quot;, then the \&quot;/thing\&quot; in the path will be transparently passed through to the underlying resource instead of being removed (which is the default) | [optional] 
**ServerContentLocation** | Pointer to **string** | The path to this deployment&#39;s files on the server. | [optional] 
**SpaMode** | Pointer to **bool** | Whether this deployment is set up to support a Single Page App by using /index.html as a fallback for all requests. | [optional] 
**Tags** | Pointer to **[]string** | Tags used for metadata. | [optional] 
**Type** | **string** | Type of deployment contents. | 
**UpdatedAt** | **string** | When the deployment was last updated (string in ISO-8601 format.) | 
**Url** | **string** | URL that this deployment will appear at. The DNS for the domain has to be set up first. | 

## Methods

### NewStaticSiteDeployment

`func NewStaticSiteDeployment(createdAt string, meta SiteMeta, name string, type_ string, updatedAt string, url string, ) *StaticSiteDeployment`

NewStaticSiteDeployment instantiates a new StaticSiteDeployment object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStaticSiteDeploymentWithDefaults

`func NewStaticSiteDeploymentWithDefaults() *StaticSiteDeployment`

NewStaticSiteDeploymentWithDefaults instantiates a new StaticSiteDeployment object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCreatedAt

`func (o *StaticSiteDeployment) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *StaticSiteDeployment) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *StaticSiteDeployment) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.


### GetExternalSource

`func (o *StaticSiteDeployment) GetExternalSource() string`

GetExternalSource returns the ExternalSource field if non-nil, zero value otherwise.

### GetExternalSourceOk

`func (o *StaticSiteDeployment) GetExternalSourceOk() (*string, bool)`

GetExternalSourceOk returns a tuple with the ExternalSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSource

`func (o *StaticSiteDeployment) SetExternalSource(v string)`

SetExternalSource sets ExternalSource field to given value.

### HasExternalSource

`func (o *StaticSiteDeployment) HasExternalSource() bool`

HasExternalSource returns a boolean if a field has been set.

### GetExternalSourceType

`func (o *StaticSiteDeployment) GetExternalSourceType() string`

GetExternalSourceType returns the ExternalSourceType field if non-nil, zero value otherwise.

### GetExternalSourceTypeOk

`func (o *StaticSiteDeployment) GetExternalSourceTypeOk() (*string, bool)`

GetExternalSourceTypeOk returns a tuple with the ExternalSourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSourceType

`func (o *StaticSiteDeployment) SetExternalSourceType(v string)`

SetExternalSourceType sets ExternalSourceType field to given value.

### HasExternalSourceType

`func (o *StaticSiteDeployment) HasExternalSourceType() bool`

HasExternalSourceType returns a boolean if a field has been set.

### GetMeta

`func (o *StaticSiteDeployment) GetMeta() SiteMeta`

GetMeta returns the Meta field if non-nil, zero value otherwise.

### GetMetaOk

`func (o *StaticSiteDeployment) GetMetaOk() (*SiteMeta, bool)`

GetMetaOk returns a tuple with the Meta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMeta

`func (o *StaticSiteDeployment) SetMeta(v SiteMeta)`

SetMeta sets Meta field to given value.


### GetName

`func (o *StaticSiteDeployment) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *StaticSiteDeployment) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *StaticSiteDeployment) SetName(v string)`

SetName sets Name field to given value.


### GetPreserveExternalPath

`func (o *StaticSiteDeployment) GetPreserveExternalPath() bool`

GetPreserveExternalPath returns the PreserveExternalPath field if non-nil, zero value otherwise.

### GetPreserveExternalPathOk

`func (o *StaticSiteDeployment) GetPreserveExternalPathOk() (*bool, bool)`

GetPreserveExternalPathOk returns a tuple with the PreserveExternalPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreserveExternalPath

`func (o *StaticSiteDeployment) SetPreserveExternalPath(v bool)`

SetPreserveExternalPath sets PreserveExternalPath field to given value.

### HasPreserveExternalPath

`func (o *StaticSiteDeployment) HasPreserveExternalPath() bool`

HasPreserveExternalPath returns a boolean if a field has been set.

### GetServerContentLocation

`func (o *StaticSiteDeployment) GetServerContentLocation() string`

GetServerContentLocation returns the ServerContentLocation field if non-nil, zero value otherwise.

### GetServerContentLocationOk

`func (o *StaticSiteDeployment) GetServerContentLocationOk() (*string, bool)`

GetServerContentLocationOk returns a tuple with the ServerContentLocation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServerContentLocation

`func (o *StaticSiteDeployment) SetServerContentLocation(v string)`

SetServerContentLocation sets ServerContentLocation field to given value.

### HasServerContentLocation

`func (o *StaticSiteDeployment) HasServerContentLocation() bool`

HasServerContentLocation returns a boolean if a field has been set.

### GetSpaMode

`func (o *StaticSiteDeployment) GetSpaMode() bool`

GetSpaMode returns the SpaMode field if non-nil, zero value otherwise.

### GetSpaModeOk

`func (o *StaticSiteDeployment) GetSpaModeOk() (*bool, bool)`

GetSpaModeOk returns a tuple with the SpaMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpaMode

`func (o *StaticSiteDeployment) SetSpaMode(v bool)`

SetSpaMode sets SpaMode field to given value.

### HasSpaMode

`func (o *StaticSiteDeployment) HasSpaMode() bool`

HasSpaMode returns a boolean if a field has been set.

### GetTags

`func (o *StaticSiteDeployment) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *StaticSiteDeployment) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *StaticSiteDeployment) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *StaticSiteDeployment) HasTags() bool`

HasTags returns a boolean if a field has been set.

### SetTagsNil

`func (o *StaticSiteDeployment) SetTagsNil(b bool)`

 SetTagsNil sets the value for Tags to be an explicit nil

### UnsetTags
`func (o *StaticSiteDeployment) UnsetTags()`

UnsetTags ensures that no value is present for Tags, not even an explicit nil
### GetType

`func (o *StaticSiteDeployment) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *StaticSiteDeployment) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *StaticSiteDeployment) SetType(v string)`

SetType sets Type field to given value.


### GetUpdatedAt

`func (o *StaticSiteDeployment) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *StaticSiteDeployment) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *StaticSiteDeployment) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetUrl

`func (o *StaticSiteDeployment) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *StaticSiteDeployment) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *StaticSiteDeployment) SetUrl(v string)`

SetUrl sets Url field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


