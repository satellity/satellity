import style from './index.module.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../../api/index.js';

class AdminComment extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {comments: []};
    this.handleDelete = this.handleDelete.bind(this);
  }

  componentDidMount() {
    this.api.comment.admin.index().then((resp) => {
      if (resp.error) {
        return
      }
      this.setState({comments: resp.data});
    });
  }

  handleDelete(e, id, body) {
    e.preventDefault();
    let c = window.confirm(`Delete: ${body}`);
    if (c) {
      this.api.comment.admin.delete(id).then((resp) => {
        if (resp.error) {
          return
        }

        let comments = this.state.comments.filter((comments) => {
          return comments.comment_id !== id;
        });
        this.setState({comments: comments});
      });
    }
  }

  render() {
    const state = this.state;

    const listComments = state.comments.map((comment) => {
      return (
        <li key={comment.comment_id}>
            {comment.user.nickname} |
          <Link to={`/comments/${comment.comment_id}`}>{comment.body}</Link>
          <div className={style.time}>
              {comment.comment_id} | {comment.created_at} |
            <Link to='' onClick={(e) => this.handleDelete(e, comment.comment_id, comment.body)} >DELETE</Link>
          </div>
        </li>
      )
    });

    return (
      <div>
        <h1 className='welcome'>
            Here is the list of  comments.
        </h1>
        <div className='panel'>
          <ul>
              {listComments}
          </ul>
        </div>
      </div>
    );
  }
}

export default AdminComment;
