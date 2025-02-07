# Access details of another user’s vehicle

To solve the challenge, you need to leak sensitive information of another user’s vehicle.
Since vehicle IDs are not sequential numbers, but GUIDs, you need to find a way to expose the vehicle ID of another user.
Find an API endpoint that receives a vehicle ID and returns information about it.

# Notes

Endpoint: `/identity/api/v2/vehicle/abc327fc-42ed-4262-8ade-daeb4c7b24e4/location`

Response:

```
{"carId":"abc327fc-42ed-4262-8ade-daeb4c7b24e4",
"vehicleLocation":{"id":4,"latitude":"38.206348",
"longitude":"-84.270172"},
"fullName":"sheep",
"email":"a@mail.com"}
```

Community page (`/forum`) expose the commenters vehicle id.
![response from community posts api](<CleanShot 2025-02-07 at 23.55.49.png>)

By replacing my own vehicle id with the ids found on the community page I'm able to see other users vehicle information `/identity/api/v2/vehicle/<vehicleid>/location` ![response when replace the vehicle id](<CleanShot 2025-02-07 at 23.59.19.png>)
