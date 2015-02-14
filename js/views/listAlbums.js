app.ListAlbumsView = Backbone.View.extend({
  events: {},

  render: function (albums) {
      $(this.el).html(this.template({user: app.CurrentUser, albums: albums}));
      return this;
  }

});