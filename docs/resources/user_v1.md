---
page_title: "Fastly: user_v1"
---

# fastly_user_v1

Provides a Fastly User, representing the configuration for a user account for interacting with Fastly.

The User resource requires a login and name, and optionally a role.

## Example Usage

Basic usage:

```hcl
resource "fastly_user_v1" "demo" {
  login = "demo@example.com"
  name  = "Demo User"
}
```

## Argument Reference

The following arguments are supported:

* `login` - (Required, Forces new resource) The email address, which is the login name, of the User.
* `name` - (Required) The real life name of the user.
* `role` - (Optional) The role of this user. Can be `user` (the default), `billing`, `engineer`, or `superuser`. For detailed information on the abilities granted to each role, see [Fastly's Documentation on User roles](https://docs.fastly.com/en/guides/configuring-user-roles-and-permissions#user-roles-and-what-they-can-do).

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` – The ID of the User.

## Import

A Fastly User can be imported using their user ID, e.g.

```
$ terraform import fastly_user_v1.demo xxxxxxxxxxxxxxxxxxxx
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **login** (String) The email address, which is the login name, of this user.
- **name** (String) The real life name of the user.

### Optional

- **id** (String) The ID of this resource.
- **role** (String) The user-assigned permissions role. Can be `user` (the default), `billing`, `engineer`, or `superuser`.