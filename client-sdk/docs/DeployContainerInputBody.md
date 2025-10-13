# DeployContainerInputBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Schema** | Pointer to **string** | A URL to the JSON Schema for this object. | [optional] [readonly] 
**ImageName** | **string** |  | 
**RegistryAuthToken** | Pointer to **string** |  | [optional] 
**RegistryUrl** | **string** |  | 
**Url** | **string** |  | 

## Methods

### NewDeployContainerInputBody

`func NewDeployContainerInputBody(imageName string, registryUrl string, url string, ) *DeployContainerInputBody`

NewDeployContainerInputBody instantiates a new DeployContainerInputBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeployContainerInputBodyWithDefaults

`func NewDeployContainerInputBodyWithDefaults() *DeployContainerInputBody`

NewDeployContainerInputBodyWithDefaults instantiates a new DeployContainerInputBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSchema

`func (o *DeployContainerInputBody) GetSchema() string`

GetSchema returns the Schema field if non-nil, zero value otherwise.

### GetSchemaOk

`func (o *DeployContainerInputBody) GetSchemaOk() (*string, bool)`

GetSchemaOk returns a tuple with the Schema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchema

`func (o *DeployContainerInputBody) SetSchema(v string)`

SetSchema sets Schema field to given value.

### HasSchema

`func (o *DeployContainerInputBody) HasSchema() bool`

HasSchema returns a boolean if a field has been set.

### GetImageName

`func (o *DeployContainerInputBody) GetImageName() string`

GetImageName returns the ImageName field if non-nil, zero value otherwise.

### GetImageNameOk

`func (o *DeployContainerInputBody) GetImageNameOk() (*string, bool)`

GetImageNameOk returns a tuple with the ImageName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageName

`func (o *DeployContainerInputBody) SetImageName(v string)`

SetImageName sets ImageName field to given value.


### GetRegistryAuthToken

`func (o *DeployContainerInputBody) GetRegistryAuthToken() string`

GetRegistryAuthToken returns the RegistryAuthToken field if non-nil, zero value otherwise.

### GetRegistryAuthTokenOk

`func (o *DeployContainerInputBody) GetRegistryAuthTokenOk() (*string, bool)`

GetRegistryAuthTokenOk returns a tuple with the RegistryAuthToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistryAuthToken

`func (o *DeployContainerInputBody) SetRegistryAuthToken(v string)`

SetRegistryAuthToken sets RegistryAuthToken field to given value.

### HasRegistryAuthToken

`func (o *DeployContainerInputBody) HasRegistryAuthToken() bool`

HasRegistryAuthToken returns a boolean if a field has been set.

### GetRegistryUrl

`func (o *DeployContainerInputBody) GetRegistryUrl() string`

GetRegistryUrl returns the RegistryUrl field if non-nil, zero value otherwise.

### GetRegistryUrlOk

`func (o *DeployContainerInputBody) GetRegistryUrlOk() (*string, bool)`

GetRegistryUrlOk returns a tuple with the RegistryUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistryUrl

`func (o *DeployContainerInputBody) SetRegistryUrl(v string)`

SetRegistryUrl sets RegistryUrl field to given value.


### GetUrl

`func (o *DeployContainerInputBody) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *DeployContainerInputBody) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *DeployContainerInputBody) SetUrl(v string)`

SetUrl sets Url field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


