# GetDeployment200Response

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
**AliasedTo** | Pointer to **string** | The URL that this deployment is an alias for. | [optional] 
**Redirect** | Pointer to **bool** | If this is true, visitors to this deployment&#39;s URL will be completely redirected to the URL that this alias is for. | [optional] 
**NoContentYet** | Pointer to **bool** | Set to true to indicate that this deployment has not yet been set up. | [optional] 

## Methods

### NewGetDeployment200Response

`func NewGetDeployment200Response(createdAt string, meta SiteMeta, name string, type_ string, updatedAt string, url string, ) *GetDeployment200Response`

NewGetDeployment200Response instantiates a new GetDeployment200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetDeployment200ResponseWithDefaults

`func NewGetDeployment200ResponseWithDefaults() *GetDeployment200Response`

NewGetDeployment200ResponseWithDefaults instantiates a new GetDeployment200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCreatedAt

`func (o *GetDeployment200Response) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *GetDeployment200Response) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *GetDeployment200Response) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.


### GetExternalSource

`func (o *GetDeployment200Response) GetExternalSource() string`

GetExternalSource returns the ExternalSource field if non-nil, zero value otherwise.

### GetExternalSourceOk

`func (o *GetDeployment200Response) GetExternalSourceOk() (*string, bool)`

GetExternalSourceOk returns a tuple with the ExternalSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSource

`func (o *GetDeployment200Response) SetExternalSource(v string)`

SetExternalSource sets ExternalSource field to given value.

### HasExternalSource

`func (o *GetDeployment200Response) HasExternalSource() bool`

HasExternalSource returns a boolean if a field has been set.

### GetExternalSourceType

`func (o *GetDeployment200Response) GetExternalSourceType() string`

GetExternalSourceType returns the ExternalSourceType field if non-nil, zero value otherwise.

### GetExternalSourceTypeOk

`func (o *GetDeployment200Response) GetExternalSourceTypeOk() (*string, bool)`

GetExternalSourceTypeOk returns a tuple with the ExternalSourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalSourceType

`func (o *GetDeployment200Response) SetExternalSourceType(v string)`

SetExternalSourceType sets ExternalSourceType field to given value.

### HasExternalSourceType

`func (o *GetDeployment200Response) HasExternalSourceType() bool`

HasExternalSourceType returns a boolean if a field has been set.

### GetMeta

`func (o *GetDeployment200Response) GetMeta() SiteMeta`

GetMeta returns the Meta field if non-nil, zero value otherwise.

### GetMetaOk

`func (o *GetDeployment200Response) GetMetaOk() (*SiteMeta, bool)`

GetMetaOk returns a tuple with the Meta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMeta

`func (o *GetDeployment200Response) SetMeta(v SiteMeta)`

SetMeta sets Meta field to given value.


### GetName

`func (o *GetDeployment200Response) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *GetDeployment200Response) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *GetDeployment200Response) SetName(v string)`

SetName sets Name field to given value.


### GetPreserveExternalPath

`func (o *GetDeployment200Response) GetPreserveExternalPath() bool`

GetPreserveExternalPath returns the PreserveExternalPath field if non-nil, zero value otherwise.

### GetPreserveExternalPathOk

`func (o *GetDeployment200Response) GetPreserveExternalPathOk() (*bool, bool)`

GetPreserveExternalPathOk returns a tuple with the PreserveExternalPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreserveExternalPath

`func (o *GetDeployment200Response) SetPreserveExternalPath(v bool)`

SetPreserveExternalPath sets PreserveExternalPath field to given value.

### HasPreserveExternalPath

`func (o *GetDeployment200Response) HasPreserveExternalPath() bool`

HasPreserveExternalPath returns a boolean if a field has been set.

### GetServerContentLocation

`func (o *GetDeployment200Response) GetServerContentLocation() string`

GetServerContentLocation returns the ServerContentLocation field if non-nil, zero value otherwise.

### GetServerContentLocationOk

`func (o *GetDeployment200Response) GetServerContentLocationOk() (*string, bool)`

GetServerContentLocationOk returns a tuple with the ServerContentLocation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServerContentLocation

`func (o *GetDeployment200Response) SetServerContentLocation(v string)`

SetServerContentLocation sets ServerContentLocation field to given value.

### HasServerContentLocation

`func (o *GetDeployment200Response) HasServerContentLocation() bool`

HasServerContentLocation returns a boolean if a field has been set.

### GetSpaMode

`func (o *GetDeployment200Response) GetSpaMode() bool`

GetSpaMode returns the SpaMode field if non-nil, zero value otherwise.

### GetSpaModeOk

`func (o *GetDeployment200Response) GetSpaModeOk() (*bool, bool)`

GetSpaModeOk returns a tuple with the SpaMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpaMode

`func (o *GetDeployment200Response) SetSpaMode(v bool)`

SetSpaMode sets SpaMode field to given value.

### HasSpaMode

`func (o *GetDeployment200Response) HasSpaMode() bool`

HasSpaMode returns a boolean if a field has been set.

### GetTags

`func (o *GetDeployment200Response) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *GetDeployment200Response) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *GetDeployment200Response) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *GetDeployment200Response) HasTags() bool`

HasTags returns a boolean if a field has been set.

### SetTagsNil

`func (o *GetDeployment200Response) SetTagsNil(b bool)`

 SetTagsNil sets the value for Tags to be an explicit nil

### UnsetTags
`func (o *GetDeployment200Response) UnsetTags()`

UnsetTags ensures that no value is present for Tags, not even an explicit nil
### GetType

`func (o *GetDeployment200Response) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *GetDeployment200Response) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *GetDeployment200Response) SetType(v string)`

SetType sets Type field to given value.


### GetUpdatedAt

`func (o *GetDeployment200Response) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *GetDeployment200Response) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *GetDeployment200Response) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetUrl

`func (o *GetDeployment200Response) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *GetDeployment200Response) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *GetDeployment200Response) SetUrl(v string)`

SetUrl sets Url field to given value.


### GetAliasedTo

`func (o *GetDeployment200Response) GetAliasedTo() string`

GetAliasedTo returns the AliasedTo field if non-nil, zero value otherwise.

### GetAliasedToOk

`func (o *GetDeployment200Response) GetAliasedToOk() (*string, bool)`

GetAliasedToOk returns a tuple with the AliasedTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAliasedTo

`func (o *GetDeployment200Response) SetAliasedTo(v string)`

SetAliasedTo sets AliasedTo field to given value.

### HasAliasedTo

`func (o *GetDeployment200Response) HasAliasedTo() bool`

HasAliasedTo returns a boolean if a field has been set.

### GetRedirect

`func (o *GetDeployment200Response) GetRedirect() bool`

GetRedirect returns the Redirect field if non-nil, zero value otherwise.

### GetRedirectOk

`func (o *GetDeployment200Response) GetRedirectOk() (*bool, bool)`

GetRedirectOk returns a tuple with the Redirect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirect

`func (o *GetDeployment200Response) SetRedirect(v bool)`

SetRedirect sets Redirect field to given value.

### HasRedirect

`func (o *GetDeployment200Response) HasRedirect() bool`

HasRedirect returns a boolean if a field has been set.

### GetNoContentYet

`func (o *GetDeployment200Response) GetNoContentYet() bool`

GetNoContentYet returns the NoContentYet field if non-nil, zero value otherwise.

### GetNoContentYetOk

`func (o *GetDeployment200Response) GetNoContentYetOk() (*bool, bool)`

GetNoContentYetOk returns a tuple with the NoContentYet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNoContentYet

`func (o *GetDeployment200Response) SetNoContentYet(v bool)`

SetNoContentYet sets NoContentYet field to given value.

### HasNoContentYet

`func (o *GetDeployment200Response) HasNoContentYet() bool`

HasNoContentYet returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


