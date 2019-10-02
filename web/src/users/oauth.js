import style from './oauth.scss';
import URLUtils from '../components/url.js';
import React, { Component } from 'react';
import API from '../api/index.js';
import Loading from '../widgets/loading.js';

class Oauth extends Component {
  constructor(props) {
    super(props);
    // TODO
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('loading', 'layout');
  }

  componentDidMount() {
    // TODO should use redirect
    const history = this.props.history;
    const code = new URLUtils().getUrlParameter('code');
    new API().user.signIn(code).then((resp) => {
      if (resp.error) {
        return
      }
      history.push('/');
    });
  }

  render() {
    return (
      <div className={style.loading}>
        <Loading class='default' />
      </div>
    );
  }
}

export default Oauth;
