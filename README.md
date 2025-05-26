# Go Fish

a fish themed ssr framework using go templating and htmlx

## Concepts

The main 3 things are:


- **Tuna** is a big fish. Served as a page. Consumes Sardines
	- Identified by mime `[ text/html ]`
	- Not cahced
- **Sardine** is a small fish. Used by tuna. Smaller templates, served standalone too
	- Identified by mime `[ text/html ]` & `_` name prefix
	- Not cahced
- **Clown** is a decorative fish. Used in head of document
	- Identified by mime `[ text/css | text/javascript ]`
	- Is cached & name from hash
- **Anchovy** is supportive of the tuna
	- Identified by mime `[ image | audio | video ]`
	- Is cached

## Example

See the example folder

## Scope

A **Tuna** has access to all other local fish, and top level fish.

## Path Value

For a file name, `.` delimited makes a new path. This is designed with path values in mind, and not to be used in place of dir structure. The even items are considered a value, enforcing a `/context/value` pattern. 

To use path values simply name a file like: 
- `user.id.html` translates to `/user/{id}`
- `user.id.edit.html` translates to `/user/{id}/edit`

How is this useful? 
- You can use the `id` in the license to restrict access
- You can use the `id` in the bait to lookup a user to help render.

## Naming

Name things whatever you like, put them wherever you like. Just know this:

- **Patterns** are lower kebab case. So if a file is called `My photo.jpg` I will recognize that to `my-photo.jpg`.
- **Template Names** are the same as file, without extension. So `about page.html` should be defined as `about page` in template.