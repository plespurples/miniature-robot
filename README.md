# Orders and reservations backend

This is an another step to fully automated Purples reservation and ticket system. In front of you, you have a websocket server that deals with locking and unlocking seats during the reservation process.

## Events

There is a couple of events which will be sent to the connected clients when they happen. There is list of them with a quick description below this paragraph.

| Event          | Description                                                                                                           |
|----------------|-----------------------------------------------------------------------------------------------------------------------|
| locked         | broadcasted except for the author (action lock), tells the client that some seat was locked                           |
| unlocked       | broadcasted except for the author (action unlock), tells the client that some seat was unlocked                       |
| lockedForYou   | only for author (action lock), tells the client that the specified seat was successfully locked for that user         |
| unlockedForYou | only for author (action lock), tells the client that the specified seat was successfully unlocked                     |
| deleted        | only for one client, tells the client that its seats has been unlocked because of the reservation creation time limit |
| reserved       | broadcasted, tells the client that some seat was reserved and is no longer selectable                                 |
| unreserved     | broadcasted, tells the client that some seat was unreserved and is selectable now                                     |
| paid           | broadcasted, tells the client that some seat was paid and is no longer selectable                                     |
| unpaid         | broadcasted, tells the client that some seat was unpaid and is selectable now                                         |
| unauthorized   | only for author (specific actions), tells the client that it hasn't provided the authorization string                 |
