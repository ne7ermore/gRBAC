## RBAC
Role-Based Access Control

### Packages
* [mongo](https://gopkg.in/mgo.v2)


### Usage


```
// install packge
> go get github.com/ne7ermore/gRBAC
```

```
import "github.com/ne7ermore/gRBAC/auth"

// init mongodb, permissions, roles, users
auth.init()

// more api infos please check below
```

## API
|Function|Description|
|--|--|
|CreatePermisson|Create a new permisson|
|GetPerm|get permission by id|
|UpdatePerm|update permission by id|
|CreateRole|Create a new role|
|GetRole|Get role by id|
|Assign|Assign a permission to a role|
|Revoke|Revoke a permission from the role|
|CreateUser|Create a new user|
|GetUser|Get one user by mongid|
|GetUserByUid|Get one user by Uid|
|AddRole|Add one role to a user|
|DelRole|Delete one role from the user|
|IsPrmitted|A user has a permission or not|
|GetAllPerms|Get all permissions|
|GetAllRoles|Get all roles|
|GetAllUsers|Get all users|
