# ranna main API
The ranna main REST API.

## Version: 1.0

### /exec

#### POST
##### Summary

Get Spec Map

##### Description

Returns the available spec map.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| payload | body | The execution payload | Yes | [models.ExecutionRequest](#modelsexecutionrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.ExecutionResponse](#modelsexecutionresponse) |
| 400 | Bad Request | [models.ErrorModel](#modelserrormodel) |
| 500 | Internal Server Error | [models.ErrorModel](#modelserrormodel) |

### /info

#### GET
##### Summary

Get System Info

##### Description

Returns general system and version information.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.ExecutionResponse](#modelsexecutionresponse) |
| 500 | Internal Server Error | [models.ErrorModel](#modelserrormodel) |

### /spec

#### GET
##### Summary

Get Spec Map

##### Description

Returns the available spec map.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.SpecMap](#modelsspecmap) |

### Models

#### models.ErrorModel

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code | integer |  | No |
| context | string |  | No |
| error | string |  | No |

#### models.ExecutionRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| arguments | [ string ] |  | No |
| code | string |  | No |
| environment | object |  | No |
| inline_expression | boolean |  | No |
| language | string |  | No |

#### models.ExecutionResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| exectimems | integer |  | No |
| stderr | string |  | No |
| stdout | string |  | No |

#### models.InlineSpec

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| import_regex | string |  | No |
| template | string |  | No |

#### models.Spec

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| cmd | string |  | No |
| entrypoint | string |  | No |
| example | string |  | No |
| filename | string |  | No |
| image | string |  | No |
| inline | [models.InlineSpec](#modelsinlinespec) |  | No |
| registry | string |  | No |
| use | string |  | No |

#### models.SpecMap

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| models.SpecMap | object |  |  |
