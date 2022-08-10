import style from './oauth.module.scss';
import React, {Component} from 'react';
import {Navigate} from 'react-router-dom';
import API from '../api/index.js';
import Loading from '../components/loading.js';

class Oauth extends Component {
  constructor(props) {
    super(props);
    // TODO
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('loading', 'layout');
    const params = new URLSearchParams('');
    this.state = {
      code: params.get('code'),
      redirect: false,
    };
  }

  componentDidMount() {
    // TODO should use redirect
    new API().user.signIn('', '', '', this.state.code).then((resp) => {
      this.setState({redirect: true});
    });
  }

  render() {
    if (this.state.redirect) {
      return (
        <Navigate to="/" replace />
      );
    }
    return (
      <div className={style.loading}>
        <Loading class='default' />
      </div>
    );
  }
}

export default Oauth;
