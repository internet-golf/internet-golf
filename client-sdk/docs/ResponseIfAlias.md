# ResponseIfAlias

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AliasedTo** | Pointer to **string** | The URL that this deployment is an alias for. | [optional] 
**ExternalSource** | Pointer to **string** | Original repository for this deployment&#39;s source. Can include a branch name. | [optional] 
**ExternalSourceType** | Pointer to **string** | Place where the original repository lives. | [optional] 
**PreserveExternalPath** | Pointer to **bool** | if this is true and the deployment url has a path like \&quot;/thing\&quot;, then the \&quot;/thing\&quot; in the path will be transparently passed through to the underlying resource instead of being removed (which is the default) | [optional] 
**Redirect** | Pointer to **bool** | If this is true, visitors to this deployment&#39;s URL will be completely redirected to the URL that this alias is for. | [optional] 
**Tags** | Pointer to **[]string** | Tags used for metadata. | [optional] 
**Type** | **string** | Type of deployment contents. | 
**Url** | **string** | URL that this deployment will appear at. The DNS for the domain has to be set up first. | 

## Methods

### NewResponseIfAlias

`func NewResponseIfAlias(type_ string, url string, ) *ResponseIfAlias`

NewResponseIfAlias instantiates a new ResponseIfAlias object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewResponseIfAliasWithDefaults

`func NewResponseIfAliasWithDefaults() *ResponseIfAlias`

NewResponseIfAliasWithDefaults instantiates a new ResponseIfAlias object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAliasedTo

`func (o *ResponseIfAlias) GetAliasedTo() string`

GetAliasedTo returns the AliasedTo field if non-nil, zero value otherwise.

### GetAliasedToOk

`func (o *ResponseIfAlias) GetAliasedToOk() (*string, bool)`

GetAliasedToOk returns a tuple with the AliasedTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAliasedTo

`func (o *ResponseIfAlias) SetAliasedTo(v string)`

SetAliasedTo sets AliasedTo field to given value.

### HasAliasedTo

`func (o *ResponseIfAlias) HasAliasedTo() bool`

HasAliasedTo returns a boolean if a field has been set.

### GetExternalSource

`func (o *ResponseIfAlias) GetExternalSource() string`

GetExternalSource returns the ExternalSource field if non-nil, zero value otherwise.

### GetExternalSourceOk

`func (o *ResponseIfAlias) GetExternalSourceOk() (*string, bool)`

GetExternalSourceOk returns a tuple with the ExternalSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSource

`func (o *ResponseIfAlias) SetExternalSource(v string)`

SetExternalSource sets ExternalSource field to given value.

### HasExternalSource

`func (o *ResponseIfAlias) HasExternalSource() bool`

HasExternalSource returns a boolean if a field has been set.

### GetExternalSourceType

`func (o *ResponseIfAlias) GetExternalSourceType() string`

GetExternalSourceType returns the ExternalSourceType field if non-nil, zero value otherwise.

### GetExternalSourceTypeOk

`func (o *ResponseIfAlias) GetExternalSourceTypeOk() (*string, bool)`

GetExternalSourceTypeOk returns a tuple with the ExternalSourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSourceType

`func (o *ResponseIfAlias) SetExternalSourceType(v string)`

SetExternalSourceType sets ExternalSourceType field to given value.

### HasExternalSourceType

`func (o *ResponseIfAlias) HasExternalSourceType() bool`

HasExternalSourceType returns a boolean if a field has been set.

### GetPreserveExternalPath

`func (o *ResponseIfAlias) GetPreserveExternalPath() bool`

GetPreserveExternalPath returns the PreserveExternalPath field if non-nil, zero value otherwise.

### GetPreserveExternalPathOk

`func (o *ResponseIfAlias) GetPreserveExternalPathOk() (*bool, bool)`

GetPreserveExternalPathOk returns a tuple with the PreserveExternalPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreserveExternalPath

`func (o *ResponseIfAlias) SetPreserveExternalPath(v bool)`

SetPreserveExternalPath sets PreserveExternalPath field to given value.

### HasPreserveExternalPath

`func (o *ResponseIfAlias) HasPreserveExternalPath() bool`

HasPreserveExternalPath returns a boolean if a field has been set.

### GetRedirect

`func (o *ResponseIfAlias) GetRedirect() bool`

GetRedirect returns the Redirect field if non-nil, zero value otherwise.

### GetRedirectOk

`func (o *ResponseIfAlias) GetRedirectOk() (*bool, bool)`

GetRedirectOk returns a tuple with the Redirect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirect

`func (o *ResponseIfAlias) SetRedirect(v bool)`

SetRedirect sets Redirect field to given value.

### HasRedirect

`func (o *ResponseIfAlias) HasRedirect() bool`

HasRedirect returns a boolean if a field has been set.

### GetTags

`func (o *ResponseIfAlias) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *ResponseIfAlias) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *ResponseIfAlias) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *ResponseIfAlias) HasTags() bool`

HasTags returns a boolean if a field has been set.

### SetTagsNil

`func (o *ResponseIfAlias) SetTagsNil(b bool)`

 SetTagsNil sets the value for Tags to be an explicit nil

### UnsetTags
`func (o *ResponseIfAlias) UnsetTags()`

UnsetTags ensures that no value is present for Tags, not even an explicit nil
### GetType

`func (o *ResponseIfAlias) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ResponseIfAlias) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ResponseIfAlias) SetType(v string)`

SetType sets Type field to given value.


### GetUrl

`func (o *ResponseIfAlias) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *ResponseIfAlias) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *ResponseIfAlias) SetUrl(v string)`

SetUrl sets Url field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


