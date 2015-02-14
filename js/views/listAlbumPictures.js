app.ListAlbumPictures = Backbone.View.extend({
  events: {},
  initialize: function(objectId, pictures) {
    this.pictures = pictures;
    if (objectId != 0) {
      this.album = new app.AlbumModel({id: objectId});
      this.album.fetch();
      this.album.on('change', this.render, this);
    } else {
      this.album = {
        toJSON: function() {
          return {name: "My pictures", tags: ""}
        }
      };
    }
  },

  render: function () {
      $(this.el).html(this.template({album: this.album.toJSON(), pics: this.pictures}));
      return this;
  }
});
