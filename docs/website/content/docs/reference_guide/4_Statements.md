---
title: "Statements"
date: 2019
weight: 2
sequence: true
---


## Statements

Our `MobileApp` does not have any detail yet on how it behaves. Let's use sysl statements to describe behaviour. Sysl supports following types of statements:
  * [Text](#text)
  * [Call](#Call)
  * [Return](#return-response)
  * [Control Statements](#control-statements)
  * [Arguments](#arguments)

#### Text
Use simple text to describe behaviour. See below for examples of text statements:
```
Server:
  Login:
    do input validation
    "Use special characters like \n to break long text into multiple lines"
    'Cannot use special characters in single quoted statements'
```

#### Call
A standalone service that does not interact with anybody is not a very useful service. Use the `call` syntax to show interaction between two services.

In the below example, MobileApp makes a call to backend Server which in turn calls database layer.
  
```
MobileApp:
  Login:
    Server <- Login

Server:
  Login(data <: LoginData):
    build query
    DB <- Query
    check result
    return Server.LoginResponse

  !type LoginData:
    username <: string
    password <: string

  !type LoginResponse:
    message <: string

DB:
  Query:
    lookup data
    return data
  Save:
    ...
```
See [/assets/call.sysl](/assets/call.sysl) for complete example.

Now we have all the ingredients to draw a sequence diagram. Here is one generated by sysl for the above example:

![](/assets/call-Seq.png)

See [Generate Diagrams](#generate-diagrams) on how to draw sequence and other types of diagrams using sysl.


## Control flows


## If/else

Sysl allows you to express high level of detail about your design. You can specify decisions, processing loops etc.

##### If, else
You can express an endpoint's critical decisions using IF/ELSE statement:
```
Server:
  HandleFormSubmit:
    validate input
    IF session exists:
      use existing session
    Else:
      create new session
    process input
```
See [/assets/if-else.sysl](/assets/if-else.sysl) for complete example.

`IF` and `ELSE` keywords are case-insensitive. Here is how sysl will render these statements:

![](/assets/if-else-Seq.png)

## For, Loop, Until, While

Express processing loop using FOR:
```
Server:
  HandleFormSubmit:
    validate input
    For each element in input:
      process element
```
See [/assets/for-loop.sysl](/assets/for-loop.sysl) for complete example.

`FOR` keyword is case insensitive. Here is how sysl will render these statements:

![](/assets/for-loop-Seq.png)

You can use `Loop`, `While`, `Until`, `Loop-N` as well (all case-insensitive).
