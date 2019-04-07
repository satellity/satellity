import style from './/index.scss';
import React, {Component} from 'react';
import { Link } from 'react-router-dom';
import API from '../api/index.js';
import LoadingView from '../loading/loading.js';

class CommentNew extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.state = {topic_id: props.topicId, body: '', submitting: false};
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
    if (this.state.submitting) {
      return
    }
    this.setState({submitting: true});
    this.api.comment.create(this.state).then((data) => {
      this.props.handleSubmit(data);
      this.setState({body: '', submitting: false});
    });
  }

  render() {
    let state = this.state;
    if (!this.api.user.loggedIn()) {
      return ''
    }
    return (
      <div className={style.form}>
        <form onSubmit={this.handleSubmit}>
          <input type='hidden' name='topic_id' defaultValue={state.topic_id} />
          <div>
            <textarea type='text' name='body' minLength='3' required placeholder='Say something ...' value={state.body} onChange={this.handleChange} />
          </div>
          <div className='action'>
            <button className='btn submit' disabled={state.submitting}>
              { state.submitting && <LoadingView style='sm-ring blank'/> }
              &nbsp;SUBMIT
            </button>
          </div>
        </form>
      </div>
    )
  }
}

export default CommentNew;
