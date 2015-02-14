app.EditAlbumView = Backbone.View.extend({
  events: {
    'submit form' : 'editAlbum'
  },

  editAlbum: function(e){
    e.preventDefault();
    var albumModel = new app.AlbumModel({
      "id": $("#albumId").val(),
      "name": $("#albumName").val(),
      "public": $("#isPublic").is(":checked")
    });

    albumModel.save({}, {
      success: function(album, response) {
        Backbone.trigger('flash', { message: "Success!", type: 'success' });
      },
      error: function(album, response) {
        Backbone.trigger('flash', { message: response.responseJSON.error, type: 'warning' });
      }
    });
  },

  render: function (album) {
      var albumTags = album.tags.join();

      dates = {
        start: new Date(album.date_range.start).toISOString().slice(0,10),
        end: new Date(album.date_range.end).toISOString().slice(0,10)
      }
      $(this.el).html(this.template({user: app.CurrentUser, album: album, albumTags: albumTags, dates: dates}));
      return this;
  }

});
