# DeploymentResponseBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AliasedTo** | Pointer to **string** | The URL that this deployment is an alias for. | [optional] 
**ExternalSource** | Pointer to **string** | Original repository for this deployment&#39;s source. Can include a branch name. | [optional] 
**ExternalSourceType** | Pointer to **string** | Place where the original repository lives. | [optional] 
**NoContentYet** | **bool** | Set to true if this deployment has not yet been set up. | 
**PreserveExternalPath** | Pointer to **bool** | if this is true and the deployment url has a path like \&quot;/thing\&quot;, then the \&quot;/thing\&quot; in the path will be transparently passed through to the underlying resource instead of being removed (which is the default) | [optional] 
**Redirect** | Pointer to **bool** | If this is true, visitors to this deployment&#39;s URL will be completely redirected to the URL that this alias is for. | [optional] 
**ServerContentLocation** | Pointer to **string** | The path to this deployment&#39;s files on the server. | [optional] 
**SpaMode** | Pointer to **bool** | Whether this deployment is set up to support a Single Page App by using /index.html as a fallback for all requests. | [optional] 
**Tags** | Pointer to **[]string** | Tags used for metadata. | [optional] 
**Type** | **string** | Type of deployment contents. | 
**Url** | **string** | URL that this deployment will appear at. The DNS for the domain has to be set up first. | 

## Methods

### NewDeploymentResponseBody

`func NewDeploymentResponseBody(noContentYet bool, type_ string, url string, ) *DeploymentResponseBody`

NewDeploymentResponseBody instantiates a new DeploymentResponseBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeploymentResponseBodyWithDefaults

`func NewDeploymentResponseBodyWithDefaults() *DeploymentResponseBody`

NewDeploymentResponseBodyWithDefaults instantiates a new DeploymentResponseBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAliasedTo

`func (o *DeploymentResponseBody) GetAliasedTo() string`

GetAliasedTo returns the AliasedTo field if non-nil, zero value otherwise.

### GetAliasedToOk

`func (o *DeploymentResponseBody) GetAliasedToOk() (*string, bool)`

GetAliasedToOk returns a tuple with the AliasedTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAliasedTo

`func (o *DeploymentResponseBody) SetAliasedTo(v string)`

SetAliasedTo sets AliasedTo field to given value.

### HasAliasedTo

`func (o *DeploymentResponseBody) HasAliasedTo() bool`

HasAliasedTo returns a boolean if a field has been set.

### GetExternalSource

`func (o *DeploymentResponseBody) GetExternalSource() string`

GetExternalSource returns the ExternalSource field if non-nil, zero value otherwise.

### GetExternalSourceOk

`func (o *DeploymentResponseBody) GetExternalSourceOk() (*string, bool)`

GetExternalSourceOk returns a tuple with the ExternalSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSource

`func (o *DeploymentResponseBody) SetExternalSource(v string)`

SetExternalSource sets ExternalSource field to given value.

### HasExternalSource

`func (o *DeploymentResponseBody) HasExternalSource() bool`

HasExternalSource returns a boolean if a field has been set.

### GetExternalSourceType

`func (o *DeploymentResponseBody) GetExternalSourceType() string`

GetExternalSourceType returns the ExternalSourceType field if non-nil, zero value otherwise.

### GetExternalSourceTypeOk

`func (o *DeploymentResponseBody) GetExternalSourceTypeOk() (*string, bool)`

GetExternalSourceTypeOk returns a tuple with the ExternalSourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSourceType

`func (o *DeploymentResponseBody) SetExternalSourceType(v string)`

SetExternalSourceType sets ExternalSourceType field to given value.

### HasExternalSourceType

`func (o *DeploymentResponseBody) HasExternalSourceType() bool`

HasExternalSourceType returns a boolean if a field has been set.

### GetNoContentYet

`func (o *DeploymentResponseBody) GetNoContentYet() bool`

GetNoContentYet returns the NoContentYet field if non-nil, zero value otherwise.

### GetNoContentYetOk

`func (o *DeploymentResponseBody) GetNoContentYetOk() (*bool, bool)`

GetNoContentYetOk returns a tuple with the NoContentYet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNoContentYet

`func (o *DeploymentResponseBody) SetNoContentYet(v bool)`

SetNoContentYet sets NoContentYet field to given value.


### GetPreserveExternalPath

`func (o *DeploymentResponseBody) GetPreserveExternalPath() bool`

GetPreserveExternalPath returns the PreserveExternalPath field if non-nil, zero value otherwise.

### GetPreserveExternalPathOk

`func (o *DeploymentResponseBody) GetPreserveExternalPathOk() (*bool, bool)`

GetPreserveExternalPathOk returns a tuple with the PreserveExternalPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreserveExternalPath

`func (o *DeploymentResponseBody) SetPreserveExternalPath(v bool)`

SetPreserveExternalPath sets PreserveExternalPath field to given value.

### HasPreserveExternalPath

`func (o *DeploymentResponseBody) HasPreserveExternalPath() bool`

HasPreserveExternalPath returns a boolean if a field has been set.

### GetRedirect

`func (o *DeploymentResponseBody) GetRedirect() bool`

GetRedirect returns the Redirect field if non-nil, zero value otherwise.

### GetRedirectOk

`func (o *DeploymentResponseBody) GetRedirectOk() (*bool, bool)`

GetRedirectOk returns a tuple with the Redirect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirect

`func (o *DeploymentResponseBody) SetRedirect(v bool)`

SetRedirect sets Redirect field to given value.

### HasRedirect

`func (o *DeploymentResponseBody) HasRedirect() bool`

HasRedirect returns a boolean if a field has been set.

### GetServerContentLocation

`func (o *DeploymentResponseBody) GetServerContentLocation() string`

GetServerContentLocation returns the ServerContentLocation field if non-nil, zero value otherwise.

### GetServerContentLocationOk

`func (o *DeploymentResponseBody) GetServerContentLocationOk() (*string, bool)`

GetServerContentLocationOk returns a tuple with the ServerContentLocation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServerContentLocation

`func (o *DeploymentResponseBody) SetServerContentLocation(v string)`

SetServerContentLocation sets ServerContentLocation field to given value.

### HasServerContentLocation

`func (o *DeploymentResponseBody) HasServerContentLocation() bool`

HasServerContentLocation returns a boolean if a field has been set.

### GetSpaMode

`func (o *DeploymentResponseBody) GetSpaMode() bool`

GetSpaMode returns the SpaMode field if non-nil, zero value otherwise.

### GetSpaModeOk

`func (o *DeploymentResponseBody) GetSpaModeOk() (*bool, bool)`

GetSpaModeOk returns a tuple with the SpaMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpaMode

`func (o *DeploymentResponseBody) SetSpaMode(v bool)`

SetSpaMode sets SpaMode field to given value.

### HasSpaMode

`func (o *DeploymentResponseBody) HasSpaMode() bool`

HasSpaMode returns a boolean if a field has been set.

### GetTags

`func (o *DeploymentResponseBody) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *DeploymentResponseBody) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *DeploymentResponseBody) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *DeploymentResponseBody) HasTags() bool`

HasTags returns a boolean if a field has been set.

### SetTagsNil

`func (o *DeploymentResponseBody) SetTagsNil(b bool)`

 SetTagsNil sets the value for Tags to be an explicit nil

### UnsetTags
`func (o *DeploymentResponseBody) UnsetTags()`

UnsetTags ensures that no value is present for Tags, not even an explicit nil
### GetType

`func (o *DeploymentResponseBody) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *DeploymentResponseBody) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *DeploymentResponseBody) SetType(v string)`

SetType sets Type field to given value.


### GetUrl

`func (o *DeploymentResponseBody) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *DeploymentResponseBody) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *DeploymentResponseBody) SetUrl(v string)`

SetUrl sets Url field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


