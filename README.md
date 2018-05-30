## RBAC
Role-Based Access Control

### Packages
* [mongo](https://gopkg.in/mgo.v2)


### Usage

install packge
```
go get github.com/ne7ermore/gRBAC
```

import
```
import (
    auth "github.com/ne7ermore/gRBAC"
)
```

We use MongoDB in this Project, and u can edit [db settings](https://github.com/ne7ermore/gRBAC/blob/master/common/const.go#L26) by urself

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
