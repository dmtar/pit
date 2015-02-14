app.ListAlbumsView = Backbone.View.extend({
  events: {},

  // initialize: function () {
  //   this.listenTo(app.AlbumModel, 'sync', this.render);
  // },

  render: function (albums) {
      $(this.el).html(this.template({user: app.CurrentUser, albums: albums}));
      return this;
  }

});