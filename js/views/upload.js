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

    values['date'] = new Date(values['date']).toJSON();

    var pictureModel = new app.PictureModel();

    pictureModel.save(values, {
      iframe: true,
      files: this.$('form :file'),
      contentType: 'multipart/form-data',
      data: values,
      success: function(picture, response) {
        if (response.error != undefined) {
          Backbone.trigger('flash', { message: response.error, type: 'warning' });
        } else {
          Backbone.trigger('flash', { message: 'Upload was successful!', type: 'success' });
          Backbone.history.navigate("#picture/view/"+picture.id);
        }
      },
      error: function(picture, response) {
        Backbone.trigger('flash', { message: response.error, type: 'warning' });
      }
    });
  },

  initialize: function() {
    app.CurrentUser.on('change', this.render, this)
  },

  render: function () {
      $(this.el).html(this.template({user: app.CurrentUser}));
      return this;
  }
});
