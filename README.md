pit (Place in Time)
===================

Place in Time is a SPA for MongoDB course completion in Faculty of Mathematics and Informatics of Sofia University. The main idea behind the project is "Managing pictures by location and tags". Supports automatic album populating by location and tags criteria, sharing, picture like and unlike. It uses [MongoDB geolocation features](http://docs.mongodb.org/manual/applications/geospatial-indexes/) and [GridFS](http://docs.mongodb.org/manual/core/gridfs/) for storing files.

Main technologies used
----------------------

  * MongoDB
  * Go (+ [goji](https://goji.io/))
  * Backbone.js

TODO
-------------------

  * Add map on landing page with all user pictures.
  * Full edit of picture and album + resynchronization by location and tags.