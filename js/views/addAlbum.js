app.AddAlbumView = Backbone.View.extend({
  events: {
    'submit form' : 'addAlbum'
  },

  addAlbum: function(e){
    e.preventDefault();
    var albumModel = new app.AlbumModel();

    albumModel.save({
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
    },
    {
      success: function(album, response) {
        Backbone.trigger('flash', { message: "Success!", type: 'success' });
        Backbone.history.navigate("#album/list", true);
      },
      error: function(album, response) {
        Backbone.trigger('flash', { message: response.responseJSON.error, type: 'danger' });
      }
    });
  },

  render: function () {
      $(this.el).html(this.template({user: app.CurrentUser}));
      return this;
  }

});
