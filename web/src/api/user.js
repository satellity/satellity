import KJUR from 'jsrsasign';
import uuid from 'uuid/v4';
import Cookies from 'js-cookie';

class User {
  constructor(api) {
    this.api = api;
    this.admin = new Admin(api);
    this.fixed_schema_header = '3059301306072a8648ce3d020106082a8648ce3d030107034200';
  }

  get ecdsa() {
    let priv = window.localStorage.getItem('token');
    let pwd = Cookies.get('sid');
    if (!priv || !pwd) {
      return "";
    }
    let ec = KJUR.KEYUTIL.getKey(priv, pwd);
    return KJUR.KEYUTIL.getPEM(ec, 'PKCS1PRV');
  }

  signIn(code) {
    let pwd = uuid().toLowerCase();
    let ec = new KJUR.crypto.ECDSA({'curve': 'secp256r1'});
    let pub = ec.generateKeyPairHex().ecpubhex;
    let priv = KJUR.KEYUTIL.getPEM(ec, 'PKCS8PRV', pwd);
    let params = {'session_secret': this.fixed_schema_header + pub, 'code': code};
    return this.api.axios.post('/oauth/github', params).then((resp) => {
      const data = resp.data;
      Cookies.set('sid', pwd, { expires: 365 });
      window.localStorage.setItem('token', priv);
      window.localStorage.setItem('uid', data.user_id);
      window.localStorage.setItem('sid', data.session_id);
      window.localStorage.setItem('user', btoa(JSON.stringify(data)));
      return data;
    });
  }

  update(params) {
    return this.api.axios.post('/me', params).then((resp) => {
      window.localStorage.setItem('user', btoa(JSON.stringify(resp.data)));
      return resp.data;
    });
  }

  show(id) {
    return this.api.axios.get(`/users/${id}`).then((resp) => {
      return resp.data;
    });
  }

  me() {
    return this.api.axios.get('/me').then((resp) => {
      window.localStorage.setItem('user', btoa(JSON.stringify(resp.data)));
      return resp.data;
    })
  }

  readMe() {
    const user = window.localStorage.getItem('user');
    if (!user) {
      return {};
    }
    return JSON.parse(atob(user));
  }

  loggedIn() {
    const user = this.readMe();
    return user.user_id !== undefined || user.username !== undefined || user.nickname !== undefined;
  }

  isAdmin() {
    return this.readMe().role === 'admin';
  }

  topics(id) {
    return this.api.axios.get(`/users/${id}/topics`).then((resp) => {
      return resp.data;
    })
  }

  clear() {
    window.localStorage.clear();
  }
}

class Admin {
  constructor(api) {
    this.api = api;
  }

  index() {
    return this.api.axios.get('/admin/users').then((resp) => {
      return resp.data;
    })
  }
}

export default User;
