# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Remote Procedure Calls for Pip.Services in Go Changelog

## <a name="1.5.2"></a> 1.5.2 (2023-01-12)
### Bug fixing
- Fixed https connection validation

## <a name="1.5.1"></a> 1.5.1 (2023-01-12)
### Features
- Update dependencies

## <a name="1.5.0"></a> 1.5.0 (2021-10-18)
### Features
* Added regexp supporting to interceptor
   Examples:
   - the interceptor route **"/dummies"** corresponds to all of this routes **"/dummies"**, **"/dummies/check"**, **"/dummies/test"**
   - the interceptor route **"/dummies$"** corresponds only for this route **"/dummies"**. The routes **"/dummies/check"**, **"/dummies/test"** aren't processing by interceptor
   Please, don't forgot, route in interceptor always automaticaly concateneted with base route, like this **service_base_route + route_in_interceptor**. 
   For example, "/api/v1/" - service base route, "/dummies$" - interceptor route, in result will be next expression - "/api/v1/dummies$"
## <a name="1.4.4"></a> 1.4.4 (2021-08-30)
### Bug fixing
* Fix retry mechnaism in REST client

## <a name="1.4.3"></a> 1.4.3 (2021-08-23)
### Bug fixing
* Updated error propagation mechanism between client and services

## <a name="1.4.2"></a> 1.4.2 (2021-07-30)
### Features
* Add configuration parameters for CORS Headers in HttpEndpoint. Use *cors_headers* and *cors_origins*.
  Example:
  ```yml
  -cors_headers: "correlation_id, access_token, Accept, Content-Type, Content-Length, X-CSRF-Token"
  -cors_origins:  "*"
  ```
## <a name="1.4.1"></a> 1.4.1 (2021-07-26)
### Bug fixing
- Fix route checks in interceptors

## <a name="1.4.0"></a> 1.4.0 (2021-07-20)
### Features
* Add methods for controll CORS Headers in HttpEndpoint
* Add configuration parameters for CORS Headers in HttpEndpoint

## <a name="1.3.3"></a> 1.3.3 (2021-06-08)
### Features
* Update Instruments and added tracers
* Fix loggers
## <a name="1.3.2"></a> 1.3.2 (2021-05-06)
### Features
* **test** Refactor test services running
* Encode URL params

## <a name="1.3.1"></a> 1.3.1 (2021-04-23) 

### Features
* Add InstrumentTiming 

## <a name="1.3.0"></a> 1.3.0 (2021-04-23) 

### Breaking Changes
* **test** Added TestRestClient
* **test** Added TestCommandableHttpClient

## <a name="1.2.0"></a> 1.2.0 (2021-04-04) 

### Breaking Changes
* Introduced IRpcServiceOverrides
* Changed signature NewRpcService to InheritRpcService
* Changed signature NewCommandableRpcService to InheritRpcService

## <a name="1.1.3"></a> 1.1.3 (2021-03-15)

### Features
* **services** Added **correlation_id** and **access_token** to CORS headers

## <a name="1.1.0"></a> 1.1.0 (2021-02-21)

### Features
* **services** Added integration with Swagger UI

## <a name="1.0.13"></a> 1.0.13 (2020-12-10) 

### Features
* Fix work with CorrelationID in RestService
* Update dependencies

## <a name="1.0.12"></a> 1.0.12 (2020-12-10) 

### Features
* Fix headers in  RestClient for properly work with others services 

## <a name="1.0.8-1.0.11"></a> 1.0.8-1.0.11 (2020-12-02) 

### Features
* Added helper methods to RestOperations
* Changed RegisterWithAuth methods

### Bug Fixes
* Fix authorizer

## <a name="1.0.7"></a> 1.0.7 (2020-11-20) 

### Features
* Added swagger support

## <a name="1.0.5-1.0.6"></a> 1.0.5-1.0.6 (2020-11-13) 

### Features
* Added helper methods

## <a name="1.0.3-1.0.4"></a> 1.0.3-1.0.4 (2020-11-12) 

### Features
* Added helper methods in RestService

### Bug Fixes
* Fix signature CallCommand in CommandableHttpClient

## <a name="1.0.1-1.0.2"></a> 1.0.1-1.0.2 (2020-08-05) 

### Features
* Added error handler in Call method of RestClient

### Bug Fixes
* Fix response error method

## <a name="1.0.0"></a> 1.0.0 (2020-01-28) 

Initial public release

### Features
* **build** HTTP service factory
* **clients** mechanisms for retrieving connection settings
* **connect** helper module to retrieve connections services and clients
* **services** basic implementation of services for connecting

