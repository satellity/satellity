import style from './index.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../api/index.js';

class UserEdit extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    const user = this.api.user.me();
    this.state = {nickname: user.nickname, biography: user.biography};
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    if (!this.api.user.loggedIn()) {
      props.history.push('/');
    }
  }

  componentDidMount() {
    this.api.user.show((resp) => {
      const user = resp.data;
      this.setState({nickname: user.nickname, biography: user.biography});
    });
  }

  handleChange(e) {
    const target = e.target;
    const name = target.name;
    this.setState({
      [name]: target.value
    });
  }

  handleSubmit(e) {
    e.preventDefault();
    const history = this.props.history;
    const data = {nickname: this.state.nickname, biography: this.state.biography};
    this.api.user.update(data, (resp) => {
      history.push('/');
    });
  }

  handleChange(e) {
    const target = e.target;
    const name = target.name;
    this.setState({
      [name]: target.value
    });
  }

  render() {
    return (
      <View onSubmit={this.handleSubmit} onChange={this.handleChange} state={this.state} />
    )
  }
}

const View = ({onSubmit, onChange, state}) => {
  return (
    <div className='container'>
      <main className='section main'>
        <div className={style.profile}>
          <h2>Update Profile</h2>
          <form onSubmit={(e) => onSubmit(e)}>
            <div>
              <label name='nickname'>Nickname</label>
              <input type='text' name='nickname' value={state.nickname} autoComplete='off' onChange={(e) => onChange(e)} />
            </div>
            <div>
              <label name='biography'>Biography</label>
              <textarea type='text' name='biography' value={state.biography} onChange={(e) => onChange(e)} />
            </div>
            <div className='action'>
              <input type='submit' value='SUBMIT' />
            </div>
          </form>
        </div>
      </main>
      <aside className='section aside'>
      </aside>
    </div>
  )
};

export default UserEdit;
