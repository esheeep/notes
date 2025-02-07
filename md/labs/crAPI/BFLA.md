# BFLA Delete a video of another user

- Leverage the predictable nature of REST APIs to find an admin endpoint to delete videos
- Delete a video of someone else

## Notes

Login to api with this endpoint `/identity/api/auth/login`, the server will return with bearer token.

Checking permissions, I can view my uploaded video at this endpoint `http://localhost:8888/identity/api/v2/user/videos/6` with my bearer token.
Using this endpoint I don't have access to any other videos.

There is an endpoint for admin to delete the videos `http://localhost:8888/identity/api/v2/admin/videos/<id>`.
Using my bearer token I can successfully delete my video. Trying a different video id such as `1`,`2`,`3`,`4`, I was able to successfully delete those videos as well.

```
{
    "message": "User video deleted successfully.",
    "status": 200
}
```
