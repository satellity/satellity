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
    this.api.comment.create(this.state, (resp) => {
      this.props.handleSubmit(resp.data);
      this.setState({body: '', submitting: false});
    });
  }

  render() {
    if (!this.api.user.loggedIn()) {
      return ''
    }
    return (
      <View onSubmit={this.handleSubmit} onChange={this.handleChange} state={this.state} />
    )
  }
}

const View = (props) => {
  return (
    <div className={style.form}>
      <form onSubmit={(e) => props.onSubmit(e)}>
        <input type='hidden' name='topic_id' defaultValue={props.state.topic_id} />
        <div>
          <textarea type='text' name='body' minLength='3' required placeholder='Say something ...' value={props.state.body} onChange={(e) => props.onChange(e)} />
        </div>
        <div className='action'>
          <button className='btn submit' disabled={props.state.submitting}>
            { props.state.submitting && <LoadingView style='sm-ring blank'/> }
            &nbsp;SUBMIT
          </button>
        </div>
      </form>
    </div>
  )
};

export default CommentNew;
