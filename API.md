# Survey Service API
This page documents the Survey service API endpoints. These endpoints will be secured using HTTP basic authentication initially. All endpoints return an `HTTP 200 OK` status code except where noted otherwise.

## List Surveys
* `GET /surveys` will return a list of known surveys.

### Example JSON Response
```json
[{
  "id": "cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87",
  "name": "BRES"
}]
```

An `HTTP 204 No Content` status code is returned if there are no known surveys.

## Get Survey
* `GET /surveys/cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87` will return the details of the survey with an ID of `cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87`.

### Example JSON Response
```json
{
  "id": "cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87",
  "name": "BRES"
}
```

An `HTTP 404 Not Found` status code is returned if the survey with the specified ID could not be found.

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