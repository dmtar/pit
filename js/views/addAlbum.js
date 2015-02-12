app.AddAlbumView = Backbone.View.extend({
  events: {
    'submit form' : 'addAlbum'
  },

  addAlbum: function(){
  	var albumModel = new app.AlbumModel();
  	var values = {
  		"name": $("#albumName").val(),
  		"location": {
  			"lat": 0,
  			"lng": 0,
  			"name": $("#albumLocationName").val()
  		},
  		"tags": $("#albumTags").val(),
  		"date_range": {
  			"start": new Date($("#startDate").val()).toJSON(),
  			"end": new Date($("#endDate").val()).toJSON()
  		}
  	};
  	
  	albumModel.save(values, {
      success: function(res) {
        Backbone.trigger('flash', { message: res, type: 'success' });
      },
      error: function(res) {
        Backbone.trigger('flash', { message: res, type: 'error' });
      }
    });
  },

  render: function () {
      $(this.el).html(this.template({user: app.CurrentUser}));
      return this;
  }

});