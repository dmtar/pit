app.AddAlbumView = Backbone.View.extend({
  events: {
    'submit form' : 'addAlbum'
  },

  addAlbum: function(e){
    e.preventDefault();
    var albumModel = new app.AlbumModel({
      "name": $("#albumName").val(),
      "location": {
        "lat": $('#albumLocationLat').val(),
        "lng": $('#albumLocationLng').val(),
        "name": $("#albumLocationName").val()
      },
      "public": $("#isPublic").is(':checked'),
      "tags": $("#albumTags").val(),
      "date_range": {
        "start": new Date($("#startDate").val()).toJSON(),
        "end": new Date($("#endDate").val()).toJSON()
      }
    });

    albumModel.save({
      success: function(album, response) {
        Backbone.trigger('flash', { message: "Success!", type: 'success' });
      },
      error: function(album, response) {
        Backbone.trigger('flash', { message: response.responseJSON.error, type: 'warning' });
      }
    });
  },

  render: function () {
      $(this.el).html(this.template({user: app.CurrentUser}));
      return this;
  }

});
