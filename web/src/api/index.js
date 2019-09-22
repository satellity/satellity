import forge from 'node-forge';
import moment from 'moment';
import KJUR from 'jsrsasign';
import uuid from 'uuid/v4';
import Cookies from 'js-cookie';
import axios from 'axios';
import Noty from 'noty';
import Config from '../components/config.js';
import Category from './category.js';
import Topic from './topic.js';
import Comment from './comment.js';
import Group from './group.js';
import Message from './message.js';
import User from './user.js';
import Me from './me.js';
import Verification from './verification.js';

Noty.overrideDefaults({
    type: 'error',
    layout: 'topCenter',
    theme: 'mint',
    killer: true,
    theme: 'nest',
    timeout: 1000,
    progressBar: false,
    animation: {
      open : 'noty_effects_open',
      close: 'noty_effects_close'
    }
});

axios.defaults.baseURL = Config.ApiHost;
axios.defaults.headers.common['Content-Type'] = 'application/json';
axios.interceptors.request.use(function(config) {
  let method = config.method, url = config.url, data = config.data;
  config.headers.common['Authorization'] = `Bearer ${token(method, url, data)}`;
  return config
}, function(error) {
  return Promise.reject(error);
});

axios.interceptors.response.use(function(response) {
  if (!!response.status && (response.status >= 200 && response.status < 300)) {
    let data = response.data;
    if (!!data.error) {
      let error = data.error;
      if (error.code === 401) {
        window.location.href = `https://github.com/login/oauth/authorize?scope=user:email&client_id=${config.GithubClientId}`;
        return
      } else if (error.code === 404) {
        window.location.href = '/404'
        return
      }
      new Noty({
        text: error.code,
      }).show();
      if (error.code === 500) {
        window.location.href = '/';
        return
      }
      return Promise.reject(error);
    }
    return data;
  }
  return response
}, function(error) {
  if (error.response) {
    if (error.response.status === 500) {
      new Noty({
        text: "Internal Server Error"
      }).show();
      return
    }
  } else if (error.request) {
    new Noty({
      text: error.message
    }).show();
    return
  } else {
    new Noty({
      text: error.message
    }).show();
    return
  }
  return Promise.reject(error);
});

function token(method, uri, body) {
  let priv = window.localStorage.getItem('token');
  let pwd = Cookies.get('sid');
  if (!priv || !pwd) {
    return "";
  }
  Cookies.set('sid', pwd, { expires: 365 });

  let uid = window.localStorage.getItem('uid');
  let sid = window.localStorage.getItem('sid');
  return sign(uid, sid, priv, method, uri, body);
}

function sign(uid, sid, privateKey, method, uri, body) {
  if (typeof body !== 'string') { body = ""; }

  let expire = moment.utc().add(30, 'minutes').unix();
  let md = forge.md.sha256.create();
  md.update(method + uri + body);

  let oHeader = {alg: 'ES256', typ: 'JWT'};
  let oPayload = {
    uid: uid,
    sid: sid,
    exp: expire,
    jti: uuid(),
    sig: md.digest().toHex()
  };
  let sHeader = JSON.stringify(oHeader);
  let sPayload = JSON.stringify(oPayload);
  let pwd = Cookies.get('sid');
  try {
    let k = KJUR.KEYUTIL.getKey(privateKey, pwd);
  } catch (e) {
    return '';
  }
  return KJUR.jws.JWS.sign('ES256', sHeader, sPayload, privateKey, pwd);
}

class API {
  constructor() {
    this.axios = axios;
    this.category = new Category(this);
    this.topic = new Topic(this);
    this.comment = new Comment(this);
    this.group = new Group(this);
    this.message = new Message(this);
    this.user = new User(this);
    this.me = new Me(this);
    this.verification = new Verification(this);
  }
}

export default API;
