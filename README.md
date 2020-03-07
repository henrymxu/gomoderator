# GoModerator

GoModerator is a library that you can use to help publically moderate your applications!

## Example Use Case

An example use case was say you were making a reddit clone and you wanted a way to publically
expose how moderators determine whether or not a post should be removed.  This is where
the GoModerator library comes in.  GoModerator can create actions for real moderators to handle, and when
the moderators come to a conclusion, GoModerator will trigger a callback informing the application how to 
handle the resolution.  

The actions that are created by GoModerator are posted on forums such as Github Issues or Reddit itself!  These 
forums are chosen as they allow the developer to decide whether or not they wish to expose the moderating process
to non authenticated accounts or just other moderators. 

## How to use

GoModerator consists of two parts, a`forum` and the `moderator`. 

### Forums

Forums represent possible places where the `moderator` can post actions and users can discuss whether or not action
needs to be taken.  Forums must be able to support and be able to provide through an API in some form: 

- Posting Actions
- Commentating on Actions
- Voting on Posts/Comments
- Ability to mark as resolved 

Forums will be provided to the `moderator` through a builder interface that can be found in `forum/forum.go`

### Moderator

The moderator is the service that will post, identify, and resolve actions.  The moderator provides the following commands:

- CreateAction
- DoesActionExist
- StartActionsPollingService
- FindAndHandleNewlyResolvedActions

### CreateAction

The moderator will create an action on the provided forum (e.g create an issue on github) with the provided details.
All actions must have an associated unique ID related to the item they are actioning on (e.g if a post requires action, the unique ID representing the post).
If an action with the same item ID already exists, the moderator will not create a new action.

### FindAndHandleNewlyResolvedActions

The moderator will go through the current unresolved actions and identify the ones that have resolutions.
If the resolution is valid, the moderator will mark the action as closed, and invoke the provided callback with the identified resolution.
This method can be ran as a service with a custom polling frequency through the `StartActionsPollingService` method.

## Forums

### Github

Requirements for Github are:

Example `config.toml`: 

```
AccessToken = "<Account Access Token>"
AccountName = "gomoderatoraccount"
RepositoryOwner = "mygithubaccount"
RepositoryName = "myrepository"
IssueLabels = ["gomoderator", "automated", "action"]
```

### Reddit

Currently not implemented
