import style from './index.scss';
import React, {Component} from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import API from '../api/index.js';
import showdown from 'showdown';

class CommentIndex extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.converter = new showdown.Converter();
    this.state = {comments: []};
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

  render() {
    return (
      <View state={this.state} />
    )
  }
}

const View = (props) => {
  const comments = props.state.comments.map((comment) => {
    return (
      <li key={comment.comment_id}>
        <div className={style.profile}>
          <img src={comment.user.avatar_url} alt={comment.user.nickname} className={style.avatar} />
          <div className={style.detail}>
            {comment.user.nickname}
            <div className={style.time}>
              <TimeAgo date={comment.created_at} />
            </div>
          </div>
        </div>
        <article className='md' dangerouslySetInnerHTML={{__html: comment.body}}>
        </article>
      </li>
    )
  });

  return (
    <div>
      <h3>Comments</h3>
      <ul className={style.comments}>
        {comments}
      </ul>
    </div>
  )
};

export default CommentIndex;
