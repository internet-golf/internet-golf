# ResponseIfEmpty

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ExternalSource** | Pointer to **string** | Original repository for this deployment&#39;s source. Can include a branch name. | [optional] 
**ExternalSourceType** | Pointer to **string** | Place where the original repository lives. | [optional] 
**NoContentYet** | **bool** | Set to true if this deployment has not yet been set up. | 
**PreserveExternalPath** | Pointer to **bool** | if this is true and the deployment url has a path like \&quot;/thing\&quot;, then the \&quot;/thing\&quot; in the path will be transparently passed through to the underlying resource instead of being removed (which is the default) | [optional] 
**Tags** | Pointer to **[]string** | Tags used for metadata. | [optional] 
**Type** | **string** | Type of deployment contents. | 
**Url** | **string** | URL that this deployment will appear at. The DNS for the domain has to be set up first. | 

## Methods

### NewResponseIfEmpty

`func NewResponseIfEmpty(noContentYet bool, type_ string, url string, ) *ResponseIfEmpty`

NewResponseIfEmpty instantiates a new ResponseIfEmpty object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewResponseIfEmptyWithDefaults

`func NewResponseIfEmptyWithDefaults() *ResponseIfEmpty`

NewResponseIfEmptyWithDefaults instantiates a new ResponseIfEmpty object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetExternalSource

`func (o *ResponseIfEmpty) GetExternalSource() string`

GetExternalSource returns the ExternalSource field if non-nil, zero value otherwise.

### GetExternalSourceOk

`func (o *ResponseIfEmpty) GetExternalSourceOk() (*string, bool)`

GetExternalSourceOk returns a tuple with the ExternalSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSource

`func (o *ResponseIfEmpty) SetExternalSource(v string)`

SetExternalSource sets ExternalSource field to given value.

### HasExternalSource

`func (o *ResponseIfEmpty) HasExternalSource() bool`

HasExternalSource returns a boolean if a field has been set.

### GetExternalSourceType

`func (o *ResponseIfEmpty) GetExternalSourceType() string`

GetExternalSourceType returns the ExternalSourceType field if non-nil, zero value otherwise.

### GetExternalSourceTypeOk

`func (o *ResponseIfEmpty) GetExternalSourceTypeOk() (*string, bool)`

GetExternalSourceTypeOk returns a tuple with the ExternalSourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSourceType

`func (o *ResponseIfEmpty) SetExternalSourceType(v string)`

SetExternalSourceType sets ExternalSourceType field to given value.

### HasExternalSourceType

`func (o *ResponseIfEmpty) HasExternalSourceType() bool`

HasExternalSourceType returns a boolean if a field has been set.

### GetNoContentYet

`func (o *ResponseIfEmpty) GetNoContentYet() bool`

GetNoContentYet returns the NoContentYet field if non-nil, zero value otherwise.

### GetNoContentYetOk

`func (o *ResponseIfEmpty) GetNoContentYetOk() (*bool, bool)`

GetNoContentYetOk returns a tuple with the NoContentYet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNoContentYet

`func (o *ResponseIfEmpty) SetNoContentYet(v bool)`

SetNoContentYet sets NoContentYet field to given value.


### GetPreserveExternalPath

`func (o *ResponseIfEmpty) GetPreserveExternalPath() bool`

GetPreserveExternalPath returns the PreserveExternalPath field if non-nil, zero value otherwise.

### GetPreserveExternalPathOk

`func (o *ResponseIfEmpty) GetPreserveExternalPathOk() (*bool, bool)`

GetPreserveExternalPathOk returns a tuple with the PreserveExternalPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreserveExternalPath

`func (o *ResponseIfEmpty) SetPreserveExternalPath(v bool)`

SetPreserveExternalPath sets PreserveExternalPath field to given value.

### HasPreserveExternalPath

`func (o *ResponseIfEmpty) HasPreserveExternalPath() bool`

HasPreserveExternalPath returns a boolean if a field has been set.

### GetTags

`func (o *ResponseIfEmpty) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *ResponseIfEmpty) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *ResponseIfEmpty) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *ResponseIfEmpty) HasTags() bool`

HasTags returns a boolean if a field has been set.

### SetTagsNil

`func (o *ResponseIfEmpty) SetTagsNil(b bool)`

 SetTagsNil sets the value for Tags to be an explicit nil

### UnsetTags
`func (o *ResponseIfEmpty) UnsetTags()`

UnsetTags ensures that no value is present for Tags, not even an explicit nil
### GetType

`func (o *ResponseIfEmpty) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ResponseIfEmpty) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ResponseIfEmpty) SetType(v string)`

SetType sets Type field to given value.


### GetUrl

`func (o *ResponseIfEmpty) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *ResponseIfEmpty) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *ResponseIfEmpty) SetUrl(v string)`

SetUrl sets Url field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


