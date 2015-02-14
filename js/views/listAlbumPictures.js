app.ListAlbumPictures = Backbone.View.extend({
  events: {},

  render: function (pics) {
  	  console.log(pics);
      $(this.el).html(this.template({pics: pics}));
      return this;
  }
});