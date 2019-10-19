import style from './oauth.module.scss';
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
    let params = new URLSearchParams(props.location.search);
    this.state = {
      code: params.get('code'),
    };
  }

  componentDidMount() {
    // TODO should use redirect
    const history = this.props.history;
    new API().user.signIn(this.state.code).then((resp) => {
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
