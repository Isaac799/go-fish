# Go Fish

a fish themed ssr framework using go templating and htmlx

## Concepts

The main 3 things are:

- **Tuna**: is the big fish of the app, a page that people visit. Available via mux
  - consumes sardines
- **Sardine**: is the small fry, prefixed with "_"
  - a template to be used within other templates 
  - can be fetched with htmlx
- **Clown**fish  is the styling or css of the app. 
  - hashed for their name to enable browser cache
  - served independently, not embedded in template

Sardine and Clown scope is show in the example below.

## Example

See the example folder