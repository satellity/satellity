const Oauth = {
  GithubDevelopmentClientId: 'b9b78f343f3a5b0d7c99',
  GithubProductionClientId: '71905afbd6e4541ad62b',
}

const Site = {
  Name: 'GoDiscourse',
  ApiDevelopmentHost: 'http://localhost:4000',
  ApiProductionHost: 'https://api.godiscourse.com',
}

let Config = {
  Name: 'Go Discourse',
  GithubClientId: function() {
    if (process.env.NODE_ENV === 'development') {
      return Oauth.GithubDevelopmentClientId;
    }
    return Oauth.GithubProductionClientId;
  },

  ApiHost: function() {
    if (process.env.NODE_ENV === 'development') {
      return Site.ApiDevelopmentHost;
    }
    return Site.ApiProductionHost;
  }
}

export default Config;
