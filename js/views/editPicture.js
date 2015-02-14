app.EditPictureView = Backbone.View.extend({
  events: {
    'submit form' : 'editPicture'
  },

  editPicture: function(e){
    e.preventDefault();
    var pictureModel = new app.PictureModel({
      "id": $("#pictureId").val(),
      "name": $("#name").val()
    });

    pictureModel.save({}, {
      success: function(picture, response) {
        Backbone.trigger('flash', { message: "Success!", type: 'success' });
        Backbone.history.navigate("#picture/view/"+picture.id, true);
      },
      error: function(picture, response) {
        Backbone.trigger('flash', { message: response.responseJSON.error, type: 'danger' });
      }
    });
  },

  render: function (picture) {

      var pictureTags = picture.tags.join();

      date = new Date(picture.date).toISOString().slice(0,10);

      $(this.el).html(this.template({user: app.CurrentUser, picture: picture, pictureTags: pictureTags, date: date}));
      return this;
  }

});
