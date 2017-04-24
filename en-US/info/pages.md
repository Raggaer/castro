---
name: Custom pages
---

# Custom pages

Castro uses lua to handle all your pages. To create a new page navigate to pages directory and create a new folder with your subtopic name. The page will be accessed as `/subtopic/:name`. You can use multiple levels of directories

`pages/community/view` can be accessed as `/subtopic/community/view`

On your new folder you can then create a `get.lua` (to handle GET requests) or `post.lua` (to handle POST requests)

Each file must contain a function with the method they correspond to. `get.lua` files should have the `function get()` and `post.lua` files should have the `function post()`