app.ListAlbumsView = Backbone.View.extend({
  events: {},

  render: function (title, albums) {
      $(this.el).html(this.template({title: title, user: app.CurrentUser, albums: albums}));
      return this;
  }

});
