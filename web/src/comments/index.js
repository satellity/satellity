import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import style from './index.scss';
import React, {Component} from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import showdown from 'showdown';
import API from '../api/index.js';
import CommentNew from './new.js';

class CommentIndex extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.converter = new showdown.Converter();
    this.state = {comments: [], comments_count: props.commentsCount};
    this.handleClick = this.handleClick.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount() {
    this.api.comment.index(this.props.topicId, (resp) => {
      let comments = resp.data.map((comment) => {
        comment.body = this.converter.makeHtml(comment.body);
        return comment
      });
      this.setState({comments: comments});
    });
  }

  handleClick(id) {
    this.api.comment.delete(id, () => {
      let comments = this.state.comments.filter(comment => comment.comment_id !== id);
      this.setState({comments: comments});
    })
  }

  handleSubmit(comment) {
    let newComments = this.state.comments.slice();
    newComments.push(comment);
    this.setState({comments: newComments, comments_count: newComments.length});
  }

  render() {
    return (
      <div>
        {this.state.comments_count > 0 && <View api={this.api}
          state={this.state}
          user={this.api.user.me()}
          handleClick={this.handleClick} />}
          <CommentNew
            topicId={this.props.topicId}
            handleSubmit={this.handleSubmit} />
      </div>
    )
  }
}

const View = (props) => {
  const comments = props.state.comments.map((comment) => {
    let delAction = '';
    if (props.user.user_id === comment.user_id) {
      delAction = (
        <FontAwesomeIcon icon={['far', 'trash-alt']} className={style.delete} onClick={() => props.handleClick(comment.comment_id)} />
      )
    }
    return (
      <li className={style.comment} key={comment.comment_id}>
        <div className={style.profile}>
          <img src={comment.user.avatar_url} alt={comment.user.nickname} className={style.avatar} />
          <div className={style.detail}>
            {comment.user.nickname}
            <div className={style.time}>
              <TimeAgo date={comment.created_at} />
            </div>
          </div>
          {delAction}
        </div>
        <article className='md' dangerouslySetInnerHTML={{__html: comment.body}}>
        </article>
      </li>
    )
  });

  return (
    <div>
      <h3>{props.state.comments_count} Comments</h3>
      <ul className={style.comments}>
        {comments}
      </ul>
    </div>
  )
};

export default CommentIndex;
