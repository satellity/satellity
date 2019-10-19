import style from './main.module.scss';
import React, { Component } from 'react';
import {Redirect} from 'react-router-dom';
import { loadReCaptcha, ReCaptcha } from 'react-recaptcha-v3';
import Config from '../components/config.js';
import API from '../api/index.js';
import Href from '../widgets/href.js';
import Button from '../widgets/button.js';

class Modal extends Component {
  constructor(props) {
    super(props);

    this.state = {
      purpose: 'SESSION',
      recaptcha: '',
      email: '',
      verification_id: '',
      username: '',
      password: '',
      session_secret: '',
      code: '',
      success: false,
      submitting: false
    }

    this.api = new API();
    this.handleChange = this.handleChange.bind(this);
    this.verifyCallback = this.verifyCallback.bind(this);
    this.handleVerification = this.handleVerification.bind(this);
    this.handleSignIn = this.handleSignIn.bind(this);
    this.handleRegister = this.handleRegister.bind(this);
    this.handleResetPassword = this.handleResetPassword.bind(this);
    this.handleClick = this.handleClick.bind(this);
  }

  handleClick(e, purpose) {
    this.setState({password: '', purpose: purpose}, () => {
      if (Config.ReCAPTCHASiteKey !== '' && this.state.verification_id === '') {
        loadReCaptcha(Config.ReCAPTCHASiteKey);
      }
    });
  }

  verifyCallback(recaptchaToken) {
    this.setState({recaptcha: recaptchaToken});
  }

  handleChange(e) {
    const {target: {name, value}} = e;
    this.setState({
      [name]: value
    });
  }

  handleVerification(e) {
    e.preventDefault();
    if (this.state.submitting) {
      return
    }
    this.api.verification.create(this.state).then((resp) => {
      if (resp.error) {
        this.setState({submitting: false});
        return
      }
      let data = resp.data;
      data.submitting = false;
      this.setState(data);
    });
    this.setState({submitting: true});
  }

  handleRegister(e) {
    e.preventDefault();
    if (this.state.submitting) {
      return
    }
    this.setState({submitting: true});
    this.api.user.verify(this.state).then((resp) => {
      if (resp.error) {
        this.setState({submitting: false});
        return
      }
      this.setState({success: true, submitting: false});
    });
  }

  handleResetPassword(e) {
    e.preventDefault();
    if (this.state.submitting) {
      return
    }
    this.api.user.verify(this.state).then((resp) => {
      if (resp.error) {
        this.setState({submitting: false});
        return
      }
      this.setState({purpose: 'SESSION', verification_id: '', email: '', password: '', submitting: false});
    });
    this.setState({submitting: true});
  }

  handleSignIn(e) {
    e.preventDefault();
    this.setState({submitting: true});
    if (this.state.submitting) {
      return
    }
    this.api.user.signIn('', this.state.email, this.state.password).then((resp) => {
      if (resp.error) {
        this.setState({submitting: false});
        return
      }
      this.setState({success: true, submitting: false});
    });
  }

  render() {
    const i18n = window.i18n;
    let state = this.state;
    if (state.success) {
      return (
        <Redirect to={{pathname: "/"}} />
      )
    }

    let signIn = (
      <div>
        <div className={style.content}>
          <Href action={`https://github.com/login/oauth/authorize?scope=user:email&client_id=${Config.GithubClientId}`} class='button' text={i18n.t('login.github')} original />
        </div>
        <div className={style.or}>
          OR
        </div>
        <form onSubmit={this.handleSignIn}>
          <div>
            <input type='text' name='email' required value={state.email} autoComplete='off' placeholder={i18n.t('account.identity')} onChange={this.handleChange} />
          </div>
          <div>
            <input type='password' name='password' required value={state.password} autoComplete='off' placeholder={i18n.t('account.password')} onChange={this.handleChange} />
          </div>
          <div>
            <Button type='submit' class='submit' disabled={state.submitting} text={i18n.t('general.submit')} />
          </div>
        </form>
        <div className={style.register} onClick={(e) => this.handleClick(e, 'USER')}>
          {i18n.t('account.new')}
        </div>
        <div className={style.register} onClick={(e) => this.handleClick(e, 'PASSWORD')}>
            {i18n.t('account.reset.password')}
        </div>
      </div>
    );

    let verification = (
      <div>
        {
          Config.ReCAPTCHASiteKey !== '' &&
          <ReCaptcha
            sitekey={Config.ReCAPTCHASiteKey}
            action='login'
            verifyCallback={this.verifyCallback}
          />
        }
        <div>
          <form onSubmit={this.handleVerification}>
            <div>
              <input type='text' name='email' required value={state.email} autoComplete='off' placeholder='Your Email *' onChange={this.handleChange} />
            </div>
            <div>
              <Button type='submit' class='submit' disabled={state.submitting} text={i18n.t('general.submit')} />
            </div>
          </form>
        </div>
      </div>
    );

    let register = (
      <div>
        <form onSubmit={this.handleRegister}>
          <div>
            <input type='text' name='code' required value={state.code} autoComplete='off' placeholder={i18n.t('account.verification')} onChange={this.handleChange} />
          </div>
          <div>
            <input type='text' name='username' required value={state.username} autoComplete='off' placeholder={i18n.t('account.username')} onChange={this.handleChange} />
          </div>
          <div>
            <input type='password' name='password' required value={state.password} autoComplete='off' placeholder={i18n.t('account.password')} onChange={this.handleChange} />
          </div>
          <div>
            <Button type='submit' class='submit' disabled={state.submitting} text={i18n.t('general.submit')} />
          </div>
        </form>
      </div>
    );

    let password = (
      <div>
        <form onSubmit={this.handleResetPassword}>
          <div>
            <input type='text' name='code' required value={state.code} autoComplete='off' placeholder={i18n.t('account.verification')} onChange={this.handleChange} />
          </div>
          <div>
            <input type='password' name='password' required value={state.password} autoComplete='off' placeholder={i18n.t('account.password')} onChange={this.handleChange} />
          </div>
          <div>
            <Button type='submit' class='submit' disabled={state.submitting} text={i18n.t('general.submit')} />
          </div>
        </form>
      </div>
    );


    return (
      <div className={style.modal}>
        <div className={style.modalContainer}>
          <div onClick={this.props.handleLoginClick} className={style.action}>✕</div>
          <div className={style.app}>
              {state.purpose==='SESSION' && i18n.t('account.sign.in')}
              {state.purpose==='USER' && i18n.t('account.sign.up')}
              {state.purpose==='PASSWORD' && i18n.t('account.reset.password')}
          </div>
            {state.purpose==='SESSION' && signIn}
            {(state.purpose==='USER' || state.purpose==='PASSWORD') && state.verification_id === '' && verification}
            {state.purpose==='USER' && state.verification_id !== '' && register}
            {state.purpose==='PASSWORD' && state.verification_id !== '' && password}
        </div>
      </div>
    )
  }
}

export default Modal;
