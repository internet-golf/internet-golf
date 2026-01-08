# DeploymentModel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AliasedTo** | Pointer to **string** | The URL that this deployment is an alias for. | [optional] 
**CreatedAt** | **string** | When the deployment was created (string in ISO-8601 format.) | 
**ExternalSource** | Pointer to **string** | Original repository for this deployment&#39;s source. Can include a branch name. | [optional] 
**ExternalSourceType** | Pointer to **string** | Place where the original repository lives. | [optional] 
**Meta** | [**SiteMeta**](SiteMeta.md) |  | 
**Name** | **string** | Name for the deployment. This is just metadata; make it whatever you want. | 
**NoContentYet** | Pointer to **bool** | Set to true to indicate that this deployment has not yet been set up. | [optional] 
**PreserveExternalPath** | Pointer to **bool** | If this is true and the deployment url has a path like \&quot;/thing\&quot;, then the \&quot;/thing\&quot; in the path will be transparently passed through to the underlying resource instead of being removed (which is the default) | [optional] 
**Redirect** | Pointer to **bool** | If this is true, visitors to this deployment&#39;s URL will be completely redirected to the URL that this alias is for. | [optional] 
**ServerContentLocation** | Pointer to **string** | The path to this deployment&#39;s files on the server. | [optional] 
**SpaMode** | Pointer to **bool** | Whether this deployment is set up to support a Single Page App by using /index.html as a fallback for all requests. | [optional] 
**Tags** | Pointer to **[]string** | Tags used for metadata. | [optional] 
**Type** | **string** | Type of deployment contents. | 
**UpdatedAt** | **string** | When the deployment was last updated (string in ISO-8601 format.) | 
**Url** | **string** | URL that this deployment will appear at. The DNS for the domain has to be set up first. | 

## Methods

### NewDeploymentModel

`func NewDeploymentModel(createdAt string, meta SiteMeta, name string, type_ string, updatedAt string, url string, ) *DeploymentModel`

NewDeploymentModel instantiates a new DeploymentModel object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeploymentModelWithDefaults

`func NewDeploymentModelWithDefaults() *DeploymentModel`

NewDeploymentModelWithDefaults instantiates a new DeploymentModel object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAliasedTo

`func (o *DeploymentModel) GetAliasedTo() string`

GetAliasedTo returns the AliasedTo field if non-nil, zero value otherwise.

### GetAliasedToOk

`func (o *DeploymentModel) GetAliasedToOk() (*string, bool)`

GetAliasedToOk returns a tuple with the AliasedTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAliasedTo

`func (o *DeploymentModel) SetAliasedTo(v string)`

SetAliasedTo sets AliasedTo field to given value.

### HasAliasedTo

`func (o *DeploymentModel) HasAliasedTo() bool`

HasAliasedTo returns a boolean if a field has been set.

### GetCreatedAt

`func (o *DeploymentModel) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *DeploymentModel) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *DeploymentModel) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.


### GetExternalSource

`func (o *DeploymentModel) GetExternalSource() string`

GetExternalSource returns the ExternalSource field if non-nil, zero value otherwise.

### GetExternalSourceOk

`func (o *DeploymentModel) GetExternalSourceOk() (*string, bool)`

GetExternalSourceOk returns a tuple with the ExternalSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSource

`func (o *DeploymentModel) SetExternalSource(v string)`

SetExternalSource sets ExternalSource field to given value.

### HasExternalSource

`func (o *DeploymentModel) HasExternalSource() bool`

HasExternalSource returns a boolean if a field has been set.

### GetExternalSourceType

`func (o *DeploymentModel) GetExternalSourceType() string`

GetExternalSourceType returns the ExternalSourceType field if non-nil, zero value otherwise.

### GetExternalSourceTypeOk

`func (o *DeploymentModel) GetExternalSourceTypeOk() (*string, bool)`

GetExternalSourceTypeOk returns a tuple with the ExternalSourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSourceType

`func (o *DeploymentModel) SetExternalSourceType(v string)`

SetExternalSourceType sets ExternalSourceType field to given value.

### HasExternalSourceType

`func (o *DeploymentModel) HasExternalSourceType() bool`

HasExternalSourceType returns a boolean if a field has been set.

### GetMeta

`func (o *DeploymentModel) GetMeta() SiteMeta`

GetMeta returns the Meta field if non-nil, zero value otherwise.

### GetMetaOk

`func (o *DeploymentModel) GetMetaOk() (*SiteMeta, bool)`

GetMetaOk returns a tuple with the Meta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMeta

`func (o *DeploymentModel) SetMeta(v SiteMeta)`

SetMeta sets Meta field to given value.


### GetName

`func (o *DeploymentModel) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DeploymentModel) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DeploymentModel) SetName(v string)`

SetName sets Name field to given value.


### GetNoContentYet

`func (o *DeploymentModel) GetNoContentYet() bool`

GetNoContentYet returns the NoContentYet field if non-nil, zero value otherwise.

### GetNoContentYetOk

`func (o *DeploymentModel) GetNoContentYetOk() (*bool, bool)`

GetNoContentYetOk returns a tuple with the NoContentYet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNoContentYet

`func (o *DeploymentModel) SetNoContentYet(v bool)`

SetNoContentYet sets NoContentYet field to given value.

### HasNoContentYet

`func (o *DeploymentModel) HasNoContentYet() bool`

HasNoContentYet returns a boolean if a field has been set.

### GetPreserveExternalPath

`func (o *DeploymentModel) GetPreserveExternalPath() bool`

GetPreserveExternalPath returns the PreserveExternalPath field if non-nil, zero value otherwise.

### GetPreserveExternalPathOk

`func (o *DeploymentModel) GetPreserveExternalPathOk() (*bool, bool)`

GetPreserveExternalPathOk returns a tuple with the PreserveExternalPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreserveExternalPath

`func (o *DeploymentModel) SetPreserveExternalPath(v bool)`

SetPreserveExternalPath sets PreserveExternalPath field to given value.

### HasPreserveExternalPath

`func (o *DeploymentModel) HasPreserveExternalPath() bool`

HasPreserveExternalPath returns a boolean if a field has been set.

### GetRedirect

`func (o *DeploymentModel) GetRedirect() bool`

GetRedirect returns the Redirect field if non-nil, zero value otherwise.

### GetRedirectOk

`func (o *DeploymentModel) GetRedirectOk() (*bool, bool)`

GetRedirectOk returns a tuple with the Redirect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirect

`func (o *DeploymentModel) SetRedirect(v bool)`

SetRedirect sets Redirect field to given value.

### HasRedirect

`func (o *DeploymentModel) HasRedirect() bool`

HasRedirect returns a boolean if a field has been set.

### GetServerContentLocation

`func (o *DeploymentModel) GetServerContentLocation() string`

GetServerContentLocation returns the ServerContentLocation field if non-nil, zero value otherwise.

### GetServerContentLocationOk

`func (o *DeploymentModel) GetServerContentLocationOk() (*string, bool)`

GetServerContentLocationOk returns a tuple with the ServerContentLocation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServerContentLocation

`func (o *DeploymentModel) SetServerContentLocation(v string)`

SetServerContentLocation sets ServerContentLocation field to given value.

### HasServerContentLocation

`func (o *DeploymentModel) HasServerContentLocation() bool`

HasServerContentLocation returns a boolean if a field has been set.

### GetSpaMode

`func (o *DeploymentModel) GetSpaMode() bool`

GetSpaMode returns the SpaMode field if non-nil, zero value otherwise.

### GetSpaModeOk

`func (o *DeploymentModel) GetSpaModeOk() (*bool, bool)`

GetSpaModeOk returns a tuple with the SpaMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpaMode

`func (o *DeploymentModel) SetSpaMode(v bool)`

SetSpaMode sets SpaMode field to given value.

### HasSpaMode

`func (o *DeploymentModel) HasSpaMode() bool`

HasSpaMode returns a boolean if a field has been set.

### GetTags

`func (o *DeploymentModel) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *DeploymentModel) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *DeploymentModel) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *DeploymentModel) HasTags() bool`

HasTags returns a boolean if a field has been set.

### SetTagsNil

`func (o *DeploymentModel) SetTagsNil(b bool)`

 SetTagsNil sets the value for Tags to be an explicit nil

### UnsetTags
`func (o *DeploymentModel) UnsetTags()`

UnsetTags ensures that no value is present for Tags, not even an explicit nil
### GetType

`func (o *DeploymentModel) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *DeploymentModel) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *DeploymentModel) SetType(v string)`

SetType sets Type field to given value.


### GetUpdatedAt

`func (o *DeploymentModel) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *DeploymentModel) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *DeploymentModel) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetUrl

`func (o *DeploymentModel) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *DeploymentModel) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *DeploymentModel) SetUrl(v string)`

SetUrl sets Url field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


