import forge from 'node-forge';
import moment from 'moment';
import KJUR from 'jsrsasign';
import uuid from 'uuid/v4';
import Cookies from 'js-cookie';

function User(api) {
  this.api = api;
}

User.prototype = {
  signIn: function(code, callback) {
    var pwd = uuid().toLowerCase();
    var ec = new KJUR.crypto.ECDSA({'curve': 'secp256r1'});
    var pub = ec.generateKeyPairHex().ecpubhex;
    var priv = KJUR.KEYUTIL.getPEM(ec, 'PKCS8PRV', pwd);
    // TODO Why use 3059301306072a8648ce3d020106082a8648ce3d030107034200
    const params = {'session_secret': '3059301306072a8648ce3d020106082a8648ce3d030107034200' + pub, 'code': code};
    this.api.request('post', '/oauth/github', params, (resp) => {
      if (resp.data) {
        const data = resp.data;
        Cookies.set('sid', pwd, { expires: 365 });
        window.localStorage.setItem('token', priv);
        window.localStorage.setItem('uid', data.user_id);
        window.localStorage.setItem('sid', data.session_id);
        window.localStorage.setItem('user', btoa(JSON.stringify(data)));
      }
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  },

  update: function(params, callback) {
    this.api.request('post', '/me', params, (resp) => {
      if (resp.data) {
        const data = resp.data;
        window.localStorage.setItem('user', btoa(JSON.stringify(data)));
      }
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  },

  self: function(callback) {
    this.api.request('get', '/me', {}, (resp) => {
      if (resp.data) {
        const data = resp.data;
        window.localStorage.setItem('user', btoa(JSON.stringify(data)));
      }
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  },

  show: function(id, callback) {
    this.api.request('get', `/users/${id}`, {}, (resp) => {
      if (resp.data) {
        const data = resp.data;
        if (typeof callback === 'function') {
          callback(resp);
        }
      }
    });
  },

  topics: function(id, callback) {
    this.api.request('get', `/users/${id}/topics`, {}, (resp) => {
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  },

  adminIndex: function(callback) {
    this.api.request('get', '/admin/users', {}, (resp) => {
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  },

  ecdsa: function() {
    var priv = window.localStorage.getItem('token');
    var pwd = Cookies.get('sid');
    if (!priv || !pwd) {
      return "";
    }
    var ec = KJUR.KEYUTIL.getKey(priv, pwd);
    return KJUR.KEYUTIL.getPEM(ec, 'PKCS1PRV');
  },

  token: function(method, uri, body) {
    var priv = window.localStorage.getItem('token');
    var pwd = Cookies.get('sid');
    if (!priv || !pwd) {
      return "";
    }
    Cookies.set('sid', pwd, { expires: 365 });

    var uid = window.localStorage.getItem('uid');
    var sid = window.localStorage.getItem('sid');
    return this.sign(uid, sid, priv, method, uri, body);
  },

  sign: function(uid, sid, privateKey, method, uri, body) {
    if (typeof body !== 'string') { body = ""; }

    let expire = moment.utc().add(30, 'minutes').unix();
    let md = forge.md.sha256.create();
    md.update(method + uri + body);

    var oHeader = {alg: 'ES256', typ: 'JWT'};
    var oPayload = {
      uid: uid,
      sid: sid,
      exp: expire,
      jti: uuid(),
      sig: md.digest().toHex()
    };
    var sHeader = JSON.stringify(oHeader);
    var sPayload = JSON.stringify(oPayload);
    var pwd = Cookies.get('sid');
    try {
      var k = KJUR.KEYUTIL.getKey(privateKey, pwd);
    } catch (e) {
      return '';
    }
    return KJUR.jws.JWS.sign('ES256', sHeader, sPayload, privateKey, pwd);
  },

  clear: function() {
    window.localStorage.clear();
  },

  me: function() {
    const user = window.localStorage.getItem('user');
    if (user === undefined || user === null) {
      return {};
    }
    return JSON.parse(atob(user));
  },

  loggedIn: function() {
    const user = this.me();
    return user.user_id !== undefined || user.username !== undefined || user.nickname !== undefined;
  },

  isAdmin: function() {
    const user = this.me();
    return user.role === 'admin';
  }
}

export default User;
