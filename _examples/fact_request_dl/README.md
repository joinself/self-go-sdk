# Deep link based information request

Your app can request your users to share attested facts about themselves via a Deep link. To do this, you'll need to generate a qr request with the facts you want the user to respond with. You can find a list of updated valid facts and their respective sources [here](https://github.com/joinself/self-go-sdk/blob/main/fact/fact.go).


As part of this process, you have to share the generated Deep link with your users, and wait for a response.

## Running this example

In order to run this example, you must have a valid app id and private key. Self credentials are issued by the [Self Developer portal](https://developer.joinself.com/) when you create a new app.

Once you have your valid `SELF_APP_ID` and `SELF_APP_DEVICE_SECRET` you can run this example with:

```bash
$ SELF_APP_ID=XXXXX SELF_APP_DEVICE_SECRET=XXXXXXXX go run fact.go
```

## Process diagram

This diagram shows how does a Deep link based information request process works internally.

![Diagram](https://static.joinself.com/images/di_facts_diagram.png)


1. Generate Self information request Deep Link
2. Share generated Deep Link code with your user
3. The user clicks the deep link
4. The user will select the requested facts and accept sharing them with you.
5. The user’s device will send back a signed response with specific facts
6. Self SDK verifies the response has been signed by the user based on its public keys.
7. Self SDK verifies each fact is signed by the user / app specified on each fact.
8. Your app gets a verified response with a list of requested verified facts.
