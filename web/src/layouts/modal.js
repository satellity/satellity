import style from './main.scss';
import React, { Component } from 'react';
import { loadReCaptcha, ReCaptcha } from 'react-recaptcha-v3';
import Config from '../components/config.js';

class Modal extends Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    loadReCaptcha(Config.ReCAPTCHASiteKey);
  }

  verifyCallback(recaptchaToken) {
    // Here you will get the final recaptchaToken!!!
    console.log(recaptchaToken, "<= your recaptcha token")
  }

  render() {
    return (
      <div className={style.modal}>
        <ReCaptcha
          sitekey={Config.ReCAPTCHASiteKey}
          action='login'
          verifyCallback={this.verifyCallback}
        />
        <div className={style.modalContainer}>
          <div onClick={this.props.handleLoginClick} className={style.action}>✕</div>
          <div className={style.app}>Login Satellity</div>
          <div className={style.content}>
            <a href={`https://github.com/login/oauth/authorize?scope=user:email&client_id=${Config.GithubClientId}`}>{i18n.t('login.github')}</a>
          </div>
        </div>
      </div>
    )
  }
}

export default Modal;
