app.ListAlbumPictures = Backbone.View.extend({
  events: {},

  render: function (pics) {
      $(this.el).html(this.template({pics: pics}));
      return this;
  }
});
