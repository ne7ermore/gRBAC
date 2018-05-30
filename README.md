## RBAC
Role-Based Access Control

### Packages
* [mongo](https://gopkg.in/mgo.v2)


### Usage

Install packge
```
go get github.com/ne7ermore/gRBAC
```

Import
```
import (
    auth "github.com/ne7ermore/gRBAC"
)
```

Use
```
func createPerm() {
    p, err := auth.CreatePermisson("p1", "form1:abc:view")
    if err != nil {
        panic(err)
    }
    // handle p
}

func createRole() {
    r, err := auth.CreateRole("role1")
    if err != nil {
        panic(err)
    }

    // add permission
    if _, err := auth.Assign(r.Id.Hex(), p.Id.Hex()); err != nil {
        panic(err)
    }
}

func createUser() {
    u, err := auth.CreateUser("use1")
    if err != nil {
        panic(err)
    }

    // add role
    if _, err := auth.AddRole(u.Id.Hex(), r.Id.Hex()); err != nil {
        panic(err)
    }
}
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
