import style from './index.module.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, {Component} from 'react';
import TimeAgo from 'react-timeago';
import showdown from 'showdown';
import API from '../api/index.js';
import Avatar from '../users/avatar.js';
import CommentNew from './new.js';

class CommentIndex extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {
      user: this.api.user.local(),
      comments: [],
      comments_count: props.commentsCount
    };

    this.converter = new showdown.Converter();
    this.handleClick = this.handleClick.bind(this);
    this.handleActionClick = this.handleActionClick.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount() {
    this.api.comment.index(this.props.topicId).then((resp) => {
      if (resp.error) {
        return
      }
      const comments = resp.data.map((comment) => {
        comment.body = this.converter.makeHtml(comment.body);
        comment.handling = false;
        return comment
      });
      this.setState({comments: comments});
    });
  }

  handleActionClick(e, id) {
    e.preventDefault();
    const comments = this.state.comments.map((comment) => {
      if (comment.comment_id === id) {
        comment.handling = !comment.handling;
      } else {
        comment.handling = false;
      }
      return comment
    });
    this.setState({comments: comments});
  }

  handleClick(e, id) {
    e.preventDefault();
    this.api.comment.delete(id).then((resp) => {
      if (resp.error) {
        return
      }
      const comments = this.state.comments.filter(comment => comment.comment_id !== id);
      this.setState({comments: comments});
    })
  }

  handleSubmit(comment) {
    let newComments = this.state.comments.slice();
    comment.body = this.converter.makeHtml(comment.body);
    newComments.push(comment);
    this.setState({comments: newComments, comments_count: newComments.length});
  }

  render() {
    const i18n = window.i18n;
    const state = this.state;
    let comments = state.comments.map((comment) => {
      let action;
      if (state.user.user_id === comment.user_id) {
        action = (
          <span className={style.station}>
            <FontAwesomeIcon icon={['fas', 'ellipsis-v']} className={style.ellipsis} onClick={(e) => this.handleActionClick(e, comment.comment_id)} />
            {
              comment.handling &&
              <div className={style.actions}>
                <div onClick={(e) => this.handleClick(e, comment.comment_id)} className={style.action}>
                  <FontAwesomeIcon icon={['far', 'trash-alt']} className={style.trash} />
                  <span className={style.delete}>{i18n.t('general.delete')}</span>
                </div>
              </div>
            }
          </span>
        )
      }

      return (
        <li className={style.comment} key={comment.comment_id}>
          <div className={style.profile}>
            <Avatar user={comment.user} />
            <div className={style.detail}>
              {comment.user.nickname}
              <div className={style.time}>
                <TimeAgo date={comment.created_at} />
              </div>
            </div>
            {action}
          </div>
          <article className='md' dangerouslySetInnerHTML={{__html: comment.body}}>
          </article>
        </li>
      )
    });

    let commentsContainer = (
      <div className={style.container}>
        <h3>{i18n.t('comment.count', {count: state.comments_count})}</h3>
        <ul className={style.comments}>
          {comments}
        </ul>
      </div>
    )

    return (
      <div>
        {this.state.comments_count > 0 && commentsContainer}
        <CommentNew
          topicId={this.props.topicId}
          handleSubmit={this.handleSubmit} />
      </div>
    )
  }
}

export default CommentIndex;
