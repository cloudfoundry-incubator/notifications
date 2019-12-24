# Notifications V1 Documentation

- System Status
	- [Check service status](#get-info)
- Sending Notifications
	- [Send a notification to a user](#post-users-guid)
	- [Send a notification to a space](#post-spaces-guid)
	- [Send a notification to an organization](#post-organizations-guid)
	- [Send a notification to all users in the system](#post-everyone-guid)
	- [Send a notification to a UAA-scope](#post-uaa-scopes)
	- [Send a notification to an email address](#post-emails)
	- [Check the status of a sent notification](#get-messages)
- Registering Notifications
	- [Register client notifications](#put-notifications)
- Updating Notifications
  - [Update a notification](#put-update-notification)
- Listing notifications
	- [List all notifications](#get-notifications)
- Managing User Preferences
	- [Retrieve options for /user_preferences endpoints](#options-user-preferences)
	- [Retrieve user preferences with a user token](#get-user-preferences)
	- [Update user preferences with a user token](#patch-user-preferences)
	- [Retrieve options for /user_preferences/{user-guid} endpoints](#options-user-preferences-guid)
	- [Retrieve user preferences with a client token](#get-user-preferences-guid)
	- [Update user preferences with a client token](#patch-user-preferences-guid)
- Managing Templates
	- [Create a new template](#post-template)
	- [Get a template](#get-template)
	- [Update a template](#put-template)
	- [Delete a template](#delete-template)
	- [List templates](#list-template)
	- [Get the default template](#get-default-template)
	- [Update the default template](#put-default-template)
	- [Assign a template to a client](#put-client-template)
	- [Assign a template to a notification](#put-client-notification-template)
	- [List template associations](#get-template-associations)

## System Status

<a name="get-info"></a>
#### Check service status

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
```

###### Route
```
GET /info
```

###### CURL example
```
$ curl -i -X GET \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  http://notifications.example.com/info

HTTP/1.1 200 OK
Connection: close
Content-Length: 13
Content-Type: application/json
Date: Tue, 30 Sep 2014 21:29:36 GMT
X-Cf-Requestid: 2cf01258-ccff-41e9-6d82-41a4441af4af

{"version": 1}
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields  | Description        |
| ------- | ------------------ |
| version | API version number |


## Sending Notifications

<a name="post-users-guid"></a>
#### Send a notification to a user

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.write` scope. Sending __critical__ notifications requires the `critical_notifications.write` scope.

###### Route
```
POST /users/{user-guid}
```
###### Params

| Key                | Description                                    |
| ------------------ | ---------------------------------------------- |
| kind_id\*          | a key to identify the type of email to be sent |
| text\*\*           | the text version of the email                  |
| html\*\*           | the html version of the email                  |
| subject\*          | the text of the subject                        |
| reply_to           | the Reply-To address for the email             |

\* required

\*\* either text or html have to be set, not both

###### CURL example
```
curl -i -X POST \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"kind_id":"example-kind-id", "subject":"what it is all about", "html":"this is a test"}' \
  http://notifications.example.com/users/user-guid

HTTP/1.1 200 OK
Connection: close
Content-Length: 129
Content-Type: application/json
Date: Tue, 30 Sep 2014 21:50:13 GMT
X-Cf-Requestid: 5c9bca88-280e-41d1-6e80-26a2a97adf4a

[{
	"notification_id":"451dd96a-ab8f-4a0b-5c3cb3bfe8ac1732",
	"recipient":"user-guid",
	"status":"queued"
}]
```
##### Response

###### Status
```
200 OK
```

###### Body
| Fields          | Description                               |
| --------------- | ----------------------------------------- |
| notification_id | Random GUID assigned to notification sent |
| recipient       | User GUID of notification recipient       |
| status          | Current delivery status of notification   |

----
<a name="post-spaces-guid"></a>
#### Send a notification to a space

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.write` scope. Sending __critical__ notifications requires the `critical_notifications.write` scope.

###### Route
```
POST /spaces/{space-guid}
```
###### Params

| Key                | Description                                    |
| ------------------ | ---------------------------------------------- |
| kind_id\*          | a key to identify the type of email to be sent |
| text\*\*           | the text version of the email                  |
| html\*\*           | the html version of the email                  |
| subject\*          | the text of the subject                        |
| reply_to           | the Reply-To address for the email             |

\* required

\*\* either text or html have to be set, not both

###### CURL example
```
$ curl -i -X POST \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"kind_id":"example-kind-id", "subject":"what it is all about", "html":"this is a test"}' \
  http://notifications.example.com/spaces/space-guid

HTTP/1.1 200 OK
Connection: close
Content-Length: 641
Content-Type: application/json
Date: Tue, 30 Sep 2014 22:01:34 GMT
X-Cf-Requestid: 4dcfc91c-9cf6-4a51-497a-8ae506ce37f5

[{
	"notification_id":"f44da2ff-e402-435d-54e8-8703970d5917",
	"recipient":"user-guid-1",
	"status":"queued"
 },
 {
 	"notification_id":"253305c8-eb72-4430-690e-76cbd8eae8ee",
 	"recipient":"user-guid-2",
 	"status":"queued"
}]
```
##### Response

###### Status
```
200 OK
```

###### Body
| Fields          | Description                               |
| --------------- | ----------------------------------------- |
| notification_id | Random GUID assigned to notification sent |
| recipient       | User GUID of notification recipient       |
| status          | Current delivery status of notification   |

----
<a name="post-organizations-guid"></a>
#### Send a notification to an organization

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.write` scope. Sending __critical__ notifications requires the `critical_notifications.write` scope.

###### Route
```
POST /organizations/{organization-guid}
```
###### Params

| Key                | Description                                    |
| ------------------ | ---------------------------------------------- |
| kind_id\*          | a key to identify the type of email to be sent |
| text\*\*           | the text version of the email                  |
| html\*\*           | the html version of the email                  |
| subject\*          | the text of the subject                        |
| reply_to           | the Reply-To address for the email             |

\* required

\*\* either text or html have to be set, not both

###### CURL example
```
$ curl -i -X POST \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"kind_id":"example-kind-id", "subject":"what it is all about", "html":"this is a test"}' \
  http://notifications.example.com/organizations/organization-guid

Connection: close
Content-Length: 897
Content-Type: application/json
Date: Thu, 06 Nov 2014 20:06:27 GMT
X-Cf-Requestid: 3a564cd9-74c8-46f6-5d31-8a8b600fc43f

[{
	"notification_id":"344f4b28-07d5-4490-468f-0a2f6fb4a65c",
	"recipient":"55498729-5749-4a4c-9e13-6893b795561b",
	"status":"queued"
	},{
	"notification_id":"96e633ef-8749-4dec-411a-f38a87f3fe79",
	"recipient":"d55067b8-cf2d-44ab-b70c-03dfd577a465",
	"status":"queued"
}]
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields          | Description                               |
| --------------- | ----------------------------------------- |
| notification_id | Random GUID assigned to notification sent |
| recipient       | User GUID of notification recipient       |
| status          | Current delivery status of notification   |

----

<a name="post-everyone-guid"></a>
#### Send a notification to all users in the system

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.write` scope. Sending __critical__ notifications requires the `critical_notifications.write` scope.

###### Route
```
POST /everyone
```
###### Params

| Key                | Description                                    |
| ------------------ | ---------------------------------------------- |
| kind_id\*          | a key to identify the type of email to be sent |
| text\*\*           | the text version of the email                  |
| html\*\*           | the html version of the email                  |
| subject\*          | the text of the subject                        |
| reply_to           | the Reply-To address for the email             |

\* required

\*\* either text or html have to be set, not both

###### CURL example
```
$ curl -i -X POST \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"kind_id":"example-kind-id", "subject":"what it is all about", "html":"this is a test"}' \
  http://notifications.example.com/everyone

Connection: close
Content-Length: 897
Content-Type: application/json
Date: Thu, 06 Nov 2014 20:06:27 GMT
X-Cf-Requestid: 3a564cd9-74c8-46f6-5d31-8a8b600fc43f

[{
	"notification_id":"344f4b28-07d5-4490-468f-0a2f6fb4a65c",
	"recipient":"55498729-5749-4a4c-9e13-6893b795561b",
	"status":"queued"
	},{
	"notification_id":"96e633ef-8749-4dec-411a-f38a87f3fe79",
	"recipient":"d55067b8-cf2d-44ab-b70c-03dfd577a465",
	"status":"queued"
}]
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields          | Description                               |
| --------------- | ----------------------------------------- |
| notification_id | Random GUID assigned to notification sent |
| recipient       | User GUID of notification recipient       |
| status          | Current delivery status of notification   |

----

<a name="post-uaa-scopes"></a>
#### Send a notification to a UAA Scope

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.write` scope. Sending __critical__ notifications requires the `critical_notifications.write` scope.

###### Route
```
POST /uaa_scopes/{scope}
```
###### Params

| Key                | Description                                    |
| ------------------ | ---------------------------------------------- |
| kind_id\*          | a key to identify the type of email to be sent |
| text\*\*           | the text version of the email                  |
| html\*\*           | the html version of the email                  |
| subject\*          | the text of the subject                        |
| reply_to           | the Reply-To address for the email             |

\* required

\*\* either text or html have to be set, not both

###### CURL example
```
$ curl -i -X POST \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"kind_id":"example-kind-id", "subject":"what it is all about", "html":"this is a test"}' \
  http://notifications.example.com/uaa_scopes/uaa.scope

Connection: close
Content-Length: 897
Content-Type: application/json
Date: Thu, 06 Nov 2014 20:06:27 GMT
X-Cf-Requestid: 3a564cd9-74c8-46f6-5d31-8a8b600fc43f

[{
	"notification_id":"344f4b28-07d5-4490-468f-0a2f6fb4a65c",
	"recipient":"55498729-5749-4a4c-9e13-6893b795561b",
	"status":"queued"
	},{
	"notification_id":"96e633ef-8749-4dec-411a-f38a87f3fe79",
	"recipient":"d55067b8-cf2d-44ab-b70c-03dfd577a465",
	"status":"queued"
}]
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields          | Description                               |
| --------------- | ----------------------------------------- |
| notification_id | Random GUID assigned to notification sent |
| recipient       | User GUID of notification recipient       |
| status          | Current delivery status of notification   |

----
<a name="post-emails"></a>
#### Send a notification to an email address

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `emails.write` scope

###### Route
```
POST /emails
```
###### Params

| Key                | Description                                    |
| ------------------ | ---------------------------------------------- |
| kind_id            | a key to identify the type of email to be sent |
| to\*               | The email address (and possibly full name) of the intended recipient in SMTP compatible format. |
| subject\*          | The desired subject line of the notification.  The final subject may be prefixed, suffixed, or truncated by the notifier, all dependent on the templates.|
| reply_to           | The email address to be included as the Reply-To address of the outgoing message. |
| text\*\*           | The message body, in plain text  (required if html is absent) |
| html\*\*           | The message body, in HTML  (required if text is absent) |

\* required

\*\* either text or html have to be set, not both

###### CURL example
```
$ curl -i -X POST \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"to":"user@example.com", "subject":"what it is all about", "html":"this is a test","kind_id":"my-notification"}' \
  http://notifications.example.com/emails

HTTP/1.1 200 OK
Connection: close
Content-Length: 108
Content-Type: application/json
Date: Tue, 30 Sep 2014 22:27:48 GMT
X-Cf-Requestid: eb7ee46c-2142-4a74-5b73-e4971eea511a

[{
	"recipient":"user@example.com",
	"notification_id":"86ad7892-8217-4359-54b1-fe3ca60d8ac9",
	"status":"queued"
}]
```
##### Response

###### Status
```
200 OK
```

###### Body
| Fields          | Description                               |
| --------------- | ----------------------------------------- |
| notification_id | Random GUID assigned to notification sent |
| recipient       | Email address of notification recipient   |
| status          | Current delivery status of notification   |


----
<a name="get-messages"></a>
#### Check the status of a sent notification

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires either the `emails.write` or the `notifications.write` scope

###### Route
```
GET /messages/{messageID}
```
###### Query parameters

| Key           | Description                                                             |
| --------------| ----------------------------------------------------------------------- |
| messageID\*   | The "notification_id" returned by any of the POST requests listed above |

\* required


###### CURL example
```
$ curl -i -X GET \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  http://notifications.example.com/messages/540cf340-03d3-4552-714f-0ec548a6cca9

200 OK
Connection: close
Content-Length: 22
Content-Type: application/json
Date: Tue, 20 Jan 2015 20:23:38 GMT
X-Cf-Requestid: 6869ab9a-c867-4271-6edd-d0c966bf7940
{"status":"delivered"}
```
##### Response

###### Status
```
200 OK
```

###### Body
| Fields          | Description                               |
| --------------- | ----------------------------------------- |
| status          | Current delivery status of notification   |

Possible `status` values:

| Value        | Meaning                                                                 |
| ------------ | ----------------------------------------------------------------------- |
| delivered    | Message delivered to the SMTP server (not necessarily the recipient)    |
| failed       | Message sending to SMTP server failed.                                  |
| queued       | Message has been added to a worker queue and will be processed shortly  |

In the case of "failed", the system will retry the delivery for up to 24 hours.

If the `messageID` is not known to the system, a `404 Not Found` response will be returned.

*Notification status info will be available for about 24 hours after a notification is first POSTed to this service. After 24 hours, status info is considered "stale" and may be purged by the system. A request for the status of a purged message will return a 404 Not Found error.*

## Registering Notifications

<a name="put-notifications"></a>
#### Register client notifications

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.write` scope. Registering __critical__ notifications requires the `critical_notifications.write` scope.

###### Route
```
PUT /notifications
```
###### Params

| Key                 | Description                                    |
| ------------------- | ---------------------------------------------- |
| source_name\* | The name of the sender, to be displayed in messages to users instead of the raw "client_id" field (which is derived from UAA) |
| notifications               | A list of notification types specified as a map (see table below for properties). |

\* required

###### Notifications Map Properties

| Key                       | Description |
| ------------------------- | ----------- |
| <name-of-notification>    | A key collecting the "description" and "critical" properties of a single notification |
| description\*              | A description of the notification, to be displayed in messages to users instead of the raw “id” field |
| critical (default: false) | A boolean describing whether this kind of notification is to be considered “critical”, usually meaning that it cannot be unsubscribed from.  Because critical notifications can be annoying to end-users, registering a critical notification kind requires the client to have an access token with the critical_notifications.write scope. |

\* required

###### CURL example
```
$ curl -i -X PUT \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"source_name":"Galactic Empire", "notifications":{"my-first-notification-id":{"description":"Example Kind Description", "critical": true}, "my-second-notification-id":{"description":"Example description", "critical":true}}}' \
  http://notifications.example.com/notifications


HTTP/1.1 204 OK
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 30 Sep 2014 22:47:50 GMT
X-Cf-Requestid: f39e22a4-6693-4a6d-6b27-006aecc924d4
```
##### Response

###### Status
```
204 No Content
```

## Updating Notifications

<a name="put-update-notification"></a>
#### Update a notification

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.manage` scope.

###### Route
```
PUT /clients/{client-id}/notifications/{notification-id}
```
###### Params

| Key                    | Description                                    |
| --------------------   | ---------------------------------------------- |
| description\*          | The description of the notification.           |
| critical\*             | A boolean describing whether this kind of notification is to be considered “critical”, usually meaning that it cannot be unsubscribed from.|
| template\*             | The GUID of the template to use when sending the notification.|

\* required

###### CURL example
```
$ curl -i -X PUT \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"description":"my excellent description", "critical":true, "template":"68C52741-C3C3-4B52-A522-787BF6159F72"}' \
  http://notifications.example.com/clients/a-good-client-id/notifications/my-notification-id


HTTP/1.1 204 OK
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 30 Sep 2014 22:47:50 GMT
X-Cf-Requestid: f39e22a4-6693-4a6d-6b27-006aecc924d4
```
##### Response

###### Status
```
204 No Content
```

## Listing Notifications

<a name="get-notifications"></a>
#### List all notifications
Returns all notifications in the system, grouped by client.  Clients without any notifications are also included.

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.manage` scope.

###### Route
```
GET /notifications
```

###### CURL example
```
$ curl -i \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  http://notifications.example.com/notifications

HTTP/1.1 200 OK
Connection: close
Content-Type: application/json
Date: Tue, 02 Dec 2014 21:40:59 GMT
X-Cf-Requestid: e3499f18-069a-4eed-720f-35baa61f1b5c
Transfer-Encoding: chunked

{
  "client-27054050-f813-4c0d-6fe4-2337f09a0aca": {
    "name": "Flynn",
    "template": "default",
    "notifications": {
      "clu": {
        "description": "CLU",
        "critical": false,
        "template": "default"
      },
      "grid": {
        "description": "A Digital Frontier...",
        "critical": false,
        "template": "EC6E8386-3096-48A4-A0C0-C0005B6933B2"
      },
      "mcp": {
        "description": "Master Control Program",
        "critical": true,
        "template": "C66DA695-C500-4D73-98F4-FC166EE0A0E9"
      }
    }
  },
 "client-36f1d81e-b9d6-400c-4f37-154ca2e8f01b": {
    "name": "Some other client",
    "template": "8BA02476-DC1F-493E-A6BF-EFE1D95ADFBD",
    "notifications": {
      "my-2nd-notification": {
        "description": "another test thingy",
        "critical": true,
        "template": "default"
      }
    }
  }
}

```
##### Response

###### Status
```
200 OK
```

###### Body
| Fields                    | Description                                                                 |
| ------------------------- | --------------------------------------------------------------------------- |
| client-id                 | Top-level keys are client GUIDs derived from UAA                            |
| name                      | The "source_name" set by the `PUT` method; displayed in messages to users   |
| template                  | The ID of the template assigned to the client                               |
| notifications             | A map, where the keys are notification IDs set by the `PUT` method          |
| notifications.description | A description of the notification.  Set by the `PUT` method                 |
| notifications.critical    | Boolean, indicating if notification is "critical".  Set by the `PUT` method |
| notifications.template    | The ID of the template assigned to the notification                         |


## Managing User Preferences

<a name="options-user-preferences"></a>
#### Retrieve Options for /user_preferences endpoints

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
```

###### Route
```
OPTIONS /user_preferences
```

###### CURL example
```
$ curl -i -X OPTIONS \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  http://notifications.example.com/user_preferences

HTTP/1.1 204 No Content
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 30 Sep 2014 22:54:40 GMT
X-Cf-Requestid: 686f601e-b6c7-4849-5699-6eed1a72004b
```
##### Response

###### Status
```
204 No Content
```

###### Headers
```
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
```
The above headers constitute a CORS contract. They indicate that the GET and PATCH endpoints for the `/user_preferences` path support the specified headers from any origin.

###### Body
| Fields   | Description |
| -------- | ----------- |
| \<none\> |             |

----
<a name="get-user-preferences"></a>
#### Retrieve user preferences with a user token

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <USER-TOKEN>
```
\* The user token requires `notification_preferences.write` scope.

###### Route
```
GET /user_preferences
```

###### CURL example
```
$ curl -i -X GET \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <USER-TOKEN>" \
  http://notifications.example.com/user_preferences

HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
Connection: close
Content-Length: 631
Content-Type: application/json
Date: Tue, 30 Sep 2014 23:19:11 GMT
X-Cf-Requestid: 92cffe86-16fe-41a8-4b80-b10987b11060

{
    "global_unsubscribe": false,
	"clients" : {
		"login-service": {
			"effa96de-2349-423a-b5e4-b1e84712a714": {
				"email": true,
				"kind_description": "Forgot Password",
				"source_description": "Login Service"
			}
		},
		"MySQL Service": {
			"6236f606-627d-4079-b0bd-f0b7e8d3d2a9": {
				"email": false,
				"kind_description": "Downtime Notification",
				"source_description": "Galactic Empire Datastore"
			},
			"fb89e98a-a1f5-47e5-9e2d-d95940b32d3d": {
				"email": true,
				"kind_description": "Provision Notification",
				"source_description": "Galactic Empire Datastore"
			}
		}
	}
}
```
##### Response

###### Status
```
200 OK
```

###### Headers
```
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
```
The above headers constitute a CORS contract. They indicate that the GET and PATCH endpoints for the `/user_preferences` path support the specified headers from any origin.

###### Response Body
| Fields             | Description                                                     |
| ------------------ | --------------------------------------------------------------- |
| global_unsubscribe | Boolean, indicates if user is unsubscribed to all notifications.  Overrides individual notification preferences |
| clients            | Map of clients

###### Client fields
| Fields             | Description |
| -------------------| ----------- |
| client_id          | Unique id of the client |
| kind_id            | Unique id of kind |
| email              | Indicates if the user is subscribed to receive the notification| 

----
<a name="patch-user-preferences"></a>
#### Update user preferences with a user token

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <USER-TOKEN>
```
\* The user token requires `notification_preferences.write` scope.

###### Route
```
PATCH /user_preferences
```

###### Request body
| Fields             | Description                                                     |
| ------------------ | --------------------------------------------------------------- |
| global_unsubscribe | Boolean, indicates if user is unsubscribed to all notifications.  Overrides individual notification preferences |
| clients            | Map of clients

###### Client fields
| Fields             | Description |
| -------------------| ----------- |
| client_id          | Unique id of the client |
| kind_id            | Unique id of kind |
| email              | Indicates if the user is subscribed to receive the notification| 

###### CURL example
```
$ curl -i -X PATCH \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <USER-TOKEN>" \
  -d '{"global_unsubscribe": false, "clients": {"login-service":{"effa96de-2349-423a-b5e4-b1e84712a714":{"email":true}}}}'
  http://notifications.example.com/user_preferences

HTTP/1.1 204 No Content
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 30 Sep 2014 23:19:11 GMT
X-Cf-Requestid: 92cffe86-16fe-41a8-4b80-b10987b11060
```
##### Response

###### Status
```
204 No Content
```

###### Headers
```
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
```
The above headers constitute a CORS contract. They indicate that the GET and PATCH endpoints for the `/user_preferences` path support the specified headers from any origin.

----
<a name="options-user-preferences-guid"></a>
#### Retrieve Options for /user_preferences/{user-guid} endpoints

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
```

###### Route
```
OPTIONS /user_preferences/{user-guid}
```

###### CURL example
```
$ curl -i -X OPTIONS \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  http://notifications.example.com/user_preferences/user-guid

HTTP/1.1 204 No Content
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 30 Sep 2014 23:07:22 GMT
X-Cf-Requestid: bfb28efe-757e-4b65-4d48-1d2c6d7a9ce6
```
##### Response

###### Status
```
204 No Content
```

###### Headers
```
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
```
The above headers constitute a CORS contract. They indicate that the GET and PATCH endpoints for the `/user_preferences/{user-guid}` path support the specified headers from any origin.

###### Body
| Fields   | Description |
| -------- | ----------- |
| \<none\> |             |

----
<a name="get-user-preferences-guid"></a>
#### Retrieve user preferences with a client token

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notification_preferences.admin` scope.

###### Route
```
GET /user_preferences/{user-guid}
```

###### CURL example
```
$ curl -i -X GET \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  http://notifications.example.com/user_preferences/user-guid

HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
Connection: close
Content-Length: 625
Content-Type: application/json
Date: Tue, 30 Sep 2014 23:19:11 GMT
X-Cf-Requestid: 92cffe86-16fe-41a8-4b80-b10987b11060

{
	"global_unsubscribe":false,	
	"clients": {
		"login-service": {
			"effa96de-2349-423a-b5e4-b1e84712a714": {
				"email": true,
				"kind_description": "Forgot Password",
				"source_description": "Login Service"
			}
		},
		"mysql-service": {
			"6236f606-627d-4079-b0bd-f0b7e8d3d2a9": {
				"email": false,
				"kind_description": "Downtime Notification",
				"source_description": "Galactic Empire Datastore"
			},
			"fb89e98a-a1f5-47e5-9e2d-d95940b32d3d": {
				"email": true,
				"kind_description": "Provision Notification",
				"source_description": "Galactic Empire Datastore"
			}
		}
	}
}
```
##### Response

###### Status
```
200 OK
```

###### Headers
```
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
```
The above headers constitute a CORS contract. They indicate that the GET and PATCH endpoints for the `/user_preferences/{user-guid}` path support the specified headers from any origin.

###### Response Body
| Fields             | Description                                                     |
| ------------------ | --------------------------------------------------------------- |
| global_unsubscribe | Boolean, indicates if user is unsubscribed to all notifications.  Overrides individual notification preferences |
| clients            | Map of clients

###### Client fields
| Fields             | Description |
| -------------------| ----------- |
| client_id          | Unique id of the client |
| kind_id            | Unique id of kind |
| email              | Indicates if the user is subscribed to receive the notification| 

----
<a name="patch-user-preferences-guid"></a>
#### Update user preferences with a client token

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notification_preferences.admin` scope.

###### Route
```
PATCH /user_preferences/user-guid
```

###### Request body
| Fields             | Description                                                     |
| ------------------ | --------------------------------------------------------------- |
| global_unsubscribe | Boolean, indicates if user is unsubscribed to all notifications.  Overrides individual notification preferences |
| clients            | Map of clients

###### Client fields
| Fields             | Description |
| -------------------| ----------- |
| client_id          | Unique id of the client |
| kind_id            | Unique id of kind |
| email              | Indicates if the user is subscribed to receive the notification| 

###### CURL example
```
$ curl -i -X PATCH \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <USER-TOKEN>" \
  -d '{"global_unsubcribe":false, "clients": {"login-service":{"effa96de-2349-423a-b5e4-b1e84712a714":{"email":true}}}}'
  http://notifications.example.com/user_preferences/user-guid

HTTP/1.1 204 No Content
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 30 Sep 2014 23:19:11 GMT
X-Cf-Requestid: 92cffe86-16fe-41a8-4b80-b10987b11060
```
##### Response

###### Status
```
204 No Content
```

###### Headers
```
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
```
The above headers constitute a CORS contract. They indicate that the GET and PATCH endpoints for the `/user_preferences/user-guid` path support the specified headers from any origin.

## Managing Templates

<a name="post-template"></a>
### Create Template

This endpoint is used to create a template and save it to the database.


##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notification_templates.write` scope

###### Route
```
POST /templates
```
###### Params

| Key      | Description                                                      |
| -------- | -----------------------------------------------------------------|
| name\*   | A human-readable template name                                   |
| html\*   | The template used for the HTML portion of the notification       |
| text     | The template used for the text portion of the notification       |
| subject  | An email subject template, defaults to "{{.Subject}}" if missing |
| metadata | Extra metadata to be stored alongside the template               |

\* required

###### CURL example
```
$ curl -i -X POST \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"name": "My template", "subject":"System notification: {{.Subject}}", "text":"Message to: {{.To}}, sent from the {{.ClientID}} UAA Client", "html": "<p>Message to: {{.To}}, sent from the {{.ClientID}} UAA Client</p>", "metadata": {"tags": "<h1>", "raptors": "scary"}}' \
  http://notifications.example.com/templates

201 Created
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 28 Oct 2014 00:18:48 GMT
X-Cf-Requestid: 8938a949-66b1-43f5-4fad-a91fc050b603


{"template-id": "E3710280-954B-4147-B7E2-AF5BF62772B5"}
```

##### Response

###### Status
```
201 Created
```

###### Body
| Fields      | Description             |
| ------------| ------------------------|
| template-id | A system-generated UUID |

<a name="get-template"></a>
### Get Template

This endpoint is used to retrieve a template that was saved to the database.


##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notification_templates.read` scope

###### Route
```
GET /templates/{my-template-id}
```
###### CURL example
```
$ curl -i -X GET \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  http://notifications.example.com/templates/my-template-id

200 OK
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 28 Oct 2014 00:18:48 GMT
X-Cf-Requestid: 8938a949-66b1-43f5-4fad-a91fc050b603


{
  "name":"My Custom Template",
  "subject" : "Hey! {{.Subject}}",
  "text" : "Dude! Stuff's Happening!",
  "html" : "\u003ch1\u003eHello!\u003c/h1\u003e",
  "metadata" : {
	"tag": "<h1>"
  }
}
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields      | Description                                  |
| ------------| ---------------------------------------------|
| name        | The human readable name of the template      |
| subject     | The subject for the template                 |
| text        | The plaintext representation of the template |
| html        | The HTML representation of the template *    |
| metadata    | Extra metadata stored alongside the template |

\* The HTML is Unicode escaped.  This is the expected behavior of the
[Golang JSON marshaller](http://golang.org/pkg/encoding/json/#Marshal)


<a name="put-template"></a>
### Update Template

This endpoint is used to update a template in the database.

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notification_templates.write` scope

###### Route
```
PUT /templates/templateID
```
###### Params

| Key      | Description                                                      |
| -------- | -----------------------------------------------------------------|
| name\*   | A human-readable template name                                   |
| subject  | An email subject template, defaults to "{{.Subject}}" if missing |
| html\*   | The template used for the HTML portion of the notification       |
| text     | The template used for the text portion of the notification       |
| metadata | Extra metadata stored alongside the template                     |

\* required

###### CURL example
```
$ curl -i -X PUT \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"name": "My template", "subject":"System notification: {{.Subject}}", "text":"Message to: {{.To}}, sent from the {{.ClientID}} UAA Client", "html": "<p>Message to: {{.To}}, sent from the {{.ClientID}} UAA Client</p>"}' \
  http://notifications.example.com/templates/templateID

204 No Content
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 28 Oct 2014 00:18:48 GMT
X-Cf-Requestid: 8938a949-66b1-43f5-4fad-a91fc050b603

```

##### Response

###### Status
```
204 No Content
```

###### Body

<a name="delete-template"></a>
### Delete Template

This endpoint is used to delete an existing template in the database.

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notification_templates.write` scope

###### Route
```
DELETE /templates/templateID
```

###### CURL example
```
$ curl -i -X DELETE \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  http://notifications.example.com/templates/template-id

204 No Content
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 28 Oct 2014 00:18:48 GMT
X-Cf-Requestid: 8938a949-66b1-43f5-4fad-a91fc050b603

```

##### Response
- If template is found and successfully deleted, then the response is `204 No Content`
- If template is not found, then the response is `404 Not Found`

<a name="list-template"></a>
### List Templates

This endpoint is used to retrieve a list of template id's and names that were saved to the database.

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notification_templates.read` scope

###### Route
```
GET /templates
```
###### CURL example
```
$ curl -i -X GET \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  http://notifications.example.com/templates

200 OK
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 28 Oct 2014 00:18:48 GMT
X-Cf-Requestid: 8938a949-66b1-43f5-4fad-a91fc050b603

{
  "F47CF7A7-43DE-4EA9-8B43-1A4C0964CDFB": {"name": "My Custom Template" },
  "584AB0E7-15EA-4BDA-B43F-BEB4EC301644": {"name": "Another template" }
}
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields      | Description                                  |
| ------------| ---------------------------------------------|
| template-id | The system-generated ID for a given template |
| name        | The human readable name of the template      |


<a name="get-default-template"></a>
### Get Default Template

This endpoint is used to retrieve the default template.

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notification_templates.read` scope

###### Route
```
GET /default_template
```
###### CURL example
```
$ curl -i -X GET \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  http://notifications.example.com/default_template

200 OK
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 28 Oct 2014 00:18:48 GMT
X-Cf-Requestid: 8938a949-66b1-43f5-4fad-a91fc050b603


{
  "name":"The Default Template",
  "subject" : "CF Notification: {{.Subject}}",
  "text" : "{{.Text}}",
  "html" : "{{.HTML}}",
  "metadata" : {}
}
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields      | Description                                  |
| ------------| ---------------------------------------------|
| name        | The human readable name of the template      |
| subject     | The subject for the template                 |
| text        | The plaintext representation of the template |
| html        | The HTML representation of the template *    |
| metadata    | Extra metadata stored alongside the template |

\* The HTML is Unicode escaped.  This is the expected behavior of the
[Golang JSON marshaller](http://golang.org/pkg/encoding/json/#Marshal)


<a name="put-default-template"></a>
### Update Default Template

This endpoint is used to update the default template.

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notification_templates.write` scope

###### Route
```
PUT /default_template
```
###### Params

| Key      | Description                                                      |
| -------- | -----------------------------------------------------------------|
| name\*   | A human-readable template name                                   |
| subject  | An email subject template, defaults to "{{.Subject}}" if missing |
| html\*   | The template used for the HTML portion of the notification       |
| text     | The template used for the text portion of the notification       |
| metadata | Extra metadata stored alongside the template                     |

\* required

###### CURL example
```
$ curl -i -X PUT \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"name": "My template", "subject":"System notification: {{.Subject}}", "text":"Message to: {{.To}}, sent from the {{.ClientID}} UAA Client", "html": "<p>Message to: {{.To}}, sent from the {{.ClientID}} UAA Client</p>"}' \
  http://notifications.example.com/default_template

204 No Content
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 28 Oct 2014 00:18:48 GMT
X-Cf-Requestid: 8938a949-66b1-43f5-4fad-a91fc050b603

```

##### Response

###### Status
```
204 No Content
```

<a name="put-client-template"></a>
### Assign a template to a client

This endpoint is used to assign an existing template to a known client.

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 1
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.manage` scope

###### Route
```
PUT /clients/:client_id/template
```
###### Params

| Key        | Description                                                                                |
| ---------- | -------------------------------------------------------------------------------------------|
| template\* | ID of template to be assigned (a value of `null` or `""` will assign the default template) |

\* required

###### CURL example
```
$ curl -i -X PUT \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"template": "4102591e-10d7-4c83-9fc9-1c88c5754f37"}' \
  http://notifications.example.com/clients/my-client/template

204 No Content
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 28 Oct 2014 00:18:48 GMT
X-Cf-Requestid: 8938a949-66b1-43f5-4fad-a91fc050b603

```

##### Response

###### Status
```
204 No Content
```

<a name="put-client-notification-template"></a>
### Assign a template to a notification

This endpoint is used to assign an existing template to a notification belonging to a known client.

##### Request

###### Headers
```
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.manage` scope

###### Route
```
PUT /clients/:client_id/notifications/:notification_id/template
```
###### Params

| Key        | Description                                                                                |
| ---------- | -------------------------------------------------------------------------------------------|
| template\* | ID of template to be assigned (a value of `null` or `""` will assign the default template) |

\* required

###### CURL example
```
$ curl -i -X PUT \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"template": "4102591e-10d7-4c83-9fc9-1c88c5754f37"}' \
  http://notifications.example.com/clients/my-client/notifications/my-notification/template

204 No Content
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 28 Oct 2014 00:18:48 GMT
X-Cf-Requestid: 8938a949-66b1-43f5-4fad-a91fc050b603

```

##### Response

###### Status
```
204 No Content
```

<a name="get-template-associations"></a>
### List template associations

This endpoint is used to list all clients and notifications associated to a template.

##### Request

###### Headers
```
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.manage` scope

###### Route
```
GET /templates/:template_id/associations
```
###### CURL example
```
$ curl -i -X GET \
  -H "X-NOTIFICATIONS-VERSION: 1" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  http://notifications.example.com/templates/template-id/associations

200 OK
Connection: close
Content-Length: 0
Content-Type: application/json
Date: Tue, 28 Oct 2014 00:18:48 GMT
X-Cf-Requestid: 8938a949-66b1-43f5-4fad-a91fc050b603

{"associations":[
    {"client":"client-id"},
    {"client":"client-id", "notification":"example-notification-id"},
    {"client":"client-id2", "notification":"example-notification-id2"}
  ]
}
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields                    | Description                                          |
| ------------------------- | ---------------------------------------------------- |
| associations              | The list of all associated clients and notifications |
| associations.client       | The client ID associated with this template          |
| associations.notification | The notification ID associated with this template    |
