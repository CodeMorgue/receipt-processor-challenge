# Receipt Processor - Implementation Description

## System Requirements
This program assumes that the target system already has the Go programming language installed. Additionally, this program is best
run in a UNIX terminal.

## Running the Program
Open a UNIX terminal and navigate to the local directory housing this project. From the root of the directory, run the following command:
```bash
go run .
```
Upon success, the terminal should display a series of debug messages ending with the following line:
```bash
[GIN-debug] Listening and serving HTTP on localhost:8080
```
Now that the program is running and ready to handle requests, open a separate UNIX terminal and once again navigate to the project directory.
In your second terminal, paste the following POST request:
```bash
 curl http://localhost:8080/receipts/process \
--include \
--header "Content-Type: application/json" \
--request "POST" \
--data '{"retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}'
```
This submits a POST request for the provided json receipt to the receipts/process path and returns a UUID. The id will be shown in the following format:
```json
{"id": {id}}
```
In order to receive the point total for the receipt, the generated id must be passed to the path receipts/{id}/points. To do so, copy and paste the request below into the UNIX terminal being sure to replace {id} in the path with the id generated from the POST request:
```bash
 curl http://localhost:8080/receipts/{id}/points
 ```
This will make an implicit GET request and return the point total for the receipt associated with the provided id in the following format:
```json
{"points": {points}}
```
For this example receipt, the point total should be 28.

Below is the general format for writing a POST request
---
```bash
curl {url] \
--include \
--header "Content-Type: application/json" \
--request "POST" \
--data '{json data}'
```

## Testing the Program
To test the program, first be sure that the program is not already running. If the program is still running in a terminal, the process can be killed using CTRL-C. Then, in a UNIX terminal set at the root directory of the project, run the following command:
```bash
go test
```
The output should look as follows if all tests pass:
```bash
PASS
ok      example/ReceiptChallenge        0.730s
```
If a test does not pass, the terminal will display an error stating the expected value and the actual value that caused an error.
The testing file tests the .json receipts included in this project's /examples directory (morning-receipt.json and simple-receipt.json).
The object morning-receipt.json is expected to return 15 points, and simple-receipt.json is expected to return 31 points

## Additional Example Requests
Below is a series of POST requests that can be copied and pasted into the UNIX terminal for program testing.
```bash
 curl http://localhost:8080/receipts/process \
--include \
--header "Content-Type: application/json" \
--request "POST" \
--data '{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}'
```
```bash
 curl http://localhost:8080/receipts/process \
--include \
--header "Content-Type: application/json" \
--request "POST" \
--data '{
    "retailer": "Walgreens",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "08:13",
    "total": "2.65",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
        {"shortDescription": "Dasani", "price": "1.40"}
    ]
}'
```
```bash
 curl http://localhost:8080/receipts/process \
--include \
--header "Content-Type: application/json" \
--request "POST" \
--data '{
    "retailer": "Target",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "13:13",
    "total": "1.25",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"}
    ]
}'
```





# Receipt Processor - Challenge Description

Build a webservice that fulfils the documented API. The API is described below. A formal definition is provided 
in the [api.yml](./api.yml) file, but the information in this README is sufficient for completion of this challenge. We will use the 
described API to test your solution.

Provide any instructions required to run your application.

Data does not need to persist when your application stops. It is sufficient to store information in memory. There are too many different database solutions, we will not be installing a database on our system when testing your application.

## Language Selection

You can assume our engineers have Go and Docker installed to run your application. Go is our preferred language, but it is not a requirement for this exercise.

If you are using a language other than Go, the engineer evaluating your submission may not have an environment ready for your language. Your instructions should include how to get an environment in any OS that can run your project. For example, if you write your project in Javascript simply stating to "run `npm start` to start the application" is not sufficient, because the engineer may not have NPM. Providing a docker file and the required docker command is a simple way to satisfy this requirement.

## Submitting Your Solution

Provide a link to a public repository, such as GitHub or BitBucket, that contains your code to the provided link through Greenhouse.

---
## Summary of API Specification

### Endpoint: Process Receipts

* Path: `/receipts/process`
* Method: `POST`
* Payload: Receipt JSON
* Response: JSON containing an id for the receipt.

Description:

Takes in a JSON receipt (see example in the example directory) and returns a JSON object with an ID generated by your code.

The ID returned is the ID that should be passed into `/receipts/{id}/points` to get the number of points the receipt
was awarded.

How many points should be earned are defined by the rules below.

Reminder: Data does not need to survive an application restart. This is to allow you to use in-memory solutions to track any data generated by this endpoint.

Example Response:
```json
{ "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }
```

## Endpoint: Get Points

* Path: `/receipts/{id}/points`
* Method: `GET`
* Response: A JSON object containing the number of points awarded.

A simple Getter endpoint that looks up the receipt by the ID and returns an object specifying the points awarded.

Example Response:
```json
{ "points": 32 }
```

---

# Rules

These rules collectively define how many points should be awarded to a receipt.

* One point for every alphanumeric character in the retailer name.
* 50 points if the total is a round dollar amount with no cents.
* 25 points if the total is a multiple of `0.25`.
* 5 points for every two items on the receipt.
* If the trimmed length of the item description is a multiple of 3, multiply the price by `0.2` and round up to the nearest integer. The result is the number of points earned.
* 6 points if the day in the purchase date is odd.
* 10 points if the time of purchase is after 2:00pm and before 4:00pm.


## Examples

```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
```
```text
Total Points: 28
Breakdown:
     6 points - retailer name has 6 characters
    10 points - 4 items (2 pairs @ 5 points each)
     3 Points - "Emils Cheese Pizza" is 18 characters (a multiple of 3)
                item price of 12.25 * 0.2 = 2.45, rounded up is 3 points
     3 Points - "Klarbrunn 12-PK 12 FL OZ" is 24 characters (a multiple of 3)
                item price of 12.00 * 0.2 = 2.4, rounded up is 3 points
     6 points - purchase day is odd
  + ---------
  = 28 points
```

----

```json
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}
```
```text
Total Points: 109
Breakdown:
    50 points - total is a round dollar amount
    25 points - total is a multiple of 0.25
    14 points - retailer name (M&M Corner Market) has 14 alphanumeric characters
                note: '&' is not alphanumeric
    10 points - 2:33pm is between 2:00pm and 4:00pm
    10 points - 4 items (2 pairs @ 5 points each)
  + ---------
  = 109 points
```

---

# FAQ

### How will this exercise be evaluated?
An engineer will review the code you submit. At a minimum they must be able to run the service and the service must provide the expected results. You
should provide any necessary documentation within the repository. While your solution does not need to be fully production ready, you are being evaluated so
put your best foot forward.

### I have questions about the problem statement
For any requirements not specified via an example, use your best judgment to determine the expected result.

### Can I provide a private repository?
If at all possible, we prefer a public repository because we do not know which engineer will be evaluating your submission. Providing a public repository
ensures a speedy review of your submission. If you are still uncomfortable providing a public repository, you can work with your recruiter to provide access to
the reviewing engineer.

### How long do I have to complete the exercise?
There is no time limit for the exercise. Out of respect for your time, we designed this exercise with the intent that it should take you a few hours. But, please
take as much time as you need to complete the work.
