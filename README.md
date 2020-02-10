# GoModerator

GoModerator is a library that you can use to help moderate your applications!

## Example Use Case

An example use case was say you were making a reddit clone and you wanted a way to publically
expose how moderators determine whether or not a post should be removed.  This is where
the GoModerator library comes in.  GoModerator can create actions for real moderators to handle, and when
the moderators come to a conclusion, GoModerator will trigger a callback informing the application how to 
handle the resolution.  

The actions that are created by GoModerator are posted on forums such as Github Issues or Reddit itself!  These 
forums are chosen as they allow the developer to decide whether or not they wish to expose the moderating process
to non authenticated accounts or just other moderators.  