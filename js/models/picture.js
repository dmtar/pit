"use strict";
app.PictureModel = Backbone.Model.extend({
  urlRoot: 'pictures/',

  save: function(attributes, options) {
    options = _.defaults((options || {}), {url: "/pictures/new"});
    return Backbone.Model.prototype.save.call(this, attributes, options);
  },

  canView: function() {
    var that = this;
    $.ajax({
      type: "GET",
      async: false,
      url: this.urlRoot + "canview/" + that.id,
      dataType: "json",
      success: function(data) {
        that.set("canView", true);
      },
      error: function(data) {
        Backbone.trigger('flash', { message: data.responseJSON.error, type: 'warning' });
        that.set("canView", false);
      }
    });
  }

});
