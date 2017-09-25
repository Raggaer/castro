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