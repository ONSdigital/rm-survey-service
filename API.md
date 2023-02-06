# Survey Service API
This page documents the Survey service API endpoints. Apart from the Service Information endpoint, all these endpoints are secured using HTTP basic authentication. All endpoints return an `HTTP 200 OK` status code except where noted otherwise.

## Service Information
* `GET /info` will return information about this service, collated from when it was last built.

### Example JSON Response
```json
{
  "name": "surveysvc",
  "version": "10.42.1",
  "origin": "git@github.com:ONSdigital/rm-survey-service.git",
  "commit": "c81fc1dc2155aed0fc201f2273333d3af75e10e0",
  "branch": "main",
  "built": "2017-07-05T18:47:28Z"
}
```

## List Surveys
* `GET /surveys` will return a list of known surveys.

### Example JSON Response
```json
[{
  "id": "cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87",
  "shortName": "BRES",
   "longName": "Business Register and Employment Survey",
   "surveyRef": "221",
   "legalBasis": "Statistics of Trade Act 1947",
   "surveyMode": "SEFT"
}]
```

An `HTTP 204 No Content` status code is returned if there are no known surveys.

## List Surveys by Survey Type
*   'GET /surveys/surveytype/<type>' Returns a list of surveys of a specific type. Type is one of Business,Social or Census. Although the endpoint is case insensitive for <Type>, Pascal case matches the database enumeration and so is preferred. i.e Business preferred over business or BUSINESS
    
### Example JSON Response    
```json
[{
  "id": "cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87",
  "shortName": "BRES",
   "longName": "Business Register and Employment Survey",
   "surveyRef": "221",
   "legalBasis": "Statistics of Trade Act 1947",
   "surveyMode": "SEFT"
}]
```

An `HTTP 204 No Content` status code is returned if there are no known surveys.

## Get Survey
* `GET /surveys/cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87` will return the details of the survey with an ID of `cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87`.

### Example JSON Response
```json
{
  "id": "cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87",
  "shortName": "BRES",
   "longName": "Business Register and Employment Survey",
   "surveyRef": "221",
   "legalBasis": "Statistics of Trade Act 1947",
   "surveyMode": "SEFT"
}
```

An `HTTP 404 Not Found` status code is returned if the survey with the specified ID could not be found.\

## Delete Survey
* `DELETE /surveys/<survey-id>` will delete the survey with the matching id, and also all the classifiers

- Returns 204 on success
- Returns 400 if the id isn't in the correct format
- Returns 403 if the http authentication isn't correct
- Returns 404 if the id of the survey isn't found

## Get Survey by Short Name
* `GET /surveys/shortname/bres` will return the details of the survey with the short name `bres` (or `BRES`).

### Example JSON Response
```json
{
  "id": "cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87",
  "shortName": "BRES",
  "longName": "Business Register and Employment Survey",
  "surveyRef": "221",
  "legalBasis": "Statistics of Trade Act 1947",
  "surveyMode": "SEFT"
}
```

An `HTTP 404 Not Found` status code is returned if the survey with the specified short name could not be found.

## Get Survey by Reference
* `GET /surveys/ref/221` will return the details of the survey with the reference `221`.

### Example JSON Response
```json
{
  "id": "cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87",
  "shortName": "BRES",
  "longName": "Business Register and Employment Survey",
  "surveyRef": "221",
  "legalBasis": "Statistics of Trade Act 1947",
  "surveyMode": "SEFT"
}
```

An `HTTP 404 Not Found` status code is returned if the survey with the specified reference could not be found.

## List Classifier Type Selectors
* `GET /surveys/cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87/classifiertypeselectors` will return a list of classifier type selectors for the survey with an ID of `cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87`.

### Example JSON Response
```json
[
  {
    "id": "efa868fb-fb80-44c7-9f33-d6800a17c4da",
    "name": "COLLECTION_INSTRUMENT"
  },
  {
    "id": "e119ffd6-6fc1-426c-ae81-67a96f9a71ba",
    "name": "COMMUNICATION_TEMPLATE"
  }
]
```

An `HTTP 404 Not Found` status code is returned if the survey with the specified ID could not be found. An `HTTP 204 No Content` status code is returned if there are no classifier type selectors for the survey with the specified ID.

## Get Classifier Types Selector
* `GET /surveys/cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87/classifiertypeselectors/efa868fb-fb80-44c7-9f33-d6800a17c4da` will return the details of the classifier type selector with an ID of `efa868fb-fb80-44c7-9f33-d6800a17c4da` for the survey with an ID of `cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87`.

### Example JSON Response
```json
{
  "id": "efa868fb-fb80-44c7-9f33-d6800a17c4da",
  "name": "COLLECTION_INSTRUMENT",
  "classifierTypes": [
    "COLLECTION_EXERCISE",
    "RU_REF"
  ]
}
```

An `HTTP 404 Not Found` status code is returned if the survey or classifier type selector with the specified ID could not be found.

## Post Survey Classifiers
* `POST /surveys/<survey_id>/classifiers`

The payload should be a classifier object, with a classifier type selector `name` and a list of `classifierTypes` as strings.

### Example JSON payload
```json

{
    "name": "COLLECTION_INSTRUMENT",
    "classifierTypes": [
      "FORM_TYPE",
      "LEGAL_BASIS"
    ]
}
```

An `HTTP 404 Not Found` status code is returned if the survey with the specified ID could not be found.

An `HTTP 409 Conflict` status code is returned if a classifier type selector already exists for any of the names in the payload.

## Post New Survey
* `POST /surveys` will create a new survey.

The payload should be a JSON document, with an `id`, a `shortName`, a `longName`, a `surveyRef`, a `legalBasis`, a `surveyType`, and a `legalBasisRef` as strings, and `classifiers` as a list.

### Example JSON payload
```json
{
    "id": "efa868fb-fb80-44c7-9f33-d6800a17c4da",
    "shortName": "test-short-name", 
    "longName": "test-long-name",
    "surveyRef": "456",
    "legalBasis": "Statistics of Trade Act 1947",
    "surveyType": "Social",
    "surveyMode": "SEFT",
    "legalBasisRef": "STA1947",
    "classifiers": [
      "LEGAL_BASIS"
    ]
}
```

An `HTTP 400 Bad Request` status code is returned if the payload has missing values and is incomplete.

## Put Survey Details on Reference
* `PUT /surveys/ref/456` will put details about a survey at a specific reference number, in this case 456.

The payload should be a JSON document, with an `id`, a `shortName`, a `longName`, a `surveyRef`, a `legalBasis`, a `surveyType`, and a `legalBasisRef` as strings, and `classifiers` as a list.

### Example JSON payload
```json
{
    "id": "efa868fb-fb80-44c7-9f33-d6800a17c4da",
    "shortName": "test-short-name", 
    "longName": "test-long-name",
    "surveyRef": "456",
    "legalBasis": "Statistics of Trade Act 1947",
    "surveyType": "Social",
    "surveyMode": "SEFT",
    "legalBasisRef": "STA1947",
    "classifiers": [
      "LEGAL_BASIS"
    ]
}
```

An `HTTP 500 Internal Server Error` status code is returned if the PUT request was unsuccessful.

## Get Legal Bases
* `GET /legal-bases` returns a list of legal bases.

### Example JSON payload
```json
[
    {"ref":"GovERD","longName":"GovERD"},
    {"ref":"STA1947","longName":"Statistics of Trade Act 1947"},
    {"ref":"STA1947_BEIS","longName":"Statistics of Trade Act 1947 - BEIS"},
    {"ref":"Vol","longName":"Voluntary Not Stated"},
    {"ref":"Vol_BEIS","longName":"Voluntary - BEIS"}
]
```