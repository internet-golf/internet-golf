# DeploymentBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AliasedTo** | Pointer to **string** | The URL that this deployment is an alias for. | [optional] 
**ExternalSource** | Pointer to **string** | Original repository for this deployment&#39;s source. Can include a branch name. | [optional] 
**ExternalSourceType** | Pointer to **string** | Place where the original repository lives. | [optional] 
**PreserveExternalPath** | Pointer to **bool** | if this is true and the deployment url has a path like \&quot;/thing\&quot;, then the \&quot;/thing\&quot; in the path will be transparently passed through to the underlying resource instead of being removed (which is the default) | [optional] 
**Redirect** | Pointer to **bool** | If this is true, visitors to this deployment&#39;s URL will be completely redirected to the URL that this alias is for. | [optional] 
**ServerContentLocation** | Pointer to **string** |  | [optional] 
**SpaMode** | Pointer to **bool** |  | [optional] 
**Tags** | Pointer to **[]string** | Tags used for metadata. | [optional] 
**Type** | Pointer to **string** | Type of deployment contents; can be StaticSite, Alias, or Empty. | [optional] 
**Url** | **string** | URL that this deployment will appear at. The DNS for the domain has to be set up first. | 

## Methods

### NewDeploymentBody

`func NewDeploymentBody(url string, ) *DeploymentBody`

NewDeploymentBody instantiates a new DeploymentBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeploymentBodyWithDefaults

`func NewDeploymentBodyWithDefaults() *DeploymentBody`

NewDeploymentBodyWithDefaults instantiates a new DeploymentBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAliasedTo

`func (o *DeploymentBody) GetAliasedTo() string`

GetAliasedTo returns the AliasedTo field if non-nil, zero value otherwise.

### GetAliasedToOk

`func (o *DeploymentBody) GetAliasedToOk() (*string, bool)`

GetAliasedToOk returns a tuple with the AliasedTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAliasedTo

`func (o *DeploymentBody) SetAliasedTo(v string)`

SetAliasedTo sets AliasedTo field to given value.

### HasAliasedTo

`func (o *DeploymentBody) HasAliasedTo() bool`

HasAliasedTo returns a boolean if a field has been set.

### GetExternalSource

`func (o *DeploymentBody) GetExternalSource() string`

GetExternalSource returns the ExternalSource field if non-nil, zero value otherwise.

### GetExternalSourceOk

`func (o *DeploymentBody) GetExternalSourceOk() (*string, bool)`

GetExternalSourceOk returns a tuple with the ExternalSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSource

`func (o *DeploymentBody) SetExternalSource(v string)`

SetExternalSource sets ExternalSource field to given value.

### HasExternalSource

`func (o *DeploymentBody) HasExternalSource() bool`

HasExternalSource returns a boolean if a field has been set.

### GetExternalSourceType

`func (o *DeploymentBody) GetExternalSourceType() string`

GetExternalSourceType returns the ExternalSourceType field if non-nil, zero value otherwise.

### GetExternalSourceTypeOk

`func (o *DeploymentBody) GetExternalSourceTypeOk() (*string, bool)`

GetExternalSourceTypeOk returns a tuple with the ExternalSourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSourceType

`func (o *DeploymentBody) SetExternalSourceType(v string)`

SetExternalSourceType sets ExternalSourceType field to given value.

### HasExternalSourceType

`func (o *DeploymentBody) HasExternalSourceType() bool`

HasExternalSourceType returns a boolean if a field has been set.

### GetPreserveExternalPath

`func (o *DeploymentBody) GetPreserveExternalPath() bool`

GetPreserveExternalPath returns the PreserveExternalPath field if non-nil, zero value otherwise.

### GetPreserveExternalPathOk

`func (o *DeploymentBody) GetPreserveExternalPathOk() (*bool, bool)`

GetPreserveExternalPathOk returns a tuple with the PreserveExternalPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreserveExternalPath

`func (o *DeploymentBody) SetPreserveExternalPath(v bool)`

SetPreserveExternalPath sets PreserveExternalPath field to given value.

### HasPreserveExternalPath

`func (o *DeploymentBody) HasPreserveExternalPath() bool`

HasPreserveExternalPath returns a boolean if a field has been set.

### GetRedirect

`func (o *DeploymentBody) GetRedirect() bool`

GetRedirect returns the Redirect field if non-nil, zero value otherwise.

### GetRedirectOk

`func (o *DeploymentBody) GetRedirectOk() (*bool, bool)`

GetRedirectOk returns a tuple with the Redirect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirect

`func (o *DeploymentBody) SetRedirect(v bool)`

SetRedirect sets Redirect field to given value.

### HasRedirect

`func (o *DeploymentBody) HasRedirect() bool`

HasRedirect returns a boolean if a field has been set.

### GetServerContentLocation

`func (o *DeploymentBody) GetServerContentLocation() string`

GetServerContentLocation returns the ServerContentLocation field if non-nil, zero value otherwise.

### GetServerContentLocationOk

`func (o *DeploymentBody) GetServerContentLocationOk() (*string, bool)`

GetServerContentLocationOk returns a tuple with the ServerContentLocation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServerContentLocation

`func (o *DeploymentBody) SetServerContentLocation(v string)`

SetServerContentLocation sets ServerContentLocation field to given value.

### HasServerContentLocation

`func (o *DeploymentBody) HasServerContentLocation() bool`

HasServerContentLocation returns a boolean if a field has been set.

### GetSpaMode

`func (o *DeploymentBody) GetSpaMode() bool`

GetSpaMode returns the SpaMode field if non-nil, zero value otherwise.

### GetSpaModeOk

`func (o *DeploymentBody) GetSpaModeOk() (*bool, bool)`

GetSpaModeOk returns a tuple with the SpaMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpaMode

`func (o *DeploymentBody) SetSpaMode(v bool)`

SetSpaMode sets SpaMode field to given value.

### HasSpaMode

`func (o *DeploymentBody) HasSpaMode() bool`

HasSpaMode returns a boolean if a field has been set.

### GetTags

`func (o *DeploymentBody) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *DeploymentBody) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *DeploymentBody) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *DeploymentBody) HasTags() bool`

HasTags returns a boolean if a field has been set.

### SetTagsNil

`func (o *DeploymentBody) SetTagsNil(b bool)`

 SetTagsNil sets the value for Tags to be an explicit nil

### UnsetTags
`func (o *DeploymentBody) UnsetTags()`

UnsetTags ensures that no value is present for Tags, not even an explicit nil
### GetType

`func (o *DeploymentBody) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *DeploymentBody) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *DeploymentBody) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *DeploymentBody) HasType() bool`

HasType returns a boolean if a field has been set.

### GetUrl

`func (o *DeploymentBody) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *DeploymentBody) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *DeploymentBody) SetUrl(v string)`

SetUrl sets Url field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


