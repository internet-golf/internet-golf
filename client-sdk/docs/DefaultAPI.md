# \DefaultAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateAlias**](DefaultAPI.md#CreateAlias) | **Put** /deploy/alias | 
[**CreateDeployment**](DefaultAPI.md#CreateDeployment) | **Put** /deploy/new | 
[**DeployAdminDash**](DefaultAPI.md#DeployAdminDash) | **Put** /admin-dash | 
[**DeployFiles**](DefaultAPI.md#DeployFiles) | **Put** /deploy/files | 
[**GetDeployment**](DefaultAPI.md#GetDeployment) | **Get** /deployment/{url} | 
[**GetDeployments**](DefaultAPI.md#GetDeployments) | **Get** /deployments | 
[**HealthCheck**](DefaultAPI.md#HealthCheck) | **Get** /alive | 
[**PostTokenGenerate**](DefaultAPI.md#PostTokenGenerate) | **Post** /token/generate | Post token generate
[**PutUserRegister**](DefaultAPI.md#PutUserRegister) | **Put** /user/register | Put user register



## CreateAlias

> SuccessOutputBody CreateAlias(ctx).DeployAliasBody(deployAliasBody).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	deployAliasBody := *openapiclient.NewDeployAliasBody("mysite.mydomain.com") // DeployAliasBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.CreateAlias(context.Background()).DeployAliasBody(deployAliasBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.CreateAlias``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateAlias`: SuccessOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.CreateAlias`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateAliasRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deployAliasBody** | [**DeployAliasBody**](DeployAliasBody.md) |  | 

### Return type

[**SuccessOutputBody**](SuccessOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateDeployment

> SuccessOutputBody CreateDeployment(ctx).DeploymentCreateInputBody(deploymentCreateInputBody).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	deploymentCreateInputBody := *openapiclient.NewDeploymentCreateInputBody("mysite.mydomain.com") // DeploymentCreateInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.CreateDeployment(context.Background()).DeploymentCreateInputBody(deploymentCreateInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.CreateDeployment``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateDeployment`: SuccessOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.CreateDeployment`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateDeploymentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deploymentCreateInputBody** | [**DeploymentCreateInputBody**](DeploymentCreateInputBody.md) |  | 

### Return type

[**SuccessOutputBody**](SuccessOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeployAdminDash

> SuccessOutputBody DeployAdminDash(ctx).DeployAdminDashBody(deployAdminDashBody).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	deployAdminDashBody := *openapiclient.NewDeployAdminDashBody("dash.mydomain.com") // DeployAdminDashBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.DeployAdminDash(context.Background()).DeployAdminDashBody(deployAdminDashBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.DeployAdminDash``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeployAdminDash`: SuccessOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.DeployAdminDash`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDeployAdminDashRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deployAdminDashBody** | [**DeployAdminDashBody**](DeployAdminDashBody.md) |  | 

### Return type

[**SuccessOutputBody**](SuccessOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeployFiles

> SuccessOutputBody DeployFiles(ctx).Url(url).Contents(contents).KeepLeadingDirectories(keepLeadingDirectories).PreserveExistingFiles(preserveExistingFiles).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	url := "url_example" // string | The URL of the deployment that you're updating.
	contents := os.NewFile(1234, "some_file") // *os.File | A .tar.gz that contains the files to be deployed. (optional)
	keepLeadingDirectories := true // bool | By default, if you upload a .tar.gz whose contents are all in one folder, the contents of that folder will be used instead of the folder itself. For example, if you upload a folder called 'dist' for the deployment 'mysite.com', the URL of your site content will not be at 'mysite.com/dist'. Setting this to true turns off that auto-unpacking. (optional) (default to false)
	preserveExistingFiles := true // bool | Leave the existing files for the current deployment in place instead of completely replacing them. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.DeployFiles(context.Background()).Url(url).Contents(contents).KeepLeadingDirectories(keepLeadingDirectories).PreserveExistingFiles(preserveExistingFiles).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.DeployFiles``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeployFiles`: SuccessOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.DeployFiles`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDeployFilesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **url** | **string** | The URL of the deployment that you&#39;re updating. | 
 **contents** | ***os.File** | A .tar.gz that contains the files to be deployed. | 
 **keepLeadingDirectories** | **bool** | By default, if you upload a .tar.gz whose contents are all in one folder, the contents of that folder will be used instead of the folder itself. For example, if you upload a folder called &#39;dist&#39; for the deployment &#39;mysite.com&#39;, the URL of your site content will not be at &#39;mysite.com/dist&#39;. Setting this to true turns off that auto-unpacking. | [default to false]
 **preserveExistingFiles** | **bool** | Leave the existing files for the current deployment in place instead of completely replacing them. | 

### Return type

[**SuccessOutputBody**](SuccessOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDeployment

> GetDeployment200Response GetDeployment(ctx, url).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	url := "url_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetDeployment(context.Background(), url).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetDeployment``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetDeployment`: GetDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetDeployment`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**url** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetDeploymentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetDeployment200Response**](GetDeployment200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDeployments

> GetDeployments200Response GetDeployments(ctx).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetDeployments(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetDeployments``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetDeployments`: GetDeployments200Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetDeployments`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetDeploymentsRequest struct via the builder pattern


### Return type

[**GetDeployments200Response**](GetDeployments200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## HealthCheck

> HealthCheckOutputBody HealthCheck(ctx).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.HealthCheck(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.HealthCheck``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `HealthCheck`: HealthCheckOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.HealthCheck`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiHealthCheckRequest struct via the builder pattern


### Return type

[**HealthCheckOutputBody**](HealthCheckOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PostTokenGenerate

> CreateBearerTokenOutputBody PostTokenGenerate(ctx).CreateBearerTokenInputBody(createBearerTokenInputBody).Execute()

Post token generate

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	createBearerTokenInputBody := *openapiclient.NewCreateBearerTokenInputBody(false) // CreateBearerTokenInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.PostTokenGenerate(context.Background()).CreateBearerTokenInputBody(createBearerTokenInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.PostTokenGenerate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PostTokenGenerate`: CreateBearerTokenOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.PostTokenGenerate`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPostTokenGenerateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createBearerTokenInputBody** | [**CreateBearerTokenInputBody**](CreateBearerTokenInputBody.md) |  | 

### Return type

[**CreateBearerTokenOutputBody**](CreateBearerTokenOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutUserRegister

> SuccessOutputBody PutUserRegister(ctx).AddExternalUserInputBody(addExternalUserInputBody).Execute()

Put user register

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	addExternalUserInputBody := *openapiclient.NewAddExternalUserInputBody("ExternalUserSource_example") // AddExternalUserInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.PutUserRegister(context.Background()).AddExternalUserInputBody(addExternalUserInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.PutUserRegister``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutUserRegister`: SuccessOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.PutUserRegister`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPutUserRegisterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addExternalUserInputBody** | [**AddExternalUserInputBody**](AddExternalUserInputBody.md) |  | 

### Return type

[**SuccessOutputBody**](SuccessOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

