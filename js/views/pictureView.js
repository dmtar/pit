"use strict";
app.PictureViewView = Backbone.View.extend({

  events: {
    "click #like": "likePicture"
  },

  details_template: _.template([
    '  <dt>Name</dt>                             ',
    '  <dd><%=picture.name || ""%></dd>          ',
    '  <% if (picture.location != undefined) {%> ',
    '    <dt>Location</dt>                       ',
    '    <dd><%=picture.location.name%></dd>     ',
    '  <%}%>                                     ',
    '  <% if (picture.tags != undefined) {%>     ',
    '  <dt>Tags</dt>                             ',
    '  <dd>                                      ',
    '    <%picture.tags.forEach(function(el){%>  ',
    '      <span class="label label-primary">    ',
    '        <%=el%>                             ',
    '      </span>                               ',
    '    <%});%>                                 ',
    '  </dd>                                     ',
    '  <%}%>                                     ',
    '  <dt>Date</dt>                             ',
    '  <dd><%=picture.date || ""%></dd>          ',
    '  <% if (album.name != undefined) {%>       ',
    '    <dt>Album</dt>                          ',
    '    <dd><%=album.name%></dd>                ',
    '  <%}%>                                     ',
    '  <dt>User</dt>                             ',
    '  <dd><%=user.display_name || ""%></dd>     ',
    '  <dt>Likes</dt>                            ',
    '  <dd><%=picture.likes.length || 0%></dd>   ',
  ].join("\n")),
  initialize: function(objectId) {
    document.body.style.overflow = 'scroll';
    this.picture = new app.PictureModel({id: objectId, likes: []});
    this.picture.fetch();
    this.picture.canView();
    this.album = new app.AlbumModel({id: this.picture.get("album")});
    this.user = new app.UserModel({id: this.picture.get("user")});
    this.picture.on('change', this.reloadAssociations, this);
    this.picture.on('change', this.render, this);
    this.album.on('change', this.updateDetails, this);
    this.user.on('change', this.updateDetails, this)
  },

  render: function () {
    $(this.el).html(this.template({
      picture: this.picture.toJSON(),
      album: this.album.toJSON(),
      user: app.CurrentUser.toJSON()
    }));
    this.updateDetails();
    return this;
  },

  reloadAssociations: function() {
    if (this.picture.get("album") != undefined) {
      this.album.id = this.picture.get("album");
      this.album.fetch();
    }
    this.user.id = this.picture.get("user");
    this.user.fetch();
  },

  updateDetails: function() {
    $("#details").html(this.details_template({picture: this.picture.toJSON(), album: this.album.toJSON(), user: this.user.toJSON()}));
  },

  likePicture: function(e) {
    if (e) { e.preventDefault(); }
    var that = this;
    $.ajax({
      type: "POST",
      url: "/pictures/like/" + that.picture.id,
      dataType: "json",
      success: function(data) {
        that.picture.fetch();
        Backbone.trigger('flash', { message: "Liked!", type: 'success' });
      },
      error: function(data) {
        Backbone.trigger('flash', { message: data.error, type: 'danger' });
      }
    });

  }

});
