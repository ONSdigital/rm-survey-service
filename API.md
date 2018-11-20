# Survey Service API
This page documents the Survey service API endpoints. Apart from the Service Information endpoint, all these endpoints are secured using HTTP basic authentication. All endpoints return an `HTTP 200 OK` status code except where noted otherwise.

## Service Information
* `GET /info` will return information about this service, collated from when it was last built.

### Example JSON Response
```json
{
  "name": "surveysvc",
  "version": "10.42.0",
  "origin": "git@github.com:ONSdigital/rm-survey-service.git",
  "commit": "c81fc1dc2155aed0fc201f2273333d3af75e10e0",
  "branch": "master",
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
   "legalBasis": "Statistics of Trade Act 1947"
}]
```

An `HTTP 204 No Content` status code is returned if there are no known surveys.

## List Surveys by Survey Type
* 'GET /surveys/surveytype/<type>'  Where type is one of Business,Social or Census
i.e
    * GET /surveys/surveytype/Business
    * GET /surveys/surveytype/Social
    * GET /surveys/surveytype/Census
    
Note: Although the endpoint is case insensitive for Type, Pascal case matches the database enumeration 
and so is preferred. i.e Business preferred over business or BUSINESS
    
### Example JSON Response    
```json
[{
  "id": "cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87",
  "shortName": "BRES",
   "longName": "Business Register and Employment Survey",
   "surveyRef": "221",
   "legalBasis": "Statistics of Trade Act 1947"
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
  "legalBasis": "Statistics of Trade Act 1947"
}
```

An `HTTP 404 Not Found` status code is returned if the survey with the specified ID could not be found.

## Get Survey by Short Name
* `GET /surveys/shortname/bres` will return the details of the survey with the short name `bres` (or `BRES`).

### Example JSON Response
```json
{
  "id": "cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87",
  "shortName": "BRES",
  "longName": "Business Register and Employment Survey",
  "surveyRef": "221",
  "legalBasis": "Statistics of Trade Act 1947"
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
  "legalBasis": "Statistics of Trade Act 1947"
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
