import style from './main.scss';
import React, { Component } from 'react';
import {Redirect} from 'react-router-dom';
import { loadReCaptcha, ReCaptcha } from 'react-recaptcha-v3';
import Config from '../components/config.js';
import API from '../api/index.js';

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
    this.handleClick = this.handleClick.bind(this);
  }

  handleClick() {
    this.setState({purpose: 'USER'}, () => {
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
    this.api.verification.create(this.state).then((data) => {
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
    this.api.user.verify(this.state).then((data) => {
      this.setState({success: true, submitting: false});
    });
    this.setState({submitting: true});
  }

  handleSignIn(e) {
    e.preventDefault();
    if (this.state.submitting) {
      return
    }
    this.api.user.signIn('', this.state.email, this.state.password).then((data) => {
      this.setState({success: true, submitting: false});
    });
    this.setState({submitting: true});
  }

  render() {
    let state = this.state;
    if (state.success) {
      return (
        <Redirect to={{pathname: "/"}} />
      )
    }

    let signIn = (
      <div>
        <div className={style.content}>
          <a href={`https://github.com/login/oauth/authorize?scope=user:email&client_id=${Config.GithubClientId}`}>{i18n.t('login.github')}</a>
        </div>
        <div className={style.or}>
          OR
        </div>
        <form onSubmit={this.handleSignIn}>
          <div>
            <input type='text' name='email' required value={state.email} autoComplete='off' placeholder='Username or Email *' onChange={this.handleChange} />
          </div>
          <div>
            <input type='password' name='password' required value={state.password} autoComplete='off' placeholder='Password *' onChange={this.handleChange} />
          </div>
          <div>
            <button type='submit' className='btn session' disabled={state.submitting}>
                &nbsp;{i18n.t('general.submit')}
            </button>
          </div>
        </form>
        <div className={style.register} onClick={this.handleClick}>
          Register A New Account
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
              <button type='submit' className='btn session' disabled={state.submitting}>
                  &nbsp;{i18n.t('general.submit')}
              </button>
            </div>
          </form>
        </div>
      </div>
    );

    let register = (
      <div>
        <form onSubmit={this.handleRegister}>
          <div>
            <input type='text' name='username' required value={state.username} autoComplete='off' placeholder='Username *' onChange={this.handleChange} />
          </div>
          <div>
            <input type='password' name='password' required value={state.password} autoComplete='off' placeholder='Password *' onChange={this.handleChange} />
          </div>
          <div>
            <input type='text' name='code' required value={state.code} autoComplete='off' placeholder='Verification Code *' onChange={this.handleChange} />
          </div>
          <div>
            <button type='submit' className='btn session' disabled={state.submitting}>
                &nbsp;{i18n.t('general.submit')}
            </button>
          </div>
        </form>
      </div>
    );

    return (
      <div className={style.modal}>
        <div className={style.modalContainer}>
          <div onClick={this.props.handleLoginClick} className={style.action}>✕</div>
          <div className={style.app}>Login</div>
          {state.purpose=='SESSION' && signIn}
          {state.purpose=='USER' && state.verification_id === '' && verification}
          {state.verification_id !== '' && register}
        </div>
      </div>
    )
  }
}

export default Modal;
