"use strict";
app.AlbumModel = Backbone.Model.extend({
  urlRoot: 'albums/',

  save: function(attributes, options) {
    options = _.defaults((options || {}), {url: "/albums/new"});
    return Backbone.Model.prototype.save.call(this, attributes, options);
  }
});