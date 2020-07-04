import forge from 'node-forge';
import moment from 'moment';
import KJUR from 'jsrsasign';
import { v4 as uuidv4 } from 'uuid';
import Cookies from 'js-cookie';
import axios from 'axios';
import Noty from 'noty';
import Category from './category.js';
import Topic from './topic.js';
import Comment from './comment.js';
import Product from './product.js';
import Client from './client.js';
import User from './user.js';
import Me from './me.js';
import Verification from './verification.js';

Noty.overrideDefaults({
    type: 'error',
    layout: 'topCenter',
    killer: true,
    theme: 'nest',
    timeout: 1000,
    progressBar: false,
    animation: {
      open : 'noty_effects_open',
      close: 'noty_effects_close'
    }
});

axios.defaults.headers.common['Content-Type'] = 'application/json';
axios.interceptors.request.use(function(config) {
  config.url = '/api' + config.url;
  const {method, url, data} = config;
  config.headers.common['Authorization'] = `Bearer ${token(method, url, data)}`;
  return config
}, function(error) {
  return Promise.reject(error);
});

axios.interceptors.response.use(function(response) {
  if (!!response.status && (response.status >= 200 && response.status < 300)) {
    const data = response.data;
    if (!!data.error) {
      const error = data.error;
      new Noty({
        text: window.i18n.t(`errors.${error.code}`)
      }).show();
      if (error.code === 401) {
        window.localStorage.removeItem('user')
        window.location.href = '/'
      } else if (error.code === 404) {
        window.location.href = '/404'
      }
    }
    return data;
  }
  return response
}, function(error) {
  let status, data;
  // TODO: should clear error.request and error
  if (error.response) {
    status = error.response.status;
    data = error.response.data;
  } else if (error.request) {
    status = 500;
    data = 'Initialize request error';
  } else {
    status = 500;
    data = error.message;
  }
  new Noty({
    text: window.i18n.t(`errors.${status}`)
  }).show();
  return {error: {code: status, description: data}};
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
    jti: uuidv4(),
    sig: md.digest().toHex()
  };
  let sHeader = JSON.stringify(oHeader);
  let sPayload = JSON.stringify(oPayload);
  let pwd = Cookies.get('sid');
  try {
    KJUR.KEYUTIL.getKey(privateKey, pwd);
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
    this.product = new Product(axios);
    this.user = new User(this);
    this.me = new Me(this);
    this.verification = new Verification(this);
    this.client = new Client(this);
  }
}

export default API;
