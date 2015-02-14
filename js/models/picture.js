"use strict";
app.PictureModel = Backbone.Model.extend({
  urlRoot: 'pictures/',

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
        Backbone.trigger('flash', { message: data.responseJSON.error, type: 'danger' });
        that.set("canView", false);
      }
    });
  }

});
