"use strict";
app.PictureViewView = Backbone.View.extend({

  details_template: _.template([
    '  <dt>Name</dt>                             ',
    '  <dd><%=picture.name || ""%></dd>          ',
    '  <% if (picture.location != undefined) {%> ',
    '    <dt>Location</dt>                       ',
    '    <dd><%=picture.location.name%></dd>     ',
    '  <%}%>                                     ',
    '  <dt>Tags</dt>                             ',
    '  <dd>                                      ',
    '    <%picture.tags.forEach(function(el){%>  ',
    '      <span class="label label-primary">    ',
    '        <%=el%>                             ',
    '      </span>                               ',
    '    <%});%>                                 ',
    '  </dd>                                     ',
    '  <dt>Date</dt>                             ',
    '  <dd><%=picture.date || ""%></dd>          ',
    '  <% if (album.name != undefined) {%>       ',
    '    <dt>Album</dt>                          ',
    '    <dd><%=album.name%></dd>                ',
    '  <%}%>                                     ',
    '  <dt>User</dt>                             ',
    '  <dd><%=user.display_name || ""%></dd>     ',
  ].join("\n")),
  initialize: function(objectId) {
    this.picture = new app.PictureModel({id: objectId});
    this.picture.fetch();
    this.picture.canView();
    this.album = new app.AlbumModel({id: this.picture.get("album")});
    this.user = new app.UserModel({id: this.picture.get("user")});
    this.picture.on('change', this.reloadAssociations, this);
    this.picture.on('change', this.updateDetails, this);
    this.album.on('change', this.updateDetails, this);
    this.user.on('change', this.updateDetails, this)
  },

  render: function () {
    $(this.el).html(this.template({picture: this.picture.toJSON(), album: this.album.toJSON()}));
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
  }

});
