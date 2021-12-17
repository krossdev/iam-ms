# KrossIAM Message Services

This repo is message services component of KrossIAM.

## What is message services

When KrossIAM need send some data to external world like email, SMS and so on,
it is send a request to the message server, ask it to complete the job.

## Why message services

Send message from KrossIAM itself is not a technology problem and more easy,
the problem is lost the flexibility of message processing.

For example, when some event happen, like user account login abnormal,
we need notice user know this event, but how to notice? email? sms? or both?
it is hard to decide, somebody like this but somebody like other.
beacuse this, KrossIAM make no assumptions, just tell the message services
this event happend, the message services decide how to do it next.

## Message broker

The message is route by [nats](https://nats.io) message broker. since nats is
publish/subscribe system, more than one app can subscribe to same `subject`,
this give very powerful flexibility to KrossIAM message services.

Think user login abnormal example above, KrossIAM tell message services the
event happen, one app may subscribe to this event to send email notice,
another app subscribe to this event to send sms, if you don't like sms,
just stop sms app.

nats has support
[many programming langauge](https://docs.nats.io/using-nats/developer),
you can write a message subscribe app use your favorite language.

one more thing, nats has authentication and authoriztion bulit in.

## Message type, async and sync

KrossIAM message has 3 types: `action`, `event`, and `audit`, where `action`
is **sync** message, `event` and `audit` is **async** message.

`action` is KrossIAM ask message services to do something, there must
has one and only one subscriber to the action, in other words, the action
must be execute exactly once, no more no less.

`action` is *sync* message, when KrossIAM ask a action, it will be block
until action finish or timeout.

example action includes send email, send sms, and so on.

`event` is *async* message, there may be any number subscriber(include 0),
KrossIAM just publish the event to the message services, does not care about
the result.

`audit` is basically same as `event`, but used for publish more sensitive
information, set audit and event on difference nats *subject*, combine with
nats authorization, we can limit which app has permissions to subscribe.

## MSS and MSC

This repo has 2 directory, `mss` and `msc`, mss stand for *Message Services Server*,
and msc stand for *Message Service Client*.

**msc** is a `go` module, has integrated into KrossIAM app. **mss** is a
standalone app. both of them is nats client.

message flow:

```
msc => nats => mss
```

see README in `mss` for more information.
