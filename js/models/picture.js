"use strict";
app.PictureModel = Backbone.Model.extend({
  urlRoot: 'pictures/',

  save: function(attributes, options) {
    options = _.defaults((options || {}), {url: "/pictures/new"});
    return Backbone.Model.prototype.save.call(this, attributes, options);
  }

});