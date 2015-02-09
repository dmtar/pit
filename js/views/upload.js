app.PictureUploadView = Backbone.View.extend({
  events: {
    'submit form' : 'uploadFile'
  },

  uploadFile: function(event) {
    var values = {};
    if(event){ event.preventDefault(); }
		_.each(this.$('form').serializeArray(), function(input){
    	values[ input.name ] = input.value;
    })

	var pictureModel = new app.PictureModel();

    pictureModel.save(values, { iframe: true,
                              files: this.$('form :file'),
                              contentType: 'multipart/form-data',
                              data: values });
  },

  initialize: function() {
    app.CurrentUser.on('change', this.render, this)
  },

  render: function () {
      $(this.el).html(this.template({user: app.CurrentUser}));
      return this;
  }
});
