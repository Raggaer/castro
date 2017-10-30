---
name: Static
---

# Static

An exntesion can also contain static assets (such as css and javascript files) you can ship your own static assets withing the `static` folder inside your extension root folder.

This aproach makes sure the main `static` folder of Castro stays clean. Below is an example of a static folder and how you would access the resources:

```html
my_extension /
    / static
        / css
            style.css
```

You can then access `style.css` using the following URL: `/extensions/my_extension/static/css/style.css`