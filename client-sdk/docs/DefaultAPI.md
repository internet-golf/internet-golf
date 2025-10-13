# \DefaultAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetAlive**](DefaultAPI.md#GetAlive) | **Get** /alive | Get alive
[**GetDeploymentByUrl**](DefaultAPI.md#GetDeploymentByUrl) | **Get** /deployment/{url} | Get deployment by URL
[**PostTokenGenerate**](DefaultAPI.md#PostTokenGenerate) | **Post** /token/generate | Post token generate
[**PutDeployContainer**](DefaultAPI.md#PutDeployContainer) | **Put** /deploy/container | Put deploy container
[**PutDeployFiles**](DefaultAPI.md#PutDeployFiles) | **Put** /deploy/files | Put deploy files
[**PutDeployInitByDomain**](DefaultAPI.md#PutDeployInitByDomain) | **Put** /deploy/init/{domain} | Put deploy init by domain
[**PutDeployNew**](DefaultAPI.md#PutDeployNew) | **Put** /deploy/new | Put deploy new
[**PutUserRegister**](DefaultAPI.md#PutUserRegister) | **Put** /user/register | Put user register



## GetAlive

> HealthCheckOutputBody GetAlive(ctx).Execute()

Get alive

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
	resp, r, err := apiClient.DefaultAPI.GetAlive(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetAlive``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAlive`: HealthCheckOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetAlive`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetAliveRequest struct via the builder pattern


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


## GetDeploymentByUrl

> GetDeploymentOutputBody GetDeploymentByUrl(ctx, url).Execute()

Get deployment by URL

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
	resp, r, err := apiClient.DefaultAPI.GetDeploymentByUrl(context.Background(), url).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetDeploymentByUrl``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetDeploymentByUrl`: GetDeploymentOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetDeploymentByUrl`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**url** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetDeploymentByUrlRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetDeploymentOutputBody**](GetDeploymentOutputBody.md)

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


## PutDeployContainer

> SuccessOutputBody PutDeployContainer(ctx).DeployContainerInputBody(deployContainerInputBody).Execute()

Put deploy container

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
	deployContainerInputBody := *openapiclient.NewDeployContainerInputBody("ImageName_example", "RegistryUrl_example", "Url_example") // DeployContainerInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.PutDeployContainer(context.Background()).DeployContainerInputBody(deployContainerInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.PutDeployContainer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutDeployContainer`: SuccessOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.PutDeployContainer`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPutDeployContainerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deployContainerInputBody** | [**DeployContainerInputBody**](DeployContainerInputBody.md) |  | 

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


## PutDeployFiles

> SuccessOutputBody PutDeployFiles(ctx).Url(url).Contents(contents).KeepLeadingDirectories(keepLeadingDirectories).PreserveExistingFiles(preserveExistingFiles).Execute()

Put deploy files

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
	contents := os.NewFile(1234, "some_file") // *os.File |  (optional)
	keepLeadingDirectories := true // bool |  (optional)
	preserveExistingFiles := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.PutDeployFiles(context.Background()).Url(url).Contents(contents).KeepLeadingDirectories(keepLeadingDirectories).PreserveExistingFiles(preserveExistingFiles).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.PutDeployFiles``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutDeployFiles`: SuccessOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.PutDeployFiles`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPutDeployFilesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **url** | **string** |  | 
 **contents** | ***os.File** |  | 
 **keepLeadingDirectories** | **bool** |  | 
 **preserveExistingFiles** | **bool** |  | 

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


## PutDeployInitByDomain

> SuccessOutputBody PutDeployInitByDomain(ctx, domain).Execute()

Put deploy init by domain

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
	domain := "domain_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.PutDeployInitByDomain(context.Background(), domain).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.PutDeployInitByDomain``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutDeployInitByDomain`: SuccessOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.PutDeployInitByDomain`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**domain** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutDeployInitByDomainRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**SuccessOutputBody**](SuccessOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutDeployNew

> SuccessOutputBody PutDeployNew(ctx).DeploymentCreateInputBody(deploymentCreateInputBody).Execute()

Put deploy new

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
	deploymentCreateInputBody := *openapiclient.NewDeploymentCreateInputBody("Url_example") // DeploymentCreateInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.PutDeployNew(context.Background()).DeploymentCreateInputBody(deploymentCreateInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.PutDeployNew``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutDeployNew`: SuccessOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.PutDeployNew`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPutDeployNewRequest struct via the builder pattern


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

