import style from './oauth.module.scss';
import React, { Component } from 'react';
import { Redirect } from 'react-router-dom';
import API from '../api/index.js';
import Loading from '../components/loading.js';

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
      redirect: false,
    };
  }

  componentDidMount() {
    // TODO should use redirect
    const props = this.props;
    new API().user.signIn('','',props.match.params.provider,this.state.code).then((resp) => {
      this.setState({redirect: true});
    });
  }

  render() {
    if (this.state.redirect) {
      return (
        <Redirect to={{ pathname: "/" }} />
      )
    }
    return (
      <div className={style.loading}>
        <Loading class='default' />
      </div>
    );
  }
}

export default Oauth;
