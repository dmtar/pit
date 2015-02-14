app.ListAlbumPictures = Backbone.View.extend({
  events: {},
  initialize: function(objectId, pictures) {
    this.pictures = pictures;
    this.album = new app.AlbumModel({id: objectId});
    this.album.fetch();
    this.album.on('change', this.render, this);
  },

  render: function () {
      $(this.el).html(this.template({album: this.album.toJSON(), pics: this.pictures}));
      return this;
  }
});
