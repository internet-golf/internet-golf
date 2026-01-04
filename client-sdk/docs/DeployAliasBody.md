# DeployAliasBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**Url** | **string** | The URL of the deployment that you&#39;re updating. | 
**AliasedTo** | Pointer to **string** | The URL that this deployment is an alias for. | [optional] 
**Redirect** | Pointer to **bool** | If this is true, visitors to this deployment&#39;s URL will be completely redirected to the URL that this alias is for. | [optional] 

## Methods

### NewDeployAliasBody

`func NewDeployAliasBody(url string, ) *DeployAliasBody`

NewDeployAliasBody instantiates a new DeployAliasBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeployAliasBodyWithDefaults

`func NewDeployAliasBodyWithDefaults() *DeployAliasBody`

NewDeployAliasBodyWithDefaults instantiates a new DeployAliasBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *DeployAliasBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *DeployAliasBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *DeployAliasBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *DeployAliasBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetUrl

`func (o *DeployAliasBody) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *DeployAliasBody) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *DeployAliasBody) SetUrl(v string)`

SetUrl sets Url field to given value.


### GetAliasedTo

`func (o *DeployAliasBody) GetAliasedTo() string`

GetAliasedTo returns the AliasedTo field if non-nil, zero value otherwise.

### GetAliasedToOk

`func (o *DeployAliasBody) GetAliasedToOk() (*string, bool)`

GetAliasedToOk returns a tuple with the AliasedTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAliasedTo

`func (o *DeployAliasBody) SetAliasedTo(v string)`

SetAliasedTo sets AliasedTo field to given value.

### HasAliasedTo

`func (o *DeployAliasBody) HasAliasedTo() bool`

HasAliasedTo returns a boolean if a field has been set.

### GetRedirect

`func (o *DeployAliasBody) GetRedirect() bool`

GetRedirect returns the Redirect field if non-nil, zero value otherwise.

### GetRedirectOk

`func (o *DeployAliasBody) GetRedirectOk() (*bool, bool)`

GetRedirectOk returns a tuple with the Redirect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirect

`func (o *DeployAliasBody) SetRedirect(v bool)`

SetRedirect sets Redirect field to given value.

### HasRedirect

`func (o *DeployAliasBody) HasRedirect() bool`

HasRedirect returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


