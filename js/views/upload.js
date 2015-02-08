app.PictureUploadView = Backbone.View.extend({
  events: {
  },

  initialize: function() {
    app.CurrentUser.on('change', this.render, this)
  },

  render: function () {
      $(this.el).html(this.template({user: app.CurrentUser}));
      return this;
  }
});
