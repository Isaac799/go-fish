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

## Naming

Name things whatever you like, put them wherever you like. Just know that I make all names **lower kebab case for patterns**. So if a file is called `My photo.jpg` I will recognize that to `my-photo.jpg`.