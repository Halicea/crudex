---
theme: default
---

# TODOS

## [IN-PROGRESS] Choose a simple and nice UI library for the default scaffold
- To consider:
  - [#A] It should not be too heavy and use a lot of files and dependencies in order to create simple things.
  - [#B] The library should be easy to use
  - [#C] Preferably will not introduce another language besides html(this is hard to achieve these days)


 - Maybies:
   - HTMX seems to be a good candidate for making requests, no js is needed works ok with server side rendering
   - Alpine.Js seems to be a good candidate for the UI, low on js, tag based, it is simple to use, I just need to check how it plays with HTMX 

 - Maybies in case i go html only 
   - Flowbite seems to be a good candidate for the UI, just css with tailwind , looks nice as well
   - Foundation 

     I started with this one, and it seems good enough, the docs are quite hard to navigate, but the components seem ok.

## [ALMOST-DONE] Search, Filtering and pagination support
   - Without this we have nothing, the search and paagination needs to be easily controllable, to map well with gorm
- We need to be able to easily create a search query from a map of strings
- The operators for the search should map easily from request parameters to gorm operators
### The search is done using odata like query params

I'm using a third party lib for the parser and integration with GORM. 

The next steps would be to create a proper response structure for the search, and to create a way to easily integrate the search with the admin interface

## [TODO] Search contract and UI
- This is/(will be) done after the search and filtering is completed. 

## [TODO] Improve the admin interface
   - The default scaffold interface needs to be more complete, to provide better widgets for the different types of fields
- To have a better look and feel, and to look more like a real admin interface
- I am not working on this improvements until I decide on a good UI library/(combo of libs) to use.

## [TODO] Display relations in the admin interface

Also ability to annotate which relations will be displayed under the details view of the model and which not
What view should be used e.t.c

## [TODO] GraphQL support

   - We need to be able to easily create a graphql schema from our models.
- Libraries like this already exist, we need to make an integration  with some of them, no need to reinvent the wheel

# Things to think about

## Debug version should build the templates dynamically, so we can see the changes without restarting the server
The idea is to listen on the scaffold templates, and rebuild the templates when they change (just to re-register them with the renderer on every change).

This will allow for the whole ui to be auto refreshed after scaffold change, without the need to restart the server.

I expect very fast feedback loops with this, and it will be very useful for the development of the admin interface so it may be worth it to do it as a high prio item.
   - Considering that my UI skills are poor(mostly non existent) this will be a very useful tool to have.

## Do we need to get into the adventure of creating a CrudDB service(and interface) so we are not directly coupled to gorm?

## Model builder
   Maybe it is good if we create a ui where we can generate new models and relations. This will mean tthat we need to touch the users code, not so sure if it is going to be a good idea
   or if the users will gain something from it. Maybe we can just provide a way to generate the models from the command line, and let the user do the rest
   Someething like `rails generate model` but for go, and with the ability to generate relations
