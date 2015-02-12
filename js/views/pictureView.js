"use strict";
app.PictureViewView = Backbone.View.extend({

  initialize: function(objectId) {
    this.model = new app.PictureModel({id: objectId});
    app.CurrentUser.on('change', this.render, this)
  },

  render: function () {
      this.model.fetch();
      $(this.el).html(this.template({picture: this.model}));
      return this;
  }
});
