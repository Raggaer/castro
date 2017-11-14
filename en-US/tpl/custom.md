---
name: Custom
---

# Custom templates

Creating your custom template is a very easy process. There are currently two routes you can take:

- Using an extension
- Overwrite default template

## Using an extension

Create a new extension over `extensions/my_custom_template`. Inside your folder we need to overwrite the default layout templates `header.html` and `footer.html`, in order to overwrite we need to create a page extension: `extensions/my_custom_template/pages/header.html` and `extensions/my_custom_template/pages/footer.html`.

Your template will go into these two html files.

## Overwrite default template

The main template is located at `views/default/header.html` and `views/default/footer.html` so you can just edit these files to change the default layout.

### Making a template compatible with Castro

A template is divided in two pieces: `header.html` and  `footer.html` where the content goes between those two. Below is a quick example:

```html
<!-- header.html -->
<html>
    <head>
        <title>My template</title>
    </head>
    <body>
        <div class="container">
```

```html
<!-- footer.html -->
        </div>
    </body>
</html>
```

All the content (subtopics) will be rendered inside `<div class="container">`. We also need to show the widgets (this is optional) by adding this piece of code (where we want them to appear):

```html
<!-- footer.html -->
        </div>
        <div class="widgets">
        {{ range $index, $element := .widgets }}
            {{ $element }}
        {{ end }}
        </div>
    </body>
</html>
```

Now the footer will show all the widgets inside `<div class="widgets">`. `.widgets` is a global template variable, [here is a list of the variables available on each request](/docs/tpl/variables#global-variables).

This is a very basic example, a template can also have [template hooks]() and [menu pages](). 