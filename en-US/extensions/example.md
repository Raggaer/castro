---
name: Example
---

# Example

On this page we will go over a full example about creating a new extension, we will create a **custom login page that uses the email/password instead of the account name/password combo**

## Getting started

First we need to create our extension folder, head over to `extensions` and create a new folder, in this example we will use `login` to override the default login page.

So we need to have the following folders created:

```
extensions / login
    / pages
        get.lua
        login.html
```

## Manifest file

We need to create an `extension.json` file in order to be able to install our extension, this manifest file will look similar to this:

```json
{
    "name": "My extension name",
    "id": "login",
    "version": "0.0.1",
    "type": "3"
  }
```

**The ID field must be the same as your extension main folder**

## Login form

Starting the action, we first need to create a `get.lua` file and a `login.html` file. Inside the `get.lua` file we will render the HTML template and check for form errors:

```lua
function get()
    if session:isLogged() then
        http:redirect("/")
        return
    end

    local data = {}

    data.error = session:getFlash("error")

    http:render("login.html", data)
end
```

- First we check if the user is already logged in.
- Load all session flash errors (flash errors are one time use).
- Render http template and pass the data we want.

Now we need to create the HTML form, very basic one for this example:

```html
{{ template "header.html" . }}
<h3>Login to your account</h3>
<hr>
{{ if .validationError }}
<div class="alert alert-danger" role="alert">
    <strong>Error!</strong> {{ .validationError }}
</div>
{{ end }}
{{ if .success }}
<div class="alert alert-success" role="alert">
    <strong>Success!</strong> {{ .success }}
</div>
{{ end }}
<form method="POST" action="{{ url "subtopic" "login" }}">
    <input type="hidden" name="_csrf" value="{{ .csrfToken }}">
    <div class="form-group">
        <label for="input-account-name">Account email</label>
        <input type="text" class="form-control" id="input-account-name" name="account-email" placeholder="Account email">
    </div>
    <div class="form-group">
        <label for="input-account-name">Password</label>
        <input autocomplete="off" type="password" class="form-control" id="input-password" name="password" placeholder="Password">
    </div>
    <div class="form-group">
        <button type="submit" class="btn btn-primary">Login</button>
    </div>
</form>
{{ template "footer.html" . }} 
```

- We need to add a hidden input field for the CSRF token (security measures)

## Parsing the submitted login form

Now when a user clicks **Login** button nothing happens, we need to create a `post.lua` file in order to handle the form data:

```lua
function post()
    if session:isLogged() then
        http:redirect("/")
        return
    end

    local account = db:singleQuery(
        "SELECT name, secret FROM accounts WHERE email = ? AND password = ?", 
        http.postValues["account-email"], 
        crypto:sha1(http.postValues.password)
    )

    if account == nil then
        session:setFlash("error", "Wrong account name or password")
        http:redirect("/subtopic/login")
        return
    end

    session:set("logged", true)
    session:set("loggedAccount", account.name)
    session:set("admin", session:isAdmin())

    http:redirect("/subtopic/account/dashboard")
end
```

- We check if the user is already logged in.
- We then load an account by the provided password and email combination (we sha1 hash the password provided).
- If we dont find the account (`account == nil`) then we save an error to the flash session and redirect to the login form.
- We add some session fields to identify the user as logged into the system.

## Conclusion

On this brief example we made a new login form using the extensions system, this is a very simple example that doesnt cover everything **Castro** has to offer.

**Dont forget to install your extension (or it wont work!)**