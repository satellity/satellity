import style from './index.scss';
import React, { Component } from 'react';
import { Redirect } from 'react-router-dom';
import API from '../api/index.js';
import LoadingView from '../loading/loading.js';

class UserEdit extends Component {
  constructor(props) {
    super(props);

    this.api = new API();
    const user = this.api.user.readMe();
    this.state = {
      nickname: user.nickname,
      biography: user.biography,
      submitting: false
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount() {
    this.api.user.me().then((user) => {
      this.setState({nickname: user.nickname, biography: user.biography});
    });
  }

  handleChange(e) {
    e.preventDefault();
    const target = e.target;
    const name = target.name;
    this.setState({
      [name]: target.value
    });
  }

  handleSubmit(e) {
    e.preventDefault();
    if (this.state.submitting) {
      return
    }
    this.setState({submitting: true});
    const history = this.props.history;
    const data = {nickname: this.state.nickname, biography: this.state.biography};
    this.api.user.update(data).then((user) => {
      this.setState({submitting: false});
      history.push('/');
    });
  }

  render() {
    if (!this.api.user.loggedIn()) {
      return (
        <Redirect to={{ pathname: "/" }} />
      )
    }

    let state = this.state;

    return (
      <div className='container'>
        <main className='section main'>
          <div className={style.profile}>
            <h2>{i18n.t('user.edit')}</h2>
            <form onSubmit={this.handleSubmit}>
              <div>
                <label name='nickname'>{i18n.t('user.nickname')}</label>
                <input type='text' name='nickname' value={state.nickname} autoComplete='off' onChange={this.handleChange} />
              </div>
              <div>
                <label name='biography'>{i18n.t('user.biography')}</label>
                <textarea type='text' name='biography' value={state.biography} onChange={this.handleChange} />
              </div>
              <div className='action'>
                <button className='btn submit' disabled={state.submitting}>
                  {state.submitting && <LoadingView style='sm-ring blank'/>}
                  &nbsp;{i18n.t('general.submit')}
                </button>
              </div>
            </form>
          </div>
        </main>
        <aside className='section aside'>
        </aside>
      </div>
    )
  }
}

export default UserEdit;
