import style from './item.scss';
import moment from 'moment';
import React, {Component} from 'react';
import {Link, Redirect} from 'react-router-dom';
import TimeAgo from 'react-timeago';
import Avatar from '../users/avatar.js';
import API from '../api/index.js';
import LoadingView from '../loading/loading.js';

class Item extends Component {
  constructor(props) {
    super(props);
    this.state = {
      message: props.message,
      body: '',
      submitting: false
    }

    this.api = new API();
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(e) {
    let target = e.target;
    let name = target.name;
    this.setState({
      [name]: target.value
    });
  }

  handleSubmit(e) {
    e.preventDefault();
    let state = this.state;
    this.setState({submitting: true}, () => {
      this.api.message.create(state.message.group_id, {body: state.body, parent_id: state.message.parent_id}).then((data) => {
        let message = this.state.message;
        message.children.push(data);
        this.setState({submitting: false, body: '', message: message});
      });
    });
  }

  render() {
    let state = this.state;
    let children = state.message.children.sort((a, b) => {
      let l = moment(a.created_at);
      let r = moment(b.created_at);
      if (l < r) {
        return -1;
      }
      if (r > l) {
        return 1;
      }
      return 0;
    }).map((msg) => {
      return (
        <div key={msg.message_id}>
          <div className={style.subProfile}>
            <img src={msg.user.avatar_url} className={style.avatar} />
            <div className={style.subName}>{msg.user.nickname}</div>
            <TimeAgo date={msg.created_at} />
          </div>
          {msg.body}
        </div>
      )
    });

    return (
      <li className={style.message}>
        <div className={style.profile}>
          <Avatar user={state.message.user} />
          <div>
              {state.message.user.nickname}
            <div className={style.time}>
              <TimeAgo date={state.message.created_at} />
            </div>
          </div>
        </div>
        {state.message.body}
        <div className={style.replies}>
          <div>
            {children}
          </div>
          {
            this.props.current.message_id !== state.message.message_id && (
            <div onClick={() => this.props.handleComment(state.message.message_id)} className={style.commit}>
              add a comment
            </div>
            )
          }
          {
            this.props.current.message_id == state.message.message_id && (
            <div>
              <form onSubmit={this.handleSubmit}>
                <input type='hidden' name='group_id' defaultValue={state.message.group_id} />
                <input type='hidden' name='parent_id' defaultValue={state.message.message_id} />
                <div>
                  <textarea
                    className={style.input}
                    type='text'
                    name='body'
                    minLength='3'
                    required
                    value={state.body}
                    onChange={this.handleChange} />
                </div>
                <div className='action'>
                  <button className={`btn submit ${style.small}`} disabled={state.submitting}>
                    { state.submitting && <LoadingView style='sm-ring blank'/> }
                    &nbsp;{i18n.t('general.submit')}
                  </button>
                </div>
              </form>
            </div>
            )
          }
        </div>
      </li>
    )
  }
}

export default Item;
